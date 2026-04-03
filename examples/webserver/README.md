# Web Server Example

Shows a realistic nested config for a web server with TLS, CORS, logging, and
security sections. Optionally starts an HTTP server that serves its own config
dump at `/config`.

## Config struct (abbreviated)

```go
type ServerConfig struct {
    Server   HTTPServerConfig `json:"server"`
    TLS      TLSConfig        `json:"tls"`
    CORS     CORSConfig       `json:"cors"`
    Logging  LoggingConfig    `json:"logging"`
    Security SecurityConfig   `json:"security"`
}
```

Secret fields (`secret:"true"`): `security.jwt_secret`, `tls.key_file`.

## Running

```bash
go run main.go
```

With environment overrides:

```bash
SERVER_HOST="0.0.0.0" SERVER_PORT=3000 LOG_LEVEL="debug" go run main.go
```

Start the live HTTP server:

```bash
START_SERVER=true go run main.go
# visit http://localhost:8080/ or http://localhost:8080/config
```

Show debug metadata:

```bash
DEBUG=true go run main.go
```

## Expected output (excerpt)

Environment variables dump (`cfg.DumpEnv()`):

```
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_ORIGINS=https://myapp.com,https://admin.myapp.com
CORS_ENABLED=true
JWT_SECRET=***
LOG_FORMAT=json
LOG_LEVEL=info
LOG_OUTPUT=stdout
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
TLS_ENABLED=false
TLS_KEY_FILE=***
TLS_PORT=8443
...
```

JSON config dump (`cfg.DumpWithOptions` with `Format: "json"`, `Content: "config"`):

```json
{
  "cors": {
    "allowed_origins": ["https://myapp.com", "https://admin.myapp.com"],
    "enabled": true
  },
  "logging": {
    "format": "json",
    "level": "info",
    "output": "stdout"
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "tls": {
    "enabled": false,
    "port": 8443
  }
}
```