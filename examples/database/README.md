# Database Example

Demonstrates a complex, deeply nested config with multiple database connections,
cache settings, migrations, and backups. Shows the struct-as-env-prefix pattern
(`env:"PRIMARY_"`) for scoping environment variables to nested structs.

## Config struct (abbreviated)

```go
type DatabaseConfig struct {
    Primary   DatabaseConnection `json:"primary"   env:"PRIMARY_"`
    ReadOnly  DatabaseConnection `json:"read_only" env:"READ_ONLY_"`
    Cache     CacheConfig        `json:"cache"`
    Migration MigrationConfig    `json:"migration"`
    Backup    BackupConfig       `json:"backup"`
}
```

Secret fields (`secret:"true"`): database passwords, Redis password, S3 secret
key.

## Running

```bash
go run main.go
```

With environment overrides:

```bash
PRIMARY_DB_HOST="prod-db" PRIMARY_DB_PASSWORD="secret" go run main.go
```

Show unmasked secrets (development only):

```bash
SHOW_SECRETS=true go run main.go
```

## Expected output (excerpt)

Environment variables dump (`cfg.DumpEnv()`):

```
PRIMARY_DB_DRIVER=postgres
PRIMARY_DB_HOST=prod-db.company.com
PRIMARY_DB_NAME=myapp_production
PRIMARY_DB_PASSWORD=***
PRIMARY_DB_PORT=5432
PRIMARY_DB_SSL_MODE=require
PRIMARY_DB_USERNAME=app_user
READ_ONLY_DB_HOST=read-replica.company.com
READ_ONLY_DB_PASSWORD=***
REDIS_HOST=redis.company.com
REDIS_PASSWORD=***
S3_SECRET_ACCESS_KEY=***
...
```

Connection information printed at the end:

```
Primary Database:
  Type: postgres
  Host: prod-db.company.com:5432
  Database: myapp_production
  Username: app_user
  SSL Mode: require
  Max Connections: 100 open, 25 idle
```