# go-config Examples

This directory contains comprehensive examples demonstrating the load and dump functionality of the go-config library. Each example showcases different use cases, configuration patterns, and operational scenarios.

## Available Examples

### 🚀 [Basic Example](./basic/)
**Perfect for getting started**
- Simple configuration structure
- File loading and environment overrides
- Basic dump formats (JSON, text, table)
- Secret masking demonstration

**Use cases**: First-time users, simple applications, learning the basics

### 🌐 [Web Server Example](./webserver/)
**Real-world web server configuration**
- Nested configuration structures (Server, TLS, CORS, Logging, Security)
- Complex field types (durations, arrays, nested objects)
- Operational dump formats for different audiences
- Live HTTP server with configuration endpoints

**Use cases**: Web applications, API servers, microservices

### 🗄️ [Database Example](./database/)
**Comprehensive database and caching setup**
- Multiple database connections (primary, read-only)
- Secret management (passwords, API keys)
- Complex nested structures (cache, migration, backup)
- Security auditing and operational views

**Use cases**: Data-heavy applications, enterprise systems, multi-database setups

### 🔧 [Debug Example](./debug/)
**Complete dump functionality showcase**
- All output formats (JSON, text, table)
- All content types (config, env, metadata, all)
- Interactive exploration mode
- Secret handling demonstrations

**Use cases**: Learning, debugging, configuration exploration, development tools

## Quick Start

### 1. Clone and Navigate
```bash
git clone https://github.com/makadev/go-config
cd go-config/examples
```

### 2. Try the Basic Example
```bash
cd basic
go run main.go
```

### 3. Experiment with Environment Variables
```bash
export APP_NAME="My Custom App"
export DEBUG=true
export PORT=3000
go run main.go
```

### 4. Explore Other Examples
```bash
# Web server with live endpoints
cd ../webserver
START_SERVER=true go run main.go

# Database configuration with secrets
cd ../database
SHOW_SECRETS=true go run main.go

# Interactive debug mode
cd ../debug
INTERACTIVE=true go run main.go
```

## Common Patterns Demonstrated

### Configuration Loading
```go
// Initialize with defaults
cfg, err := config.NewConfig(nil, defaultConfig)

// Load from file (optional)
cfg.LoadFromFile("config.json")

// Load from environment (overrides file)
cfg.LoadFromEnv()
```

### Dump Formats
```go
// JSON dump (default)
cfg.DumpConfig()

// Environment variables
cfg.DumpEnv()

// Custom options
cfg.Dump(config.DumpOptions{
    Format:  "table",
    Content: "metadata",
    ShowSecrets: false,
})
```

### Secret Handling
```go
type Config struct {
    PublicKey  string `json:"public_key" env:"PUBLIC_KEY"`
    PrivateKey string `json:"private_key" env:"PRIVATE_KEY" secret:"true"`
}
```

## Output Format Guide

### JSON Format
**Best for**: API responses, documentation, structured data
```json
{
  "app.name": "my-app",
  "app.port": 8080,
  "database.password": "***"
}
```

### Text Format  
**Best for**: Shell exports, simple key-value pairs
```
APP_NAME=my-app
APP_PORT=8080
DATABASE_PASSWORD=***
```

### Table Format
**Best for**: Human-readable overviews, operations
```
CONFIG_KEY          VALUE    SECRET
----------          -----    ------
app.name            my-app   
app.port            8080     
database.password   ***      yes
```

## Content Type Guide

| Content Type | Description | Use Case |
|--------------|-------------|----------|
| `config` | Configuration keys and values | Application configuration |
| `env` | Environment variable mappings | Container deployment |
| `metadata` | Complete field information | Development, debugging |
| `all` | Everything including internal paths | Deep debugging |

## Security Best Practices

### ✅ DO
- Mark sensitive fields with `secret:"true"`
- Use masked dumps in production logs
- Keep `ShowSecrets: true` for development only
- Validate configuration before using secrets

### ❌ DON'T
- Log configuration with `ShowSecrets: true` in production
- Commit configuration files with real secrets
- Expose configuration dumps in public APIs without masking
- Use default passwords in production

## Environment Variable Patterns

### Single Service
```bash
export APP_NAME="my-service"
export APP_PORT=8080
export DATABASE_URL="postgres://..."
```

### Microservices
```bash
export SERVICE_NAME="user-service"
export DATABASE_URL="postgres://..."
export CACHE_URL="redis://..."
export API_GATEWAY_URL="https://..."
```

### Docker Compose
```yaml
environment:
  - APP_NAME=my-service
  - DATABASE_URL=postgres://db:5432/myapp
  - REDIS_URL=redis://redis:6379
```

### Kubernetes
```yaml
env:
- name: APP_NAME
  value: "my-service"
- name: DATABASE_URL
  valueFrom:
    secretKeyRef:
      name: db-secret
      key: url
```

## Use Case Matrix

| Example | Simple Setup | Complex Config | Secrets | Live Demo | Interactive |
|---------|-------------|----------------|---------|-----------|-------------|
| Basic | ✅ | ❌ | ✅ | ❌ | ❌ |
| Web Server | ❌ | ✅ | ✅ | ✅ | ❌ |
| Database | ❌ | ✅ | ✅ | ❌ | ❌ |
| Debug | ❌ | ✅ | ✅ | ❌ | ✅ |

## Development Workflow

### 1. Start with Basic
Learn the fundamentals with the basic example

### 2. Choose Your Domain
Pick the example closest to your use case:
- Web applications → Web Server
- Data applications → Database  
- Learning/debugging → Debug

### 3. Adapt and Extend
Copy example patterns to your own configuration

### 4. Validate and Deploy
Use dump functions to verify configuration before deployment

## Contributing

Found an issue or want to add an example? 

1. Check existing examples for similar patterns
2. Follow the established structure (main.go, config file, README.md)
3. Include comprehensive documentation
4. Test with various environment combinations
5. Submit a pull request

## Need Help?

- 📖 **Documentation**: Check individual example READMEs
- 🔧 **Debugging**: Use the debug example for interactive exploration  
- 💡 **Patterns**: Look at similar examples for configuration patterns
- 🚀 **Getting Started**: Begin with the basic example

Each example is self-contained and thoroughly documented. Start with the basic example and progress to more complex scenarios as needed.