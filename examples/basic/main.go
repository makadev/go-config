package main

import (
	"fmt"
	"log"
	"os"

	"github.com/makadev/go-config"
)

// AppConfig represents a simple application configuration
type AppConfig struct {
	AppName  string `json:"app_name" env:"APP_NAME"`
	Version  string `json:"version" env:"VERSION"`
	Debug    bool   `json:"debug" env:"DEBUG"`
	Port     int    `json:"port" env:"PORT"`
	APIKey   string `json:"api_key" env:"API_KEY" secret:"true"`
	LogLevel string `json:"log_level" env:"LOG_LEVEL"`
}

func main() {
	fmt.Println("=== Basic Config Load & Dump Example ===")

	// Create default configuration
	defaultConfig := &AppConfig{
		AppName:  "my-app",
		Version:  "1.0.0",
		Debug:    false,
		Port:     8080,
		APIKey:   "default-secret-key",
		LogLevel: "info",
	}

	// Initialize configuration with file loading
	opts := config.NewOptions()
	opts.ConfigPaths = []string{"config.json"}
	cfg, err := config.NewConfig(opts, defaultConfig)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	// Load configuration from file and environment variables
	if err := cfg.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Demo: Show current configuration values
	fmt.Println("Current Configuration:")
	fmt.Printf("  App Name: %s\n", cfg.Data.AppName)
	fmt.Printf("  Version: %s\n", cfg.Data.Version)
	fmt.Printf("  Debug: %v\n", cfg.Data.Debug)
	fmt.Printf("  Port: %d\n", cfg.Data.Port)
	fmt.Printf("  API Key: %s (secret)\n", cfg.Data.APIKey)
	fmt.Printf("  Log Level: %s\n", cfg.Data.LogLevel)
	fmt.Println()

	// Demo: Dump configuration in different formats
	fmt.Println("=== Configuration Dumps ===")

	// 1. Default dump (table format, secrets masked)
	fmt.Println("1. Default Dump (table, secrets masked):")
	jsonDump, err := cfg.Dump()
	if err != nil {
		log.Fatalf("Failed to dump config: %v", err)
	}
	fmt.Println(jsonDump)
	fmt.Println()

	// 2. Environment variables dump
	fmt.Println("2. Environment Variables Dump:")
	envDump, err := cfg.DumpEnv()
	if err != nil {
		log.Fatalf("Failed to dump env: %v", err)
	}
	fmt.Println(envDump)
	fmt.Println()

	// 3. Table format dump
	fmt.Println("3. Table Format Dump:")
	tableDump, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "table",
		Content: "config",
	})
	if err != nil {
		log.Fatalf("Failed to dump table: %v", err)
	}
	fmt.Println(tableDump)
	fmt.Println()

	// 4. Debug dump with secrets (for development only!)
	if os.Getenv("SHOW_SECRETS") == "true" {
		fmt.Println("4. Debug Dump (with secrets - development only!):")
		debugDump, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:      "json",
			Content:     "config",
			MaskSecrets: false,
		})
		if err != nil {
			log.Fatalf("Failed to dump debug: %v", err)
		}
		fmt.Println(debugDump)
		fmt.Println()
	}

	fmt.Println("=== Example Usage Instructions ===")
	fmt.Println("1. Create a config.json file with your settings")
	fmt.Println("2. Set environment variables to override specific values:")
	fmt.Println("   export APP_NAME=\"my-custom-app\"")
	fmt.Println("   export DEBUG=true")
	fmt.Println("   export PORT=3000")
	fmt.Println("3. Run with SHOW_SECRETS=true to see actual secret values")
	fmt.Println("4. Use the dump functions for debugging and deployment verification")
}
