package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/makadev/go-config"
)

// ServerConfig represents a complete web server configuration
type ServerConfig struct {
	Server   HTTPServerConfig `json:"server"`
	TLS      TLSConfig        `json:"tls"`
	CORS     CORSConfig       `json:"cors"`
	Logging  LoggingConfig    `json:"logging"`
	Security SecurityConfig   `json:"security"`
}

type HTTPServerConfig struct {
	Host         string        `json:"host" env:"SERVER_HOST"`
	Port         int           `json:"port" env:"SERVER_PORT"`
	ReadTimeout  time.Duration `json:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `json:"idle_timeout" env:"SERVER_IDLE_TIMEOUT"`
}

type TLSConfig struct {
	Enabled  bool   `json:"enabled" env:"TLS_ENABLED"`
	CertFile string `json:"cert_file" env:"TLS_CERT_FILE"`
	KeyFile  string `json:"key_file" env:"TLS_KEY_FILE" secret:"true"`
	Port     int    `json:"port" env:"TLS_PORT"`
}

type CORSConfig struct {
	Enabled        bool     `json:"enabled" env:"CORS_ENABLED"`
	AllowedOrigins []string `json:"allowed_origins" env:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods []string `json:"allowed_methods" env:"CORS_ALLOWED_METHODS"`
	AllowedHeaders []string `json:"allowed_headers" env:"CORS_ALLOWED_HEADERS"`
}

type LoggingConfig struct {
	Level  string `json:"level" env:"LOG_LEVEL"`
	Format string `json:"format" env:"LOG_FORMAT"`
	Output string `json:"output" env:"LOG_OUTPUT"`
}

type SecurityConfig struct {
	JWTSecret      string        `json:"jwt_secret" env:"JWT_SECRET" secret:"true"`
	SessionTimeout time.Duration `json:"session_timeout" env:"SESSION_TIMEOUT"`
	RateLimitRPS   int           `json:"rate_limit_rps" env:"RATE_LIMIT_RPS"`
	MaxRequestSize int64         `json:"max_request_size" env:"MAX_REQUEST_SIZE"`
}

func main() {
	fmt.Println("=== Web Server Configuration Example ===")
	fmt.Println()

	// Create default configuration
	defaultConfig := &ServerConfig{
		Server: HTTPServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		TLS: TLSConfig{
			Enabled:  false,
			CertFile: "server.crt",
			KeyFile:  "server.key",
			Port:     8443,
		},
		CORS: CORSConfig{
			Enabled:        true,
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
		Security: SecurityConfig{
			JWTSecret:      "default-jwt-secret-key",
			SessionTimeout: 24 * time.Hour,
			RateLimitRPS:   100,
			MaxRequestSize: 10 * 1024 * 1024, // 10MB
		},
	}

	// Initialize configuration
	opts := config.NewOptions()
	opts.ConfigPaths = []string{"server.json"}
	cfg, err := config.NewConfig(opts, defaultConfig)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}
	// Load configuration from file and environment variables
	if err := cfg.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := validateConfig(cfg.Data); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}
	fmt.Println("Configuration validation passed")
	fmt.Println()

	// Demo: Configuration dumps for different purposes
	demoConfigurationDumps(cfg)

	// Demo: Start a simple HTTP server with the configuration
	if os.Getenv("START_SERVER") == "true" {
		startDemoServer(cfg.Data)
	} else {
		fmt.Println("Set START_SERVER=true to actually start the server")
	}
}

func validateConfig(cfg *ServerConfig) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.TLS.Enabled {
		if cfg.TLS.CertFile == "" || cfg.TLS.KeyFile == "" {
			return fmt.Errorf("TLS enabled but cert/key files not specified")
		}
	}

	if cfg.Security.JWTSecret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}

	return nil
}

func demoConfigurationDumps(cfg *config.Config[ServerConfig]) {
	fmt.Println("=== Configuration Dumps ===")
	fmt.Println()

	// 1. Operations dump - table format for overview
	fmt.Println("1. Operations Overview (Table Format):")
	opsDump, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "table",
		Content: "config",
	})
	if err != nil {
		log.Printf("Failed to dump ops config: %v", err)
	} else {
		fmt.Println(opsDump)
	}
	fmt.Println()

	// 2. Environment export for deployment
	fmt.Println("2. Environment Variables for Deployment:")
	envDump, err := cfg.DumpEnv()
	if err != nil {
		log.Printf("Failed to dump env: %v", err)
	} else {
		fmt.Println(envDump)
	}
	fmt.Println()

	// 3. Complete configuration for documentation
	fmt.Println("3. Complete Configuration (JSON):")
	jsonDump, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "json",
		Content: "config",
	})
	if err != nil {
		log.Printf("Failed to dump config: %v", err)
	} else {
		fmt.Println(jsonDump)
	}
	fmt.Println()

	// 4. Debug metadata for troubleshooting
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("4. Debug Metadata (All Fields):")
		debugDump, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:  "json",
			Content: "metadata",
		})
		if err != nil {
			log.Printf("Failed to dump debug: %v", err)
		} else {
			fmt.Println(debugDump)
		}
		fmt.Println()
	}
}

func startDemoServer(cfg *ServerConfig) {
	fmt.Printf("Starting HTTP server on %s:%d\n", cfg.Server.Host, cfg.Server.Port)

	// Create HTTP server with configuration
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Simple handler that shows configuration info
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server running with configuration:\n")
		fmt.Fprintf(w, "- Host: %s\n", cfg.Server.Host)
		fmt.Fprintf(w, "- Port: %d\n", cfg.Server.Port)
		fmt.Fprintf(w, "- TLS Enabled: %v\n", cfg.TLS.Enabled)
		fmt.Fprintf(w, "- CORS Enabled: %v\n", cfg.CORS.Enabled)
		fmt.Fprintf(w, "- Log Level: %s\n", cfg.Logging.Level)
		fmt.Fprintf(w, "- Rate Limit: %d RPS\n", cfg.Security.RateLimitRPS)
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		// Create temporary config instance for dumping
		tempCfg, _ := config.NewConfig(nil, cfg)
		dump, err := tempCfg.Dump()
		if err != nil {
			http.Error(w, "Failed to dump config", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(dump))
	})

	// Start server with graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	fmt.Printf("Server started! Visit:\n")
	fmt.Printf("- http://%s:%d/ for server info\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("- http://%s:%d/config for configuration dump\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Println("Press Ctrl+C to stop")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}
