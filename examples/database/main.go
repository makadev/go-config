package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/makadev/go-config"
)

// DatabaseConfig represents database connection and configuration
type DatabaseConfig struct {
	Primary   DatabaseConnection `json:"primary" env:"PRIMARY_"`
	ReadOnly  DatabaseConnection `json:"read_only" env:"READ_ONLY_"`
	Cache     CacheConfig        `json:"cache"`
	Migration MigrationConfig    `json:"migration"`
	Backup    BackupConfig       `json:"backup"`
}

type DatabaseConnection struct {
	Driver          string        `json:"driver" env:"DB_DRIVER"`
	Host            string        `json:"host" env:"DB_HOST"`
	Port            int           `json:"port" env:"DB_PORT"`
	Database        string        `json:"database" env:"DB_NAME"`
	Username        string        `json:"username" env:"DB_USERNAME"`
	Password        string        `json:"password" env:"DB_PASSWORD" secret:"true"`
	SSLMode         string        `json:"ssl_mode" env:"DB_SSL_MODE"`
	ConnectTimeout  time.Duration `json:"connect_timeout" env:"DB_CONNECT_TIMEOUT"`
	MaxOpenConns    int           `json:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `json:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME"`
}

type CacheConfig struct {
	Enabled  bool              `json:"enabled" env:"CACHE_ENABLED"`
	Type     string            `json:"type" env:"CACHE_TYPE"`
	Redis    RedisConfig       `json:"redis"`
	Memcache MemcacheConfig    `json:"memcache"`
	TTL      map[string]string `json:"ttl"`
}

type RedisConfig struct {
	Host     string `json:"host" env:"REDIS_HOST"`
	Port     int    `json:"port" env:"REDIS_PORT"`
	Password string `json:"password" env:"REDIS_PASSWORD" secret:"true"`
	Database int    `json:"database" env:"REDIS_DB"`
	PoolSize int    `json:"pool_size" env:"REDIS_POOL_SIZE"`
}

type MemcacheConfig struct {
	Servers []string `json:"servers" env:"MEMCACHE_SERVERS"`
}

type MigrationConfig struct {
	Enabled     bool          `json:"enabled" env:"MIGRATION_ENABLED"`
	Directory   string        `json:"directory" env:"MIGRATION_DIR"`
	Table       string        `json:"table" env:"MIGRATION_TABLE"`
	LockKey     string        `json:"lock_key" env:"MIGRATION_LOCK_KEY"`
	LockTimeout time.Duration `json:"lock_timeout" env:"MIGRATION_LOCK_TIMEOUT"`
}

type BackupConfig struct {
	Enabled   bool          `json:"enabled" env:"BACKUP_ENABLED"`
	Schedule  string        `json:"schedule" env:"BACKUP_SCHEDULE"`
	Retention time.Duration `json:"retention" env:"BACKUP_RETENTION"`
	S3Config  S3Config      `json:"s3"`
}

type S3Config struct {
	Bucket          string `json:"bucket" env:"S3_BUCKET"`
	Region          string `json:"region" env:"S3_REGION"`
	AccessKeyID     string `json:"access_key_id" env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `json:"secret_access_key" env:"S3_SECRET_ACCESS_KEY" secret:"true"`
	Endpoint        string `json:"endpoint" env:"S3_ENDPOINT"`
}

func main() {
	fmt.Println("=== Database Configuration Example ===")
	fmt.Println()

	// Create default configuration
	defaultConfig := &DatabaseConfig{
		Primary: DatabaseConnection{
			Driver:          "postgres",
			Host:            "localhost",
			Port:            5432,
			Database:        "myapp",
			Username:        "postgres",
			Password:        "postgres",
			SSLMode:         "disable",
			ConnectTimeout:  30 * time.Second,
			MaxOpenConns:    25,
			MaxIdleConns:    10,
			ConnMaxLifetime: 5 * time.Minute,
		},
		ReadOnly: DatabaseConnection{
			Driver:          "postgres",
			Host:            "localhost",
			Port:            5433,
			Database:        "myapp",
			Username:        "readonly",
			Password:        "readonly",
			SSLMode:         "disable",
			ConnectTimeout:  30 * time.Second,
			MaxOpenConns:    15,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Cache: CacheConfig{
			Enabled: true,
			Type:    "redis",
			Redis: RedisConfig{
				Host:     "localhost",
				Port:     6379,
				Password: "",
				Database: 0,
				PoolSize: 10,
			},
			Memcache: MemcacheConfig{
				Servers: []string{"localhost:11211"},
			},
			TTL: map[string]string{
				"users":    "1h",
				"sessions": "24h",
				"cache":    "5m",
			},
		},
		Migration: MigrationConfig{
			Enabled:     true,
			Directory:   "./migrations",
			Table:       "schema_migrations",
			LockKey:     "migration_lock",
			LockTimeout: 10 * time.Minute,
		},
		Backup: BackupConfig{
			Enabled:   false,
			Schedule:  "0 2 * * *",
			Retention: 30 * 24 * time.Hour,
			S3Config: S3Config{
				Bucket:          "myapp-backups",
				Region:          "us-east-1",
				AccessKeyID:     "",
				SecretAccessKey: "",
				Endpoint:        "",
			},
		},
	}

	// Initialize configuration
	opts := config.NewOptions()
	opts.ConfigPaths = []string{"database.json"}
	cfg, err := config.NewConfig(opts, defaultConfig)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	// Load configuration from file and environment variables
	if err := cfg.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := validateDatabaseConfig(cfg.Data); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}
	fmt.Println("Configuration validation passed")
	fmt.Println()

	// Demo various configuration dumps
	demoDatabaseConfigDumps(cfg)

	// Demo connection string generation (without secrets)
	demoConnectionStrings(cfg.Data)
}

func validateDatabaseConfig(cfg *DatabaseConfig) error {
	// Validate primary database
	if err := validateConnection("primary", &cfg.Primary); err != nil {
		return err
	}

	// Validate read-only database
	if err := validateConnection("read-only", &cfg.ReadOnly); err != nil {
		return err
	}

	// Validate cache configuration
	if cfg.Cache.Enabled {
		switch cfg.Cache.Type {
		case "redis":
			if cfg.Cache.Redis.Host == "" {
				return fmt.Errorf("redis host cannot be empty when cache type is redis")
			}
		case "memcache":
			if len(cfg.Cache.Memcache.Servers) == 0 {
				return fmt.Errorf("memcache servers cannot be empty when cache type is memcache")
			}
		default:
			return fmt.Errorf("unsupported cache type: %s", cfg.Cache.Type)
		}
	}

	// Validate backup configuration
	if cfg.Backup.Enabled {
		if cfg.Backup.S3Config.Bucket == "" {
			return fmt.Errorf("S3 bucket cannot be empty when backup is enabled")
		}
		if cfg.Backup.S3Config.AccessKeyID == "" || cfg.Backup.S3Config.SecretAccessKey == "" {
			return fmt.Errorf("S3 credentials cannot be empty when backup is enabled")
		}
	}

	return nil
}

func validateConnection(name string, conn *DatabaseConnection) error {
	if conn.Host == "" {
		return fmt.Errorf("%s database host cannot be empty", name)
	}
	if conn.Port <= 0 || conn.Port > 65535 {
		return fmt.Errorf("%s database port must be between 1 and 65535", name)
	}
	if conn.Database == "" {
		return fmt.Errorf("%s database name cannot be empty", name)
	}
	if conn.Username == "" {
		return fmt.Errorf("%s database username cannot be empty", name)
	}
	return nil
}

func demoDatabaseConfigDumps(cfg *config.Config[DatabaseConfig]) {
	fmt.Println("=== Database Configuration Dumps ===")
	fmt.Println()

	// 1. Security audit - show which fields contain secrets
	fmt.Println("1. Security Audit (Metadata with Secret Detection):")
	securityDump, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "table",
		Content: "metadata",
	})
	if err != nil {
		log.Printf("Failed to dump security audit: %v", err)
	} else {
		fmt.Println(securityDump)
	}
	fmt.Println()

	// 2. Environment variables for containerized deployment
	fmt.Println("2. Environment Variables for Docker/Kubernetes:")
	envDump, err := cfg.DumpEnv()
	if err != nil {
		log.Printf("Failed to dump env: %v", err)
	} else {
		fmt.Println(envDump)
	}
	fmt.Println()

	// 3. Configuration overview for documentation
	fmt.Println("3. Configuration Documentation (JSON):")
	jsonDump, err := cfg.Dump()
	if err != nil {
		log.Printf("Failed to dump config: %v", err)
	} else {
		fmt.Println(jsonDump)
	}
	fmt.Println()

	// 4. Production deployment verification
	fmt.Println("4. Production Deployment Verification:")
	verifyDump, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "text",
		Content: "config",
	})
	if err != nil {
		log.Printf("Failed to dump verification: %v", err)
	} else {
		// Show only non-secret fields for production verification
		fmt.Println("Non-secret configuration values:")
		fmt.Println(verifyDump)
	}
	fmt.Println()

	// 5. Development dump with secrets (if explicitly enabled)
	if os.Getenv("SHOW_SECRETS") == "true" {
		fmt.Println("5. Development Dump (WITH SECRETS - DEV ONLY!):")
		devDump, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:      "json",
			Content:     "config",
			MaskSecrets: false,
		})
		if err != nil {
			log.Printf("Failed to dump dev config: %v", err)
		} else {
			fmt.Println(devDump)
		}
		fmt.Println()
	}
}

func demoConnectionStrings(cfg *DatabaseConfig) {
	fmt.Println("=== Database Connection Information ===")
	fmt.Println()

	// Generate safe connection info (without passwords)
	fmt.Println("Primary Database:")
	fmt.Printf("  Type: %s\n", cfg.Primary.Driver)
	fmt.Printf("  Host: %s:%d\n", cfg.Primary.Host, cfg.Primary.Port)
	fmt.Printf("  Database: %s\n", cfg.Primary.Database)
	fmt.Printf("  Username: %s\n", cfg.Primary.Username)
	fmt.Printf("  SSL Mode: %s\n", cfg.Primary.SSLMode)
	fmt.Printf("  Max Connections: %d open, %d idle\n", cfg.Primary.MaxOpenConns, cfg.Primary.MaxIdleConns)
	fmt.Println()

	fmt.Println("Read-Only Database:")
	fmt.Printf("  Type: %s\n", cfg.ReadOnly.Driver)
	fmt.Printf("  Host: %s:%d\n", cfg.ReadOnly.Host, cfg.ReadOnly.Port)
	fmt.Printf("  Database: %s\n", cfg.ReadOnly.Database)
	fmt.Printf("  Username: %s\n", cfg.ReadOnly.Username)
	fmt.Printf("  SSL Mode: %s\n", cfg.ReadOnly.SSLMode)
	fmt.Printf("  Max Connections: %d open, %d idle\n", cfg.ReadOnly.MaxOpenConns, cfg.ReadOnly.MaxIdleConns)
	fmt.Println()

	if cfg.Cache.Enabled {
		fmt.Printf("Cache: %s\n", cfg.Cache.Type)
		if cfg.Cache.Type == "redis" {
			fmt.Printf("  Redis: %s:%d (DB %d)\n", cfg.Cache.Redis.Host, cfg.Cache.Redis.Port, cfg.Cache.Redis.Database)
		} else if cfg.Cache.Type == "memcache" {
			fmt.Printf("  Memcache Servers: %v\n", cfg.Cache.Memcache.Servers)
		}
		fmt.Println()
	}

	if cfg.Migration.Enabled {
		fmt.Printf("Migrations: Enabled\n")
		fmt.Printf("  Directory: %s\n", cfg.Migration.Directory)
		fmt.Printf("  Table: %s\n", cfg.Migration.Table)
		fmt.Println()
	}

	if cfg.Backup.Enabled {
		fmt.Printf("Backup: Enabled\n")
		fmt.Printf("  Schedule: %s\n", cfg.Backup.Schedule)
		fmt.Printf("  S3 Bucket: %s (%s)\n", cfg.Backup.S3Config.Bucket, cfg.Backup.S3Config.Region)
		fmt.Printf("  Retention: %v\n", cfg.Backup.Retention)
		fmt.Println()
	}

	fmt.Println("Note: Passwords and secret keys are masked for security")
	fmt.Println("Set SHOW_SECRETS=true to see actual values (development only)")
}
