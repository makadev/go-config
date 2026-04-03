# go-config

A lightweight, generic Go library for loading and managing application configuration. It provides a thin layer of structure on top of standard YAML/JSON file loading, with first-class support for environment variable overrides and user-friendly configuration dumping.

## Scope

This library is a simple, composable config layer for Go services, not an all-in-one configuration platform.

## When to use go-config

Use this library when your service needs:
- A typed config struct loaded from local YAML/JSON files.
- Environment variable overrides for deployment-specific values.
- Safe config dumps with masked secrets for debugging and support.

Consider a larger config framework when you need:
- Remote config sources (Vault, Consul, AWS Parameter Store, etc.).
- Live reload and distributed config updates.
- Built-in schema validation, policy enforcement, or advanced config lifecycles.

## Installation

```bash
go get github.com/makadev/go-config
```

## Usage

The following examples cover the main features of this library. All examples use this import:

```go
import config "github.com/makadev/go-config"
```

---

### Generic `Config[T]` type

`Config[T]` is fully typed. The `Data` field gives direct, type-safe access to your config struct without any casting.

```go
type AppConfig struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

cfg, err := config.NewConfig(config.NewOptions(), &AppConfig{
    Host: "localhost",
    Port: 8080,
})
if err != nil {
    log.Fatal(err)
}
if err := cfg.Load(); err != nil {
    log.Fatal(err)
}

// cfg.Data is *AppConfig – no type assertion needed
fmt.Println(cfg.Data.Host, cfg.Data.Port)
```

### Thread safety

Calls that go through `Config` methods are synchronized internally:

- `Load`
- `Dump` / `DumpWithOptions` / `DumpEnv`
- `GetFieldValue` / `SetFieldValue`
- `GetConfigValue` / `SetConfigValue`

Direct access to exported members is **not** synchronized:

- `cfg.Data`
- `cfg.Metadata`
- `cfg.Options`

To protect direct field access from multiple goroutines, use the locking helpers built into `Config`:

```go
// WithLock / WithRLock — preferred; lock is always released via defer
cfg.WithLock(func() {
    cfg.Data.Counter++
})

var snapshot MyConfig
cfg.WithRLock(func() {
    snapshot = *cfg.Data
})

// Raw Lock / Unlock when defer-based scoping does not fit
cfg.Lock()
cfg.Data.Counter++
cfg.Unlock()

// Raw RLock / RUnlock
cfg.RLock()
value := cfg.Data.Counter
cfg.RUnlock()
```

> **Deadlock warning:** Do **not** call any `Config` methods (`Load`, `Dump`, `Get*`, `Set*`, …) inside a `WithLock`/`WithRLock` callback or between `Lock`/`Unlock` (or `RLock`/`RUnlock`) calls. Those methods acquire the same mutex internally, and `sync.RWMutex` is not reentrant — doing so will deadlock.

> **Race warning:** Mixing direct field access (`cfg.Data`, `cfg.Options`, `cfg.Metadata`) with concurrent method calls — without holding the lock — is a **data race**. Always use the helpers above, or the raw lock/unlock pairs, whenever you need to read or write exported fields from multiple goroutines.

> **`Options` note:** `Config.Options` should be configured before the first `Load` call and not mutated afterwards. Changing `Options` fields while `Load` is running in another goroutine is not safe even with the locking helpers.

---

### File loading (YAML & JSON)

Pass a list of file paths in `Options.ConfigPaths`. The first file that exists is loaded; its format is detected from the extension (`.yaml`, `.yml`, or `.json`).

```go
opts := config.NewOptions()
opts.ConfigPaths = []string{"config.yaml"} // or "config.json"

cfg, _ := config.NewConfig(opts, &AppConfig{})
cfg.Load() // reads config.yaml into cfg.Data
```

**config.yaml**

```yaml
host: "0.0.0.0"
port: 9090
```

**config.json** (same fields, JSON format)

```json
{ "host": "0.0.0.0", "port": 9090 }
```

---

### Priority file lookup

If `ConfigPaths` is left at its default, go-config searches for files in this order and loads the first one it finds:

```
config.local.yaml  →  config.local.json  →  config.yaml  →  config.json
```

`config.local.*` files take precedence, so local developer overrides work without any extra code.

```go
// Uses default priority order; config.local.yaml wins if it exists
cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{})
cfg.Load()
```

You can provide a fully custom priority list:

```go
opts := config.NewOptions()
opts.ConfigPaths = []string{
    "/etc/myapp/config.yaml", // system-level file
    "config.yaml",            // project-local file
}
```

---

### Environment variable loading

Tag each field with `env:"VAR_NAME"`. After files are loaded, matching environment variables override the values.

```go
type AppConfig struct {
    Host   string `yaml:"host"   env:"APP_HOST"`
    Port   int    `yaml:"port"   env:"APP_PORT"`
    APIKey string `yaml:"apikey" env:"APP_API_KEY" secret:"true"`
}

// export APP_HOST=prod.example.com
// export APP_PORT=443
cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{Host: "localhost", Port: 8080})
cfg.Load() // cfg.Data.Host == "prod.example.com", cfg.Data.Port == 443
```

---

### `AutoEnv` mode

Setting `Options.AutoEnv = true` derives env var names automatically from the full config key, so no `env` tags are required. Dots become underscores and the name is uppercased.

```go
type ServerConfig struct {
    Host string `yaml:"host"` // env var → SERVER_HOST
    Port int    `yaml:"port"` // env var → SERVER_PORT
}
type AppConfig struct {
    Server ServerConfig `yaml:"server"`
}

opts := config.NewOptions()
opts.AutoEnv = true

// export SERVER_HOST=api.example.com
cfg, _ := config.NewConfig(opts, &AppConfig{})
cfg.Load() // cfg.Data.Server.Host == "api.example.com"
```

An optional `EnvPrefix` prepends a global prefix to every auto-generated name:

```go
opts.AutoEnv    = true
opts.EnvPrefix  = "MYAPP_"
// "server.port" → MYAPP_SERVER_PORT
```

---

### Struct-as-env-prefix pattern

When `AutoEnv` is `false` (the default), placing `env:"PREFIX_"` on a struct field makes that value the prefix for all nested fields' env var names.

```go
type ServerConfig struct {
    Host string `env:"HOST"` // resolved as SERVER_HOST
    Port int    `env:"PORT"` // resolved as SERVER_PORT
}
type AppConfig struct {
    Server ServerConfig `env:"SERVER_"`
}

// export SERVER_HOST=0.0.0.0
// export SERVER_PORT=443
cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{})
cfg.Load()
```

---

### Struct tag mapping

The config key for a field is resolved in this priority order: `config` → `yaml` → `json`. If none are present, the lowercased field name is used.

```go
type AppConfig struct {
    Host string `config:"hostname" yaml:"host" json:"host"` // config key = "hostname"
    Port int    `yaml:"port"`                               // config key = "port"
    Name string                                             // config key = "name"
}
```

---

### Configurable tag priority

`Options.ConfigTags` controls which struct tags are checked and in what order. The default is `["config", "yaml", "json"]`.

```go
opts := config.NewOptions()
opts.ConfigTags = []string{"json"} // only look at json tags for config key names

type AppConfig struct {
    Host string `json:"host" yaml:"hostname"` // config key = "host" (json wins)
}
```

---

### `SkipFiles` / `SkipEnv` flags

Either loading source can be individually disabled.

```go
// Only use environment variables – ignore all config files
optsEnvOnly := config.NewOptions()
optsEnvOnly.SkipFiles = true

// Only use config files – ignore all environment variables
optsFilesOnly := config.NewOptions()
optsFilesOnly.SkipEnv = true
```

---

### Secret masking in dumps

Tag any field with `secret:"true"` to have its value replaced with `***` in all dump output by default.

```go
type AppConfig struct {
    Host     string `yaml:"host"`
    Password string `yaml:"password" env:"DB_PASSWORD" secret:"true"`
}

cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{Host: "localhost", Password: "s3cr3t"})
cfg.Load()

output, _ := cfg.Dump() // Password appears as "***"
fmt.Println(output)
// CONFIG_KEY   VALUE       SECRET
// ----------   -----       ------
// host         localhost
// password     ***         yes
```

To reveal secrets in development (use with care):

```go
output, _ := cfg.DumpWithOptions(&config.DumpOptions{
    Format:      "yaml",
    Content:     "config",
    MaskSecrets: false, // shows actual secret values
})
```

---

### Multiple dump formats

`Dump` / `DumpWithOptions` support four output formats: `json`, `yaml`, `text` (key=value), and `table`.

```go
type AppConfig struct {
    Host     string `yaml:"host"     env:"APP_HOST"`
    Port     int    `yaml:"port"     env:"APP_PORT"`
    Password string `yaml:"password" env:"APP_PASSWORD" secret:"true"`
}

cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{Host: "localhost", Port: 8080, Password: "s3cr3t"})
cfg.Load()

// Table (default via Dump)
tableOut, _ := cfg.Dump()

// JSON
jsonOut, _ := cfg.DumpWithOptions(&config.DumpOptions{Format: "json", Content: "config", MaskSecrets: true})

// Plain text key=value
textOut, _ := cfg.DumpWithOptions(&config.DumpOptions{Format: "text", Content: "config", MaskSecrets: true})
// host=localhost
// port=8080

// YAML
yamlOut, _ := cfg.DumpWithOptions(&config.DumpOptions{Format: "yaml", Content: "config", MaskSecrets: true})
// host: localhost
// port: 8080
// password: '***'
```

`DumpEnv` is a convenience wrapper that dumps env var names → values as `ENV_VAR=value` lines:

```go
envOut, _ := cfg.DumpEnv()
// APP_HOST=localhost
// APP_PORT=8080
// APP_PASSWORD=***
```

---

### Multiple dump content modes

The `Content` field of `DumpOptions` controls what information is included in the output.

| Value      | Shows |
|------------|-------|
| `"config"` | Config keys → current values |
| `"env"`    | Env var names → current values |
| `"metadata"` | Config key + env var name per field |
| `"all"`    | Everything above plus the Go struct field path |

```go
// Using the same AppConfig struct and cfg from the "Multiple dump formats" section above.

// Show env-variable-to-value mapping
cfg.DumpWithOptions(&config.DumpOptions{Format: "table", Content: "env", MaskSecrets: true})
// ENV_VAR      VALUE       SECRET
// -------      -----       ------
// APP_HOST     localhost
// APP_PORT     8080
// APP_PASSWORD ***         yes

// Full metadata for troubleshooting
cfg.DumpWithOptions(&config.DumpOptions{Format: "table", Content: "all", MaskSecrets: true})
// CONFIG_KEY   ENV_VAR      FIELD_PATH      VALUE       SECRET
// ----------   -------      ----------      -----       ------
// host         APP_HOST     Host            localhost
// port         APP_PORT     Port            8080
// password     APP_PASSWORD Password        ***         yes
```

---

### Getter/Setter API

Fields can be read and written at runtime by their **Go struct path** (`GetFieldValue` / `SetFieldValue`) or by their **config key** (`GetConfigValue` / `SetConfigValue`).

```go
type AppConfig struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{Host: "localhost", Port: 8080})
cfg.Load()

// Read by struct field path
host, _ := cfg.GetFieldValue("Host")
fmt.Println(host) // "localhost"

// Write by struct field path
cfg.SetFieldValue("Port", 9090)
fmt.Println(cfg.Data.Port) // 9090

// Read by config key
port, _ := cfg.GetConfigValue("port")
fmt.Println(port) // 9090

// Write by config key
cfg.SetConfigValue("host", "0.0.0.0")
fmt.Println(cfg.Data.Host) // "0.0.0.0"
```

For nested structs use dot-separated paths:

```go
type AppConfig struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
}

cfg.GetFieldValue("Server.Host")      // by struct path
cfg.GetConfigValue("server.host")     // by config key
cfg.SetFieldValue("Server.Port", 443)
cfg.SetConfigValue("server.port", 443)
```

---

### Nil pointer auto-initialization

When a `*struct` pointer field is `nil` and a path traversal (e.g. `SetFieldValue`) or env loading needs to pass through it, go-config automatically allocates the struct rather than returning an error.

```go
type Inner struct{ Value string }
type AppConfig struct {
    Inner *Inner `yaml:"inner"`
}

cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{}) // Inner is nil
cfg.SetFieldValue("Inner.Value", "hello")                     // Inner is allocated automatically
fmt.Println(cfg.Data.Inner.Value)                             // "hello"
```

---

### Slice env var parsing

Slice fields are populated from a single comma-separated environment variable.

```go
type AppConfig struct {
    Hosts []string `yaml:"hosts" env:"HOSTS"`
    Ports []int    `yaml:"ports" env:"PORTS"`
}

// export HOSTS=a.example.com,b.example.com,c.example.com
// export PORTS=8080,8081,8082
cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{})
cfg.Load()
fmt.Println(cfg.Data.Hosts) // [a.example.com b.example.com c.example.com]
fmt.Println(cfg.Data.Ports) // [8080 8081 8082]
```

---

### Map env var parsing

`map[string]T` fields are populated from a `KEY1=VAL1,KEY2=VAL2` formatted environment variable.

```go
type AppConfig struct {
    Labels map[string]string `yaml:"labels" env:"LABELS"`
}

// export LABELS=env=production,region=us-east-1,version=2
cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{})
cfg.Load()
fmt.Println(cfg.Data.Labels) // map[env:production region:us-east-1 version:2]
```

---

### `time.Duration` env support

Fields of type `time.Duration` are parsed with `time.ParseDuration`, accepting any string accepted by that function.

```go
type AppConfig struct {
    Timeout time.Duration `yaml:"timeout" env:"TIMEOUT"`
}

// export TIMEOUT=1h30m
cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{Timeout: 30 * time.Second})
cfg.Load()
fmt.Println(cfg.Data.Timeout) // 1h30m0s
```

---

### Extended boolean parsing

Boolean fields accept a wider range of values than Go's `strconv.ParseBool` when loaded from environment variables.

| Accepted values | Result |
|-----------------|--------|
| `true`, `t`, `yes`, `y`, `1`, `on` | `true` |
| `false`, `f`, `no`, `n`, `0`, `off` | `false` |

An empty string is also treated as `false`.

All values are case-insensitive.

```go
type AppConfig struct {
    Debug   bool `yaml:"debug"   env:"DEBUG"`
    Enabled bool `yaml:"enabled" env:"ENABLED"`
}

// export DEBUG=yes
// export ENABLED=on
cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{})
cfg.Load()
fmt.Println(cfg.Data.Debug, cfg.Data.Enabled) // true true
```

---

### Duplicate detection

Duplicate config keys or env var names across the struct tree are detected during `NewConfig` and returned as an error immediately.

```go
type Bad struct {
    HostA string `yaml:"host" env:"HOST"` // config key "host", env "HOST"
    HostB string `yaml:"host" env:"HOST"` // same config key AND same env var!
}

_, err := config.NewConfig(config.NewOptions(), &Bad{})
fmt.Println(err) // failed to list fields: duplicate config key "host" found
```

---

### Direct metadata access

`cfg.Metadata` is publicly exported and provides three lookup maps for advanced introspection or tooling.

```go
cfg, _ := config.NewConfig(config.NewOptions(), &AppConfig{})

// Look up by Go struct field path
info := cfg.Metadata.FieldPathMap["Server.Port"]
fmt.Println(info.ConfigKey) // "server.port"
fmt.Println(info.EnvVar)    // "APP_PORT"
fmt.Println(info.Secret)    // false

// Look up by env var name
info = cfg.Metadata.EnvMap["APP_PORT"]
fmt.Println(info.FieldPath) // "Server.Port"

// Look up by config key
info = cfg.Metadata.KeyMap["server.port"]
fmt.Println(info.EnvVar)    // "APP_PORT"

// Iterate all fields
for path, info := range cfg.Metadata.FieldPathMap {
    fmt.Printf("%s → config key: %s, env: %s\n", path, info.ConfigKey, info.EnvVar)
}
```
