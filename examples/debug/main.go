package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/makadev/go-config"
)

// DebugConfig demonstrates all dump functionality and formats
type DebugConfig struct {
	App      AppSettings      `json:"app"`
	Database DatabaseSettings `json:"database"`
	External ExternalServices `json:"external"`
	Features FeatureFlags     `json:"features"`
}

type AppSettings struct {
	Name        string            `json:"name" env:"APP_NAME"`
	Version     string            `json:"version" env:"APP_VERSION"`
	Environment string            `json:"environment" env:"APP_ENV"`
	Debug       bool              `json:"debug" env:"DEBUG"`
	LogLevel    string            `json:"log_level" env:"LOG_LEVEL"`
	Metadata    map[string]string `json:"metadata"`
}

type DatabaseSettings struct {
	URL      string `json:"url" env:"DATABASE_URL" secret:"true"`
	MaxConns int    `json:"max_conns" env:"DB_MAX_CONNS"`
	Timeout  string `json:"timeout" env:"DB_TIMEOUT"`
}

type ExternalServices struct {
	PaymentAPI APIConfig `json:"payment_api" env:"PAYMENT_API_"`
	EmailAPI   APIConfig `json:"email_api" env:"EMAIL_API_"`
	StorageAPI APIConfig `json:"storage_api" env:"STORAGE_API_"`
}

type APIConfig struct {
	URL     string `json:"url" env:"URL"`
	Key     string `json:"key" env:"KEY" secret:"true"`
	Secret  string `json:"secret" env:"SECRET" secret:"true"`
	Timeout string `json:"timeout" env:"TIMEOUT"`
}

type FeatureFlags struct {
	NewUI        bool `json:"new_ui" env:"FEATURE_NEW_UI"`
	BetaFeatures bool `json:"beta_features" env:"FEATURE_BETA"`
	Analytics    bool `json:"analytics" env:"FEATURE_ANALYTICS"`
}

func main() {
	fmt.Println("=== Configuration Debug & Dump Formats Example ===")
	fmt.Println()

	// Create a complex configuration for demonstration
	cfg := createDemoConfig()

	// Initialize go-config
	opts := config.NewOptions()
	opts.ConfigPaths = []string{"debug.json"}
	goCfg, err := config.NewConfig(opts, cfg)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	fmt.Println("Configuration loaded successfully!")
	fmt.Println()

	// Demonstrate all dump formats and options
	demonstrateAllFormats(goCfg)

	// Demonstrate content filtering
	demonstrateContentFilters(goCfg)

	// Demonstrate secret handling
	demonstrateSecretHandling(goCfg)

	// Interactive mode
	if os.Getenv("INTERACTIVE") == "true" {
		runInteractiveMode(goCfg)
	} else {
		fmt.Println("\n=== Usage Instructions ===")
		fmt.Println("Set INTERACTIVE=true to explore dump options interactively")
		fmt.Println("Set various environment variables to see configuration changes:")
		fmt.Println("  export APP_NAME=\"My Custom App\"")
		fmt.Println("  export DEBUG=true")
		fmt.Println("  export DATABASE_URL=\"postgres://user:pass@localhost/db\"")
		fmt.Println("  export API_KEY=\"your-api-key\"")
	}
}

func createDemoConfig() *DebugConfig {
	return &DebugConfig{
		App: AppSettings{
			Name:        "debug-example",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
			LogLevel:    "debug",
			Metadata: map[string]string{
				"build":     "local",
				"commit":    "abc123",
				"deploy_by": "developer",
			},
		},
		Database: DatabaseSettings{
			URL:      "postgres://user:password@localhost:5432/debug_db",
			MaxConns: 10,
			Timeout:  "30s",
		},
		External: ExternalServices{
			PaymentAPI: APIConfig{
				URL:     "https://api.stripe.com/v1",
				Key:     "pk_test_123456789",
				Secret:  "sk_test_987654321",
				Timeout: "30s",
			},
			EmailAPI: APIConfig{
				URL:     "https://api.sendgrid.com/v3",
				Key:     "SG.abcdef123456",
				Secret:  "sg_secret_token",
				Timeout: "15s",
			},
			StorageAPI: APIConfig{
				URL:     "https://s3.amazonaws.com",
				Key:     "AKIAIOSFODNN7EXAMPLE",
				Secret:  "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				Timeout: "60s",
			},
		},
		Features: FeatureFlags{
			NewUI:        true,
			BetaFeatures: false,
			Analytics:    true,
		},
	}
}

func demonstrateAllFormats(cfg *config.Config[DebugConfig]) {
	fmt.Println("=== ALL DUMP FORMATS ===")
	fmt.Println()

	formats := []string{"json", "text", "table"}

	for _, format := range formats {
		fmt.Printf("--- %s FORMAT ---\n", strings.ToUpper(format))

		result, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:  format,
			Content: "config",
		})
		if err != nil {
			log.Printf("Failed to dump %s format: %v", format, err)
			continue
		}

		fmt.Println(result)
		fmt.Println()
	}
}

func demonstrateContentFilters(cfg *config.Config[DebugConfig]) {
	fmt.Println("=== CONTENT FILTERS ===")
	fmt.Println()

	contents := []struct {
		name        string
		content     string
		description string
	}{
		{"Config Keys", "config", "Shows configuration keys and values"},
		{"Environment Variables", "env", "Shows environment variable mappings"},
		{"Metadata", "metadata", "Shows config keys, env vars, and metadata"},
		{"All Fields", "all", "Shows everything including internal field paths"},
	}

	for _, c := range contents {
		fmt.Printf("--- %s ---\n", c.name)
		fmt.Printf("Description: %s\n\n", c.description)

		result, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:  "json",
			Content: c.content,
		})
		if err != nil {
			log.Printf("Failed to dump %s content: %v", c.content, err)
			continue
		}

		// Truncate long output for readability
		lines := strings.Split(result, "\n")
		if len(lines) > 15 {
			fmt.Println(strings.Join(lines[:10], "\n"))
			fmt.Printf("... (truncated, showing first 10 lines of %d total)\n", len(lines))
		} else {
			fmt.Println(result)
		}
		fmt.Println()
	}
}

func demonstrateSecretHandling(cfg *config.Config[DebugConfig]) {
	fmt.Println("=== SECRET HANDLING ===")
	fmt.Println()

	scenarios := []struct {
		name        string
		options     config.DumpOptions
		description string
	}{
		{
			"Default (Secrets Masked)",
			config.DumpOptions{Format: "table", Content: "config"},
			"Default behavior - secrets are masked with ***",
		},
		{
			"Show Secrets (Development)",
			config.DumpOptions{Format: "table", Content: "config", MaskSecrets: false},
			"Development mode - actual secret values are shown",
		},
		{
			"Security Audit",
			config.DumpOptions{Format: "table", Content: "metadata"},
			"Shows which fields are marked as secrets",
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("--- %s ---\n", scenario.name)
		fmt.Printf("Description: %s\n\n", scenario.description)

		result, err := cfg.DumpWithOptions(&scenario.options)
		if err != nil {
			log.Printf("Failed to dump %s: %v", scenario.name, err)
			continue
		}

		fmt.Println(result)
		fmt.Println()
	}
}

func runInteractiveMode(cfg *config.Config[DebugConfig]) {
	fmt.Println()
	fmt.Println("=== INTERACTIVE MODE ===")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  formats  - Show all output formats")
	fmt.Println("  content  - Show all content types")
	fmt.Println("  secrets  - Show secret handling options")
	fmt.Println("  custom   - Create custom dump options")
	fmt.Println("  quit     - Exit interactive mode")
	fmt.Println()

	for {
		fmt.Print("Enter command: ")
		var command string
		fmt.Scanln(&command)

		switch strings.ToLower(command) {
		case "formats":
			demonstrateAllFormats(cfg)
		case "content":
			demonstrateContentFilters(cfg)
		case "secrets":
			demonstrateSecretHandling(cfg)
		case "custom":
			runCustomDump(cfg)
		case "quit", "exit", "q":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}

func runCustomDump(cfg *config.Config[DebugConfig]) {
	fmt.Println("\n--- Custom Dump Options ---")

	var format, content string
	var showSecrets bool

	fmt.Print("Format (json/text/table): ")
	fmt.Scanln(&format)

	fmt.Print("Content (config/env/metadata/all): ")
	fmt.Scanln(&content)

	fmt.Print("Show secrets? (true/false): ")
	fmt.Scanln(&showSecrets)

	options := config.DumpOptions{
		Format:      format,
		Content:     content,
		MaskSecrets: !showSecrets,
	}

	fmt.Printf("\nDumping with options: %+v\n\n", options)

	result, err := cfg.DumpWithOptions(&options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(result)
	fmt.Println()
}
