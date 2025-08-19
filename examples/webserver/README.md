# Web Server Configuration Example

This example demonstrates a real-world web server configuration with nested structures, security settings, and multiple dump formats for operations.

## Features Demonstrated

- **Nested Configuration**: Server, TLS, CORS, Logging, Security sections
- **Security Fields**: JWT secrets and TLS keys marked as secret
- **Time Durations**: Proper handling of timeout configurations
- **Array Fields**: CORS origins, methods, headers
- **Operational Dumps**: Different formats for different use cases
- **Live Server**: Actually runs an HTTP server with the configuration

## Configuration Structure

```go
type ServerConfig struct {
    Server   HTTPServerConfig // Host, port, timeouts
    TLS      TLSConfig        // TLS/SSL configuration
    CORS     CORSConfig       // Cross-origin resource sharing
    Logging  LoggingConfig    // Logging configuration
    Security SecurityConfig   // Security and rate limiting
}
```

## Running the Example

### 1. Configuration Validation and Dumps
```bash
go run main.go
```

### 2. With Environment Overrides
```bash
export SERVER_HOST="0.0.0.0"
export SERVER_PORT=3000
export LOG_LEVEL="debug"
export JWT_SECRET="my-custom-secret"
go run main.go
```

### 3. Start Actual Server
```bash
START_SERVER=true go run main.go
```

### 4. Debug Mode (show metadata)
```bash
DEBUG=true go run main.go
```

## Example Outputs

### Operations Table View
```
CONFIG_KEY              VALUE                          SECRET
----------              -----                          ------
cors.allowed_headers    [Content-Type Authorization]   
cors.allowed_methods    [GET POST PUT DELETE OPTIONS]  
cors.allowed_origins    [https://myapp.com]           
cors.enabled            true                           
logging.format          json                           
logging.level           info                           
security.jwt_secret     ***                           yes
server.host             0.0.0.0                        
server.port             8080                           
tls.cert_file           /etc/ssl/certs/server.crt      
tls.enabled             true                           
tls.key_file            ***                           yes
```

### Environment Variables for Deployment
```
CORS_ALLOWED_HEADERS=[Content-Type Authorization X-Requested-With]
CORS_ALLOWED_METHODS=[GET POST PUT DELETE OPTIONS]
CORS_ALLOWED_ORIGINS=[https://myapp.com https://admin.myapp.com]
CORS_ENABLED=true
JWT_SECRET=***
LOG_FORMAT=json
LOG_LEVEL=info
LOG_OUTPUT=/var/log/server.log
MAX_REQUEST_SIZE=52428800
RATE_LIMIT_RPS=200
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SESSION_TIMEOUT=24h0m0s
TLS_CERT_FILE=/etc/ssl/certs/server.crt
TLS_ENABLED=true
TLS_KEY_FILE=***
TLS_PORT=8443
```

### Complete JSON Configuration
```json
{
  "cors.allowed_headers": ["Content-Type", "Authorization", "X-Requested-With"],
  "cors.allowed_methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
  "cors.allowed_origins": ["https://myapp.com", "https://admin.myapp.com"],
  "cors.enabled": true,
  "logging.format": "json",
  "logging.level": "info",
  "logging.output": "/var/log/server.log",
  "security.jwt_secret": "***",
  "security.max_request_size": 52428800,
  "security.rate_limit_rps": 200,
  "security.session_timeout": "24h0m0s",
  "server.host": "0.0.0.0",
  "server.idle_timeout": "2m0s",
  "server.port": 8080,
  "server.read_timeout": "30s",
  "server.write_timeout": "30s",
  "tls.cert_file": "/etc/ssl/certs/server.crt",
  "tls.enabled": true,
  "tls.key_file": "***",
  "tls.port": 8443
}
```

## Server Endpoints

When running with `START_SERVER=true`:

- **http://localhost:8080/** - Server information page
- **http://localhost:8080/config** - Live configuration dump (JSON)

## Use Cases

### Development
- **Quick validation**: Check configuration without starting services
- **Debug mode**: See all metadata and field mappings
- **Override testing**: Test environment variable overrides

### Operations
- **Deployment verification**: Table view for quick config review
- **Environment export**: Generate environment variables for containers
- **Documentation**: JSON dump for configuration documentation

### Production
- **Health checks**: `/config` endpoint for monitoring
- **Graceful shutdown**: Proper signal handling
- **Secret masking**: Never expose secrets in logs or dumps

## Key Features

1. **Nested Structures**: Demonstrates complex configuration hierarchies
2. **Type Safety**: Proper Go types for all configuration values
3. **Security**: Secret fields are properly masked
4. **Validation**: Custom validation logic before server start
5. **Live Configuration**: Running server serves its own configuration
6. **Multiple Formats**: Different dumps for different operational needs