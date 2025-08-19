# Debug & Dump Formats Example

This comprehensive example demonstrates all available dump formats, content types, and secret handling options. It's designed as both a learning tool and a debugging utility for the go-config library.

## Features Demonstrated

- **All Output Formats**: JSON, Text, Table
- **All Content Types**: Config, Environment, Metadata, All
- **Secret Handling**: Masked, visible, security audit
- **Interactive Mode**: Explore dump options dynamically
- **Complex Configuration**: Nested structures with various data types
- **Comprehensive Examples**: Every possible dump combination

## Configuration Structure

```go
type DebugConfig struct {
    App      AppSettings      // Application metadata and settings
    Database DatabaseSettings // Database connection with secrets
    External ExternalServices // Multiple external APIs with credentials
    Features FeatureFlags     // Boolean feature toggles
}
```

## Running the Example

### 1. Basic Demo (All Formats)
```bash
go run main.go
```

### 2. Interactive Mode
```bash
INTERACTIVE=true go run main.go
```

### 3. With Environment Overrides
```bash
export APP_NAME="Custom Debug App"
export DEBUG=true
export DATABASE_URL="postgres://user:secret@prod:5432/db"
export API_KEY="real-api-key"
go run main.go
```

### 4. Load from Configuration File
```bash
# Loads debug.json automatically if present
go run main.go
```

## Output Format Examples

### JSON Format
```json
{
  "app.debug": true,
  "app.environment": "testing",
  "app.log_level": "info",
  "app.name": "debug-showcase",
  "app.version": "2.0.0",
  "database.max_conns": 20,
  "database.timeout": "45s",
  "database.url": "***",
  "external.email_api.key": "***",
  "external.email_api.secret": "***",
  "external.payment_api.key": "***"
}
```

### Text Format
```
app.debug=true
app.environment=testing
app.log_level=info
app.name=debug-showcase
app.version=2.0.0
database.max_conns=20
database.timeout=45s
database.url=***
external.email_api.key=***
external.email_api.secret=***
```

### Table Format
```
CONFIG_KEY                    VALUE                   SECRET
----------                    -----                   ------
app.debug                     true                    
app.environment               testing                 
app.log_level                 info                    
app.name                      debug-showcase          
app.version                   2.0.0                   
database.max_conns            20                      
database.timeout              45s                     
database.url                  ***                     yes
external.email_api.key        ***                     yes
external.email_api.secret     ***                     yes
```

## Content Type Examples

### Config Content
Shows configuration keys and their values:
```json
{
  "app.name": "debug-showcase",
  "app.version": "2.0.0",
  "database.url": "***"
}
```

### Environment Content
Shows environment variable mappings:
```
APP_NAME=debug-showcase
APP_VERSION=2.0.0
DATABASE_URL=***
API_KEY=***
```

### Metadata Content  
Shows configuration keys, environment variables, and metadata:
```json
[
  {
    "config_key": "app.name",
    "env_var": "APP_NAME", 
    "value": "debug-showcase",
    "is_secret": false
  },
  {
    "config_key": "database.url",
    "env_var": "DATABASE_URL",
    "value": "***",
    "is_secret": true,
    "is_masked": true
  }
]
```

### All Content
Includes internal field paths for debugging:
```json
[
  {
    "config_key": "app.name",
    "env_var": "APP_NAME",
    "field_path": "App.Name",
    "value": "debug-showcase",
    "is_secret": false
  }
]
```

## Secret Handling Examples

### Default (Masked)
```
database.url=***
external.payment_api.key=***
external.payment_api.secret=***
```

### Show Secrets (Development Only)
```
database.url=postgres://testuser:testpass@test-db:5432/test_db
external.payment_api.key=pk_test_updated_key
external.payment_api.secret=sk_test_updated_secret
```

### Security Audit
```
CONFIG_KEY                    ENV_VAR                 SECRET
----------                    -------                 ------
database.url                  DATABASE_URL            yes
external.email_api.key        API_KEY                 yes
external.email_api.secret     API_SECRET              yes
external.payment_api.key      API_KEY                 yes
external.payment_api.secret   API_SECRET              yes
external.storage_api.key      API_KEY                 yes
external.storage_api.secret   API_SECRET              yes
```

## Interactive Commands

When running with `INTERACTIVE=true`:

- **formats** - Demonstrate all output formats (JSON, text, table)
- **content** - Show all content types (config, env, metadata, all)
- **secrets** - Show secret handling options (masked, visible, audit)
- **custom** - Create custom dump with your own options
- **quit** - Exit interactive mode

## Use Cases

### Development
- **Configuration Debugging**: See exactly what values are loaded
- **Format Testing**: Try different output formats for your needs
- **Secret Verification**: Ensure secrets are properly marked and masked
- **Environment Testing**: Test environment variable overrides

### Operations
- **Configuration Auditing**: Review all configuration in production
- **Secret Detection**: Identify which fields contain sensitive data
- **Documentation**: Generate configuration documentation
- **Deployment Verification**: Verify configuration before deployment

### Learning
- **Feature Exploration**: Learn all dump capabilities
- **Format Comparison**: Compare different output formats
- **Interactive Testing**: Experiment with different options
- **Best Practices**: See recommended patterns for different scenarios

## Key Learning Points

1. **Format Selection**: Choose the right format for your use case
   - JSON: Structured data, API responses, documentation
   - Text: Shell exports, simple key-value pairs
   - Table: Human-readable, operational overviews

2. **Content Filtering**: Control what information is included
   - Config: Application configuration
   - Env: Environment variable mapping
   - Metadata: Complete field information
   - All: Debug-level information with internal paths

3. **Secret Security**: Always protect sensitive information
   - Default masking prevents accidental exposure
   - Development mode for debugging only
   - Security audit to identify sensitive fields

4. **Interactive Exploration**: Use interactive mode to learn and experiment
   - Safe environment to try different options
   - Real-time feedback on configuration changes
   - Custom dump creation for specific needs