# Basic Configuration Example

This example demonstrates the fundamental usage of go-config for loading and dumping configuration data.

## Features Demonstrated

- Loading configuration from JSON files
- Environment variable overrides
- Secret field handling
- Multiple dump formats (JSON, text, table)
- Development vs production secret handling

## Configuration Structure

```go
type AppConfig struct {
    AppName    string `json:"app_name" env:"APP_NAME"`
    Version    string `json:"version" env:"VERSION"`
    Debug      bool   `json:"debug" env:"DEBUG"`
    Port       int    `json:"port" env:"PORT"`
    APIKey     string `json:"api_key" env:"API_KEY" secret:"true"`
    LogLevel   string `json:"log_level" env:"LOG_LEVEL"`
}
```

## Running the Example

### 1. Basic Run (with file config)
```bash
go run main.go
```

### 2. With Environment Variables
```bash
export APP_NAME="my-custom-app"
export DEBUG=true
export PORT=3000
export API_KEY="env-secret-key"
go run main.go
```

### 3. Show Secrets (development only)
```bash
SHOW_SECRETS=true go run main.go
```

## Expected Output

### JSON Config Dump
```json
{
  "api_key": "***",
  "app_name": "basic-example",
  "debug": true,
  "log_level": "debug",
  "port": 9000,
  "version": "1.2.0"
}
```

### Environment Variables Dump
```
API_KEY=***
APP_NAME=basic-example
DEBUG=true
LOG_LEVEL=debug
PORT=9000
VERSION=1.2.0
```

### Table Format Dump
```
CONFIG_KEY	VALUE	SECRET
----------	-----	------
api_key	***	yes
app_name	basic-example	
debug	true	
log_level	debug	
port	9000	
version	1.2.0	
```

## Key Learning Points

1. **Configuration Priority**: Environment variables override file values
2. **Secret Handling**: Fields marked with `secret:"true"` are masked in dumps
3. **Multiple Formats**: Choose the right dump format for your use case
4. **Development vs Production**: Use `SHOW_SECRETS=true` only in development
5. **Flexibility**: Easy to switch between file-based and environment-based configuration

## Use Cases

- Application initialization and validation
- Configuration debugging and troubleshooting  
- Deployment verification
- Development environment setup
- Documentation and examples