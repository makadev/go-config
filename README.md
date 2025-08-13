# Go Config Package

A flexible and powerful configuration loading package for Go applications that supports struct tag annotations, multiple file formats, environment variable overrides, and nested configurations.

## Features

- ✅ **Struct Tag Annotations**: Use `config`, `env`, `default`, `required`, and `secret` tags
- ✅ **Multiple File Formats**: Support for YAML and JSON configuration files
- ✅ **Priority Loading**: `config.local.yaml` → `config.local.json` → `config.yaml` → `config.json`
- ✅ **Environment Variable Overrides**: Override any configuration value via environment variables
- ✅ **Nested Structures**: Full support for nested configuration structures
- ✅ **Type Safety**: Automatic type conversion with validation
- ✅ **Default Values**: Set default values directly in struct tags
- ✅ **Validation**: Required field validation and custom validation options
- ✅ **Secret Management**: Mark sensitive fields with `secret:"true"` for redaction
- ✅ **Configuration Dumping**: Export configuration to YAML, JSON, env, or flat formats
- ✅ **Generic API**: Type-safe `Config[T]` struct with methods for advanced usage
- ✅ **Flexible API**: Simple `Load()` function with advanced `LoadWithOptions()` and `NewConfig[T]()`

## Installation

```bash
go get github.com/matthias/go-config
```

## Quick Start

### 1. Define Your Configuration Struct

```go
type Config struct {
    // Simple fields with defaults and environment variable support
    AppName string `config:"app_name" env:"APP_NAME" default:"MyApp"`
    Debug   bool   `config:"debug" env:"DEBUG" default:"false"`
    Port    int    `config:"port" env:"PORT" default:"8080"`
    
    // Nested configuration
    Database DatabaseConfig `config:"database"`
    
    // Required fields
    APIKey string `config:"api_key" env:"API_KEY" required:"true"`
    
    // Secret fields (will be redacted in dumps)
    Secret string `config:"secret" env:"SECRET" secret:"true"`
}

type DatabaseConfig struct {
    Host     string `config:"host" env:"DB_HOST" default:"localhost"`
    Port     int    `config:"port" env:"DB_PORT" default:"5432"`
    Username string `config:"username" env:"DB_USER" required:"true"`
    Password string `config:"password" env:"DB_PASSWORD" required:"true" secret:"true"`
}
```

### 2. Load Configuration

```go
package main

import (
    "log"
    "github.com/matthias/go-config"
)

func main() {
    var cfg Config
    
    if err := config.Load(&cfg); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Use your configuration
    fmt.Printf("Starting %s on port %d\n", cfg.AppName, cfg.Port)
}
```

### 3. Create Configuration Files

**config.yaml** (base configuration):
```yaml
app_name: "ProductionApp"
debug: false
port: 3000

database:
  host: "prod-db.example.com"
  username: "prod_user"
  password: "prod_password"
```

**config.local.yaml** (local overrides, highest priority):
```yaml
app_name: "DevApp"
debug: true
port: 8080

database:
  host: "localhost"
  username: "dev_user"
  password: "dev_password"
```

### 4. Environment Variable Overrides

```bash
export APP_NAME="EnvApp"
export DEBUG="true"
export DB_HOST="env-db.example.com"
export API_KEY="secret-api-key"
```

## Struct Tag Reference

| Tag | Description | Example |
|-----|-------------|---------|
| `config` | Key name in configuration files | `config:"server_port"` |
| `env` | Environment variable name | `env:"SERVER_PORT"` |
| `default` | Default value if not set | `default:"8080"` |
| `required` | Mark field as required | `required:"true"` |
| `secret` | Mark field as secret (redacted in dumps) | `secret:"true"` |

## Supported Types

- **Basic types**: `string`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`
- **Time**: `time.Duration` (parsed from strings like "5m", "30s", "1h")
- **Collections**: `[]string`, `[]int`, etc. (comma-separated in env vars)
- **Maps**: `map[string]string`, `map[string]int`, etc. (format: "key1=value1,key2=value2")
- **Nested structs**: Unlimited nesting depth
- **Pointers**: Pointers to any supported type

## Advanced Usage

### Custom Configuration Paths

```go
options := &config.Options{
    ConfigPaths: []string{"app.yaml", "app.json"},
    EnvPrefix:   "MYAPP",
    Secret:      true,              // Enable secret redaction in dumps
    SecretWith:  "[HIDDEN]",        // Custom redaction text
}

err := config.LoadWithOptions(&cfg, options)
```

### Skip File or Environment Loading

```go
options := &config.Options{
    SkipFiles: true,  // Only use defaults and environment variables
    SkipEnv:   false, // Only use defaults and configuration files
}
```

### Generic Config API

```go
// Create a generic config instance
config, err := config.NewConfig[MyConfig](nil)
if err != nil {
    log.Fatal(err)
}

// Load configuration
if err := config.Load(); err != nil {
    log.Fatal(err)
}

// Access the loaded data
fmt.Printf("Host: %s\n", config.Data.Database.Host)

// Get field metadata
info, err := config.GetFieldInfo("Database.Host")
fmt.Printf("Field: %s, Env: %s, Default: %s\n", 
    info.FieldPath, info.EnvVar, info.DefaultValue)

// Get/Set field values by path
value, err := config.GetFieldValue("Database.Host")
err = config.SetFieldValue("Database.Host", "new-host.com")

// Create redacted copy for logging/debugging
redacted, err := config.RedactedCopy()
fmt.Printf("Config (redacted): %+v\n", redacted.Data)

// Dump configuration in various formats
yamlDump, err := config.Dump("yaml")
jsonDump, err := config.Dump("json")
envDump, err := config.Dump("env")
flatDump, err := config.Dump("flat")
```

### Configuration Dumping

```go
var cfg Config
config.Load(&cfg)

// Dump to YAML format
yamlOutput, err := config.Dump(&cfg, metadata, "yaml", true, "[REDACTED]", "")

// Dump to JSON format  
jsonOutput, err := config.Dump(&cfg, metadata, "json", true, "[REDACTED]", "")

// Dump to .env format (only fields with env tags)
envOutput, err := config.Dump(&cfg, metadata, "env", false, "", "MYAPP_")

// Dump to flat format
flatOutput, err := config.Dump(&cfg, metadata, "flat", false, "", "")
```

### Field Introspection

```go
// Get all configurable fields metadata
metadata, err := config.GetFieldInfoMap(&cfg)
for path, info := range metadata {
    fmt.Printf("Field: %s, Env: %s, Default: %s, Secret: %t\n", 
        path, info.EnvVar, info.DefaultValue, info.Secret)
}

// Get field value by path
value, err := config.GetFieldValue(&cfg, "Database.Host")

// Set field value by path
err = config.SetFieldValue(&cfg, "Database.Host", "new-host.com")

// Find available config files
configFile, err := config.FindConfigFile([]string{"config.yaml", "config.json"})

// Get config file format
format, err := config.GetConfigFormat("config.yaml") // returns "yaml"
```

## Configuration Loading Priority

The configuration loading follows this priority order (later sources override earlier ones):

1. **Default values** (from struct tags)
2. **Configuration files** (first found file is used):
   - `config.local.yaml`
   - `config.local.json`
   - `config.yaml`
   - `config.json`
3. **Environment variables** (highest priority)

## Environment Variable Formats

### Boolean Values
Accepted values for boolean fields:
- **True**: `true`, `t`, `yes`, `y`, `1`, `on` (case-insensitive)
- **False**: `false`, `f`, `no`, `n`, `0`, `off`, `` (empty string, case-insensitive)

### Duration Values
Use Go duration format: `5m`, `30s`, `1h30m`, `100ms`

### Slice Values
Comma-separated values: `value1,value2,value3`

### Map Values
Key-value pairs: `key1=value1,key2=value2,key3=value3`

## Error Handling

The package provides detailed error messages for common issues:

- Invalid struct input (not a pointer to struct)
- Type conversion errors
- Missing required fields
- Invalid file formats
- Field path resolution errors

## 💡 Usage

```go
var cfg AppConfig
err := config.Load(&cfg)  // Simple!

// Or with advanced options:
opts := &config.Options{
    ConfigPaths: []string{"custom.yaml"},
    EnvPrefix:   "MYAPP",
    Secret:      true,
    SecretWith:  "[REDACTED]",
}
err := config.LoadWithOptions(&cfg, opts)
```

## 🔧 Environment Variable Examples

```bash
# Override server configuration
export SERVER_HOST="production.example.com"
export SERVER_PORT="443"

# Override nested database settings  
export DB_HOST="db.prod.com"
export DB_PASSWORD="secure_password"

# Set slice values (comma-separated)
export ALLOWED_IPS="192.168.1.0/24,10.0.0.0/8"

# Set map values (key=value pairs)
export METADATA="env=prod,region=eu-west,version=2.0"

# Boolean values (flexible format)
export DEBUG="yes"           # true
export ENABLE_TLS="1"        # true  
export PROFILING="off"       # false
```

## 🛠 Complete Implementation Steps

To use this package in your project:

1. **Create the package structure:**
```bash
mkdir config
cd config
```

2. **Add the dependencies:**
```bash
go mod init your-module/config
go get gopkg.in/yaml.v3
```

3. **Copy the implementation files:**
   - Copy all the Go files I provided above
   - Update the module path in `go.mod` to match your project

4. **Create your configuration struct:**
```go
type MyAppConfig struct {
    Server   ServerConfig `config:"server"`
    Database DBConfig     `config:"database"`  
    APIKey   string       `config:"api_key" env:"API_KEY" required:"true"`
}
```

5. **Use in your application:**
```go
var cfg MyAppConfig
if err := config.Load(&cfg); err != nil {
    log.Fatal(err)
}
```

## 🧪 Testing

The package includes comprehensive tests covering:
- Default value application
- Environment variable overrides  
- File loading priority
- Nested struct handling
- Type conversion edge cases
- Error conditions
- Performance benchmarks

Run tests with:
```bash
go test -v
go test -race
go test -bench=.
```

## Dependencies

- `gopkg.in/yaml.v3` - YAML parsing
- Standard library only for JSON parsing

## Examples

See the `example/` directory for a complete working example demonstrating all features.

## License

MIT License
