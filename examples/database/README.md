# Database Configuration Example

This example demonstrates a comprehensive database configuration with multiple connections, caching, migrations, and backups. It showcases advanced secret handling and operational configuration management.

## Features Demonstrated

- **Multiple Database Connections**: Primary and read-only replicas
- **Secret Management**: Passwords, API keys properly marked and masked
- **Complex Nested Structures**: Cache, migration, backup configurations
- **Validation Logic**: Comprehensive configuration validation
- **Operational Dumps**: Different views for different operational needs
- **Security Auditing**: Identify which fields contain sensitive information

## Configuration Structure

```go
type DatabaseConfig struct {
    Primary   DatabaseConnection // Main database connection
    ReadOnly  DatabaseConnection // Read replica connection  
    Cache     CacheConfig        // Redis/Memcache configuration
    Migration MigrationConfig    // Database migration settings
    Backup    BackupConfig       // Backup and S3 configuration
}
```

### Security Fields
All sensitive fields are marked with `secret:"true"`:
- Database passwords
- Redis authentication
- S3 secret access keys

## Running the Example

### 1. Basic Configuration Validation
```bash
go run main.go
```

### 2. With Environment Overrides
```bash
export DB_HOST="production-db.company.com"
export DB_PASSWORD="production-secret"
export REDIS_PASSWORD="redis-auth-token"
export S3_SECRET_ACCESS_KEY="real-s3-secret"
go run main.go
```

### 3. Show Secrets (Development Only)
```bash
SHOW_SECRETS=true go run main.go
```

## Example Outputs

### Security Audit Table
```
CONFIG_KEY                    ENV_VAR                 VALUE                    SECRET
----------                    -------                 -----                    ------
backup.s3.access_key_id      S3_ACCESS_KEY_ID        AKIAIOSFODNN7EXAMPLE     
backup.s3.secret_access_key  S3_SECRET_ACCESS_KEY    ***                      yes
cache.redis.password         REDIS_PASSWORD          ***                      yes
primary.password             DB_PASSWORD             ***                      yes
read_only.password           DB_PASSWORD             ***                      yes
```

### Environment Variables for Deployment
```
BACKUP_ENABLED=true
BACKUP_RETENTION=2160h0m0s
BACKUP_SCHEDULE=0 3 * * *
CACHE_ENABLED=true
CACHE_TYPE=redis
DB_CONNECT_TIMEOUT=30s
DB_CONN_MAX_LIFETIME=1h0m0s
DB_DRIVER=postgres
DB_HOST=prod-db.company.com
DB_MAX_IDLE_CONNS=25
DB_MAX_OPEN_CONNS=100
DB_NAME=myapp_production
DB_PASSWORD=***
DB_PORT=5432
DB_SSL_MODE=require
DB_USERNAME=app_user
MIGRATION_DIR=/app/migrations
MIGRATION_ENABLED=true
MIGRATION_LOCK_KEY=migration_lock_prod
MIGRATION_LOCK_TIMEOUT=15m0s
MIGRATION_TABLE=schema_migrations
REDIS_DB=1
REDIS_HOST=redis.company.com
REDIS_PASSWORD=***
REDIS_POOL_SIZE=20
REDIS_PORT=6379
S3_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
S3_BUCKET=myapp-prod-backups
S3_REGION=us-west-2
S3_SECRET_ACCESS_KEY=***
```

### Database Connection Information
```
Primary Database:
  Type: postgres
  Host: prod-db.company.com:5432
  Database: myapp_production
  Username: app_user
  SSL Mode: require
  Max Connections: 100 open, 25 idle

Read-Only Database:
  Type: postgres
  Host: read-replica.company.com:5432
  Database: myapp_production
  Username: readonly_user
  SSL Mode: require
  Max Connections: 50 open, 15 idle

Cache: redis
  Redis: redis.company.com:6379 (DB 1)

Migrations: Enabled
  Directory: /app/migrations
  Table: schema_migrations

Backup: Enabled
  Schedule: 0 3 * * *
  S3 Bucket: myapp-prod-backups (us-west-2)
  Retention: 2160h0m0s
```

## Use Cases

### Development
- **Local Setup**: Override database hosts for local development
- **Debug Configuration**: See all values including secrets
- **Validation**: Ensure configuration is complete before starting services

### Production Deployment
- **Security Audit**: Identify all secret fields before deployment
- **Environment Generation**: Create container environment variables
- **Deployment Verification**: Verify configuration without exposing secrets
- **Connection Testing**: Test database connectivity with masked passwords

### Operations
- **Health Monitoring**: Check database connection settings
- **Backup Verification**: Ensure backup configuration is correct
- **Performance Tuning**: Review connection pool settings
- **Migration Management**: Verify migration configuration

## Key Features

1. **Multiple Database Types**: Primary/replica pattern with different settings
2. **Comprehensive Caching**: Redis and Memcache support with TTL configuration
3. **Migration Management**: Database migration settings with locking
4. **Backup Integration**: S3-based backup configuration
5. **Secret Security**: All sensitive fields properly masked
6. **Validation**: Extensive validation before service startup
7. **Operational Views**: Different dump formats for different operational needs

## Security Best Practices

- All passwords and keys marked as `secret:"true"`
- Secrets masked in all dumps by default
- Development-only secret viewing with explicit flag
- Validation prevents empty secrets in production
- Connection strings generated without exposing passwords