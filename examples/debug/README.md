# Debug Example

Demonstrates every dump format (`json`, `text`, `table`) and content type
(`config`, `env`, `metadata`, `all`). Includes an interactive mode for
exploring options at runtime.

## Config struct (abbreviated)

```go
type DebugConfig struct {
    App      AppSettings      `json:"app"`
    Database DatabaseSettings `json:"database"`
    External ExternalServices `json:"external"`
    Features FeatureFlags     `json:"features"`
}
```

Fields marked `secret:"true"`: database URL, API keys and secrets.

## Running

```bash
go run main.go
```

Interactive mode (prompts for format/content/secret options):

```bash
INTERACTIVE=true go run main.go
```

With environment overrides:

```bash
APP_NAME="custom" DEBUG=true go run main.go
```

## Expected output (excerpt)

JSON format, config content:

```json
{
  "app": {
    "debug": false,
    "environment": "testing",
    "name": "debug-showcase",
    "version": "2.0.0"
  },
  "database": {
    "max_conns": 20,
    "timeout": "45s",
    "url": "postgres://testuser:testpass@test-db:5432/test_db"
  }
}
```

Text format, config content:

```
app.debug=false
app.environment=testing
app.name=debug-showcase
app.version=2.0.0
database.max_conns=20
database.timeout=45s
database.url=postgres://testuser:testpass@test-db:5432/test_db
```

Table format with secret masking (`MaskSecrets: true`):

```
CONFIG_KEY              VALUE   SECRET
----------              -----   ------
database.url            ***     yes
external.email_api.key  ***     yes
...
```

> **Note:** When `MaskSecrets` is not explicitly set in `DumpOptions` it
> defaults to `false`, so secrets are visible. Use `cfg.Dump()` (which sets
> `MaskSecrets: true` internally) or pass `MaskSecrets: true` explicitly.