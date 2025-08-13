package main

import (
	"fmt"
	"log"
	"time"

	"github.com/makadev/go-config"
)

// AppConfig represents the main application configuration
type AppConfig struct {
	// Server configuration
	Server ServerConfig `config:"server"`

	// Database configuration
	Database DatabaseConfig `config:"database"`

	// Logging configuration
	Logging LoggingConfig `config:"logging"`

	// Feature flags
	Features FeatureFlags `config:"features"`

	// Simple fields with defaults and env vars
	AppName    string            `config:"app_name" env:"APP_NAME" default:"MyApp"`
	Version    string            `config:"version" env:"VERSION" default:"1.0.0"`
	Debug      bool              `config:"debug" env:"DEBUG" default:"false"`
	Timeout    time.Duration     `config:"timeout" env:"TIMEOUT" default:"30s"`
	MaxRetries int               `config:"max_retries" env:"MAX_RETRIES" default:"3"`
	AllowedIPs []string          `config:"allowed_ips" env:"ALLOWED_IPS"`
	Metadata   map[string]string `config:"metadata" env:"METADATA"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host         string        `config:"host" env:"SERVER_HOST" default:"localhost"`
	Port         int           `config:"port" env:"SERVER_PORT" default:"8080"`
	TLS          bool          `config:"tls" env:"SERVER_TLS" default:"false"`
	CertFile     string        `config:"cert_file" env:"SERVER_CERT_FILE"`
	KeyFile      string        `config:"key_file" env:"SERVER_KEY_FILE"`
	ReadTimeout  time.Duration `config:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"5s"`
	WriteTimeout time.Duration `config:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"5s"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string        `config:"host" env:"DB_HOST" default:"localhost"`
	Port         int           `config:"port" env:"DB_PORT" default:"5432"`
	Name         string        `config:"name" env:"DB_NAME" required:"true"`
	Username     string        `config:"username" env:"DB_USER" required:"true"`
	Password     string        `config:"password" env:"DB_PASSWORD" required:"true" secret:"true"`
	SSLMode      string        `config:"ssl_mode" env:"DB_SSL_MODE" default:"prefer"`
	MaxConns     int           `config:"max_connections" env:"DB_MAX_CONNS" default:"10"`
	MaxIdleConns int           `config:"max_idle_connections" env:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnTimeout  time.Duration `config:"connection_timeout" env:"DB_CONN_TIMEOUT" default:"10s"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `config:"level" env:"LOG_LEVEL" default:"info"`
	Format     string `config:"format" env:"LOG_FORMAT" default:"json"`
	Output     string `config:"output" env:"LOG_OUTPUT" default:"stdout"`
	Filename   string `config:"filename" env:"LOG_FILENAME"`
	MaxSize    int    `config:"max_size" env:"LOG_MAX_SIZE" default:"100"`
	MaxBackups int    `config:"max_backups" env:"LOG_MAX_BACKUPS" default:"3"`
}

// FeatureFlags holds feature toggle configuration
type FeatureFlags struct {
	EnableMetrics   bool `config:"enable_metrics" env:"FEATURE_METRICS" default:"true"`
	EnableTracing   bool `config:"enable_tracing" env:"FEATURE_TRACING" default:"false"`
	EnableProfiling bool `config:"enable_profiling" env:"FEATURE_PROFILING" default:"false"`
	BetaFeatures    bool `config:"beta_features" env:"FEATURE_BETA" default:"false"`
}

func main() {
	var cfg AppConfig

	// Load configuration with default options
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print loaded configuration
	fmt.Println("=== Loaded Configuration ===")
	fmt.Printf("App: %s v%s (Debug: %v)\n", cfg.AppName, cfg.Version, cfg.Debug)
	fmt.Printf("Server: %s:%d (TLS: %v)\n", cfg.Server.Host, cfg.Server.Port, cfg.Server.TLS)
	fmt.Printf("Database: %s@%s:%d/%s\n", cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	fmt.Printf("Logging: %s level to %s\n", cfg.Logging.Level, cfg.Logging.Output)
	fmt.Printf("Timeout: %v, Max Retries: %d\n", cfg.Timeout, cfg.MaxRetries)

	if len(cfg.AllowedIPs) > 0 {
		fmt.Printf("Allowed IPs: %v\n", cfg.AllowedIPs)
	}

	if len(cfg.Metadata) > 0 {
		fmt.Printf("Metadata: %v\n", cfg.Metadata)
	}

	// Example with custom options
	fmt.Println("\n=== Custom Options Example ===")
	var customCfg AppConfig

	options := &config.Options{
		ConfigPaths: []string{"app.yaml", "app.json", "config.yaml"},
		EnvPrefix:   "MYAPP",
	}

	if err := config.LoadWithOptions(&customCfg, options); err != nil {
		log.Printf("Custom config failed (expected): %v", err)
	}

	// Example: List all configurable fields
	fmt.Println("\n=== Available Configuration Fields ===")
	fields, err := config.GetFieldInfoMap(&cfg)
	if err != nil {
		log.Printf("Failed to list fields: %v", err)
		return
	}

	for fieldPath, info := range fields {
		fmt.Printf("Field: %-20s Config: %-15s Env: %-15s Default: %s\n",
			fieldPath, info.ConfigKey, info.EnvVar, info.DefaultValue)
	}

	// Example: Get and set individual field values
	fmt.Println("\n=== Field Value Manipulation ===")

	// Get a field value
	if port, err := config.GetFieldValue(&cfg, "Server.Port"); err == nil {
		fmt.Printf("Current server port: %v\n", port)
	}

	// Set a field value
	if err := config.SetFieldValue(&cfg, "Server.Port", 9090); err == nil {
		fmt.Printf("Updated server port to: %d\n", cfg.Server.Port)
	}

	out, err := config.Dump(&cfg, fields, "json", true, "***", "MY_APP_")
	if err != nil {
		log.Printf("Failed to dump config: %v", err)
		return
	}
	fmt.Println("=== Dumped Configuration ===")
	fmt.Println(string(out))
}
