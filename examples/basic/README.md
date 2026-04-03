# Basic Example

Loads a flat config struct from a JSON file, applies environment variable
overrides, and dumps the result in several formats with secret masking.

## Config struct

```go
type AppConfig struct {
    AppName  string `json:"app_name" env:"APP_NAME"`
    Version  string `json:"version"  env:"VERSION"`
    Debug    bool   `json:"debug"    env:"DEBUG"`
    Port     int    `json:"port"     env:"PORT"`
    APIKey   string `json:"api_key"  env:"API_KEY" secret:"true"`
    LogLevel string `json:"log_level" env:"LOG_LEVEL"`
}
```

## Running

```bash
go run main.go
```

With environment overrides:

```bash
APP_NAME="custom" PORT=3000 API_KEY="my-key" go run main.go
```

Show unmasked secrets (development only):

```bash
SHOW_SECRETS=true go run main.go
```

## Expected output

Default dump (table format, secrets masked via `cfg.Dump()`):

```
CONFIG_KEY   VALUE           SECRET
----------   -----           ------
api_key      ***             yes
app_name     basic-example
debug        true
log_level    debug
port         9000
version      1.2.0
```

Environment variables dump (`cfg.DumpEnv()`):

```
API_KEY=***
APP_NAME=basic-example
DEBUG=true
LOG_LEVEL=debug
PORT=9000
VERSION=1.2.0
```