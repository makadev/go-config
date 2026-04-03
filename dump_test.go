package config_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/makadev/go-config"
)

type TestConfigDump struct {
	Host     string `json:"host" env:"HOST"`
	Port     int    `json:"port" env:"PORT"`
	Password string `json:"password" env:"PASSWORD" secret:"true"`
	Debug    bool   `json:"debug" env:"DEBUG"`
}

func TestDumpConfigJSON(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    true,
	})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// Test DumpConfig (JSON format, config content, secrets masked)
	result, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:      "json",
		Content:     "config",
		MaskSecrets: true,
		MaskWith:    "***",
	})
	if err != nil {
		t.Fatalf("DumpConfig failed: %v", err)
	}

	// Parse JSON to verify structure
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	// Verify content
	if data["host"] != "localhost" {
		t.Errorf("Expected host=localhost, got %v", data["host"])
	}
	if data["port"].(float64) != 8080 {
		t.Errorf("Expected port=8080, got %v", data["port"])
	}
	if data["password"] != "***" {
		t.Errorf("Expected password to be masked, got %v", data["password"])
	}
	if data["debug"] != true {
		t.Errorf("Expected debug=true, got %v", data["debug"])
	}
}

func TestDumpConfigYAML(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    true,
	})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// Test DumpConfig (YAML format, config content, secrets masked)
	result, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:      "yaml",
		Content:     "config",
		MaskSecrets: true,
		MaskWith:    "***",
	})
	if err != nil {
		t.Fatalf("DumpConfig failed: %v", err)
	}

	// Verify YAML structure (manual verification)
	if !strings.Contains(result, "host: localhost") {
		t.Errorf("Expected host=localhost in YAML, got %v", result)
	}
	if !strings.Contains(result, "port: 8080") {
		t.Errorf("Expected port=8080 in YAML, got %v", result)
	}
	if !strings.Contains(result, "password: '***'") {
		t.Errorf("Expected password to be masked in YAML, got %v", result)
	}
	if !strings.Contains(result, "debug: true") {
		t.Errorf("Expected debug=true in YAML, got %v", result)
	}
}

func TestDumpEnv(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    true,
	})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// Test DumpEnv (text format, env content, secrets masked)
	result, err := cfg.DumpEnv()
	if err != nil {
		t.Fatalf("DumpEnv failed: %v", err)
	}

	lines := strings.Split(result, "\n")
	expected := map[string]string{
		"HOST":     "localhost",
		"PORT":     "8080",
		"PASSWORD": "***",
		"DEBUG":    "true",
	}

	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d: %v", len(expected), len(lines), lines)
	}

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			t.Errorf("Invalid line format: %s", line)
			continue
		}

		key, value := parts[0], parts[1]
		if expectedValue, ok := expected[key]; ok {
			if value != expectedValue {
				t.Errorf("Expected %s=%s, got %s=%s", key, expectedValue, key, value)
			}
		} else {
			t.Errorf("Unexpected environment variable: %s", key)
		}
	}
}

func TestDumpWithSecrets(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    true,
	})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// Test with ShowSecrets=true
	options := config.DumpOptions{
		Format:      "json",
		Content:     "config",
		MaskSecrets: false,
	}

	result, err := cfg.DumpWithOptions(&options)
	if err != nil {
		t.Fatalf("Dump with secrets failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	// Verify secret is not masked
	if data["password"] != "secret123" {
		t.Errorf("Expected password=secret123 when ShowSecrets=true, got %v", data["password"])
	}
}

func TestDumpMetadata(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    true,
	})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// Test metadata dump
	options := config.DumpOptions{
		Format:      "json",
		Content:     "metadata",
		MaskSecrets: true,
		MaskWith:    "***",
	}

	result, err := cfg.DumpWithOptions(&options)
	if err != nil {
		t.Fatalf("Dump metadata failed: %v", err)
	}

	var entries []config.DumpEntry
	if err := json.Unmarshal([]byte(result), &entries); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	// Verify we have entries with metadata
	found := false
	for _, entry := range entries {
		if entry.ConfigKey == "password" {
			found = true
			if entry.EnvVar != "PASSWORD" {
				t.Errorf("Expected EnvVar=PASSWORD, got %s", entry.EnvVar)
			}
			if !entry.IsSecret {
				t.Errorf("Expected IsSecret=true for password field")
			}
			if !entry.IsMasked {
				t.Errorf("Expected IsMasked=true for password field")
			}
			break
		}
	}

	if !found {
		t.Error("Password entry not found in metadata dump")
	}
}

func TestDumpEnvYAML(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    true,
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	result, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:      "yaml",
		Content:     "env",
		MaskSecrets: true,
		MaskWith:    "***",
	})
	if err != nil {
		t.Fatalf("Dump yaml env failed: %v", err)
	}

	if !strings.Contains(result, "HOST: localhost") {
		t.Fatalf("expected HOST in YAML env output, got %s", result)
	}
	if !strings.Contains(result, "PASSWORD: '***'") {
		t.Fatalf("expected masked PASSWORD in YAML env output, got %s", result)
	}
	if !strings.Contains(result, "DEBUG: true") {
		t.Fatalf("expected DEBUG in YAML env output, got %s", result)
	}
}

func TestDumpTableAllIncludesFieldPath(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    true,
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	result, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:      "table",
		Content:     "all",
		MaskSecrets: true,
		MaskWith:    "***",
	})
	if err != nil {
		t.Fatalf("Dump table all failed: %v", err)
	}

	if !strings.Contains(result, "CONFIG_KEY") || !strings.Contains(result, "CONFIG_NAME") || !strings.Contains(result, "ENV_VAR") ||
		!strings.Contains(result, "FIELD_PATH") || !strings.Contains(result, "VALUE") || !strings.Contains(result, "SECRET") {
		t.Fatalf("expected all table header, got %s", result)
	}
	if !strings.Contains(result, "password") || !strings.Contains(result, "PASSWORD") ||
		!strings.Contains(result, "Password") || !strings.Contains(result, "***") || !strings.Contains(result, "yes") {
		t.Fatalf("expected password row with field path, got %s", result)
	}
}

func TestDumpInvalidContentReturnsEmptyTable(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host: "localhost",
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	result, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "table",
		Content: "invalid",
	})
	if err != nil {
		t.Fatalf("expected invalid content to produce empty output, got %v", err)
	}
	if result != "" {
		t.Fatalf("expected empty output for invalid content, got %q", result)
	}
}

func TestDumpTableMetadataIncludesConfigName(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    true,
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	result, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:      "table",
		Content:     "metadata",
		MaskSecrets: true,
		MaskWith:    "***",
	})
	if err != nil {
		t.Fatalf("Dump table metadata failed: %v", err)
	}

	if !strings.Contains(result, "CONFIG_KEY") || !strings.Contains(result, "CONFIG_NAME") ||
		!strings.Contains(result, "ENV_VAR") || !strings.Contains(result, "VALUE") || !strings.Contains(result, "SECRET") {
		t.Fatalf("expected metadata table header with CONFIG_NAME, got %s", result)
	}
	if !strings.Contains(result, "password") || !strings.Contains(result, "PASSWORD") ||
		!strings.Contains(result, "***") || !strings.Contains(result, "yes") {
		t.Fatalf("expected password row in metadata table, got %s", result)
	}
}

func TestDumpTableMetadataNestedConfigNameDiffers(t *testing.T) {
	type ServerCfg struct {
		Host string `json:"host" env:"SERVER_HOST"`
		Port int    `json:"port" env:"SERVER_PORT"`
	}
	type AppCfg struct {
		Name   string    `json:"name" env:"APP_NAME"`
		Server ServerCfg `json:"server"`
	}

	cfg, err := config.NewConfig(nil, &AppCfg{
		Name:   "myapp",
		Server: ServerCfg{Host: "localhost", Port: 8080},
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// metadata table: config_key="server.host", config_name="host" — they must differ
	result, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "table",
		Content: "metadata",
	})
	if err != nil {
		t.Fatalf("Dump table metadata failed: %v", err)
	}

	if !strings.Contains(result, "CONFIG_NAME") {
		t.Fatalf("expected CONFIG_NAME header in metadata table, got %s", result)
	}
	// "server.host" is the config_key, "host" is the config_name
	if !strings.Contains(result, "server.host") {
		t.Fatalf("expected config_key 'server.host' in metadata table, got %s", result)
	}
	if !strings.Contains(result, "server.port") {
		t.Fatalf("expected config_key 'server.port' in metadata table, got %s", result)
	}

	// Verify config_name != config_key by checking rows contain both the full key and the leaf name
	lines := strings.Split(result, "\n")
	foundDiffering := false
	for _, line := range lines {
		// A row with "server.host" should also have "host" as config_name
		if strings.Contains(line, "server.host") && strings.Contains(line, "SERVER_HOST") {
			foundDiffering = true
			// The line should contain "host" as config_name (separate from "server.host")
			// Split by whitespace and check for the leaf name
			fields := strings.Fields(line)
			if len(fields) < 3 {
				t.Fatalf("expected at least 3 fields in row, got %v", fields)
			}
			// fields[0] = config_key ("server.host"), fields[1] = config_name ("host")
			if fields[0] != "server.host" {
				t.Errorf("expected config_key 'server.host', got %s", fields[0])
			}
			if fields[1] != "host" {
				t.Errorf("expected config_name 'host' (different from config_key), got %s", fields[1])
			}
			break
		}
	}
	if !foundDiffering {
		t.Fatalf("expected a row where config_key and config_name differ, got %s", result)
	}

	// Also verify "all" table includes CONFIG_NAME for nested fields
	allResult, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "table",
		Content: "all",
	})
	if err != nil {
		t.Fatalf("Dump table all failed: %v", err)
	}
	if !strings.Contains(allResult, "CONFIG_NAME") {
		t.Fatalf("expected CONFIG_NAME header in all table, got %s", allResult)
	}
	if !strings.Contains(allResult, "server.port") {
		t.Fatalf("expected config_key 'server.port' in all table, got %s", allResult)
	}
}

func TestDumpTable(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "secret123",
		Debug:    false,
	})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// Test table format
	options := config.DumpOptions{
		Format:      "table",
		Content:     "config",
		MaskSecrets: true,
		MaskWith:    "***",
	}

	result, err := cfg.DumpWithOptions(&options)
	if err != nil {
		t.Fatalf("Dump table failed: %v", err)
	}

	lines := strings.Split(result, "\n")
	if len(lines) < 3 { // header + separator + at least one entry
		t.Fatalf("Expected at least 3 lines in table output, got %d", len(lines))
	}

	// Check header
	if !strings.Contains(lines[0], "CONFIG_KEY") {
		t.Errorf("Expected header to contain CONFIG_KEY, got: %s", lines[0])
	}

	// Check that password is masked in table
	foundPassword := false
	for _, line := range lines {
		if strings.Contains(line, "password") && strings.Contains(line, "***") {
			foundPassword = true
			break
		}
	}
	if !foundPassword {
		t.Error("Expected to find masked password in table output")
	}
}


func TestDumpUnsupportedFormat(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	options := config.DumpOptions{
		Format:  "xml", // unsupported format
		Content: "config",
	}

	_, err = cfg.DumpWithOptions(&options)
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("Expected 'unsupported format' error, got: %v", err)
	}
}

func TestNewDumpOptionsDefaults(t *testing.T) {
	opts := config.NewDumpOptions()
	if opts == nil {
		t.Fatal("NewDumpOptions returned nil")
	}
	if opts.Format != "table" {
		t.Errorf("Expected Format=table, got %s", opts.Format)
	}
	if opts.Content != "config" {
		t.Errorf("Expected Content=config, got %s", opts.Content)
	}
	if !opts.MaskSecrets {
		t.Errorf("Expected MaskSecrets=true, got %v", opts.MaskSecrets)
	}
	if opts.MaskWith != "***" {
		t.Errorf("Expected MaskWith=***, got %s", opts.MaskWith)
	}
}
	func TestDumpWithOptions_NilOptions(t *testing.T) {
		cfg, err := config.NewConfig(nil, &TestConfigDump{})
		if err != nil {
			t.Fatalf("failed to initialize config: %v", err)
		}

		_, err = cfg.DumpWithOptions(nil)
		if err == nil {
			t.Fatal("expected error for nil dump options")
		}
		if !strings.Contains(err.Error(), "dump options cannot be nil") {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	func TestDumpEnvJSON(t *testing.T) {
		cfg, err := config.NewConfig(nil, &TestConfigDump{
			Host:     "localhost",
			Port:     8080,
			Password: "secret123",
			Debug:    true,
		})
		if err != nil {
			t.Fatalf("failed to initialize config: %v", err)
		}

		result, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:      "json",
			Content:     "env",
			MaskSecrets: true,
			MaskWith:    "***",
		})
		if err != nil {
			t.Fatalf("DumpWithOptions failed: %v", err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(result), &data); err != nil {
			t.Fatalf("failed to parse JSON result: %v", err)
		}

		if data["HOST"] != "localhost" {
			t.Errorf("Expected HOST=localhost, got %v", data["HOST"])
		}
		if data["PASSWORD"] != "***" {
			t.Errorf("Expected PASSWORD to be masked, got %v", data["PASSWORD"])
		}
		if data["DEBUG"] != true {
			t.Errorf("Expected DEBUG=true, got %v", data["DEBUG"])
		}
	}

	func TestDumpMetadataYAMLAndTextCollections(t *testing.T) {
		type DumpCollectionsConfig struct {
			Names  []string          `json:"names" env:"NAMES"`
			Labels map[string]string `json:"labels" env:"LABELS"`
			Secret string            `json:"secret" env:"SECRET" secret:"true"`
		}

		cfg, err := config.NewConfig(nil, &DumpCollectionsConfig{
			Names:  []string{"alice", "bob"},
			Labels: map[string]string{"role": "admin"},
			Secret: "hidden",
		})
		if err != nil {
			t.Fatalf("failed to initialize config: %v", err)
		}

		yamlResult, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:      "yaml",
			Content:     "metadata",
			MaskSecrets: true,
			MaskWith:    "***",
		})
		if err != nil {
			t.Fatalf("YAML metadata dump failed: %v", err)
		}
		if !strings.Contains(yamlResult, "configkey: names") {
			t.Fatalf("expected YAML metadata output, got %s", yamlResult)
		}
		if !strings.Contains(yamlResult, "ismasked: true") {
			t.Fatalf("expected masked metadata output, got %s", yamlResult)
		}

		textResult, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:      "text",
			Content:     "all",
			MaskSecrets: true,
			MaskWith:    "***",
		})
		if err != nil {
			t.Fatalf("Text dump failed: %v", err)
		}
		if !strings.Contains(textResult, "ConfigKey=names") || !strings.Contains(textResult, "Value=alice,bob") {
			t.Fatalf("expected slice values in text output, got %s", textResult)
		}
		if !strings.Contains(textResult, "ConfigKey=labels") || !strings.Contains(textResult, "Value=role=admin") {
			t.Fatalf("expected map values in text output, got %s", textResult)
		}
		if !strings.Contains(textResult, "ConfigKey=secret") || !strings.Contains(textResult, "(secret)") {
			t.Fatalf("expected secret metadata in text output, got %s", textResult)
		}
	}

	func TestDumpTableEnv(t *testing.T) {
		cfg, err := config.NewConfig(nil, &TestConfigDump{
			Host:     "localhost",
			Port:     8080,
			Password: "secret123",
			Debug:    true,
		})
		if err != nil {
			t.Fatalf("failed to initialize config: %v", err)
		}

		result, err := cfg.DumpWithOptions(&config.DumpOptions{
			Format:      "table",
			Content:     "env",
			MaskSecrets: true,
			MaskWith:    "***",
		})
		if err != nil {
			t.Fatalf("Dump table failed: %v", err)
		}

		if !strings.Contains(result, "ENV_VAR") || !strings.Contains(result, "VALUE") || !strings.Contains(result, "SECRET") {
			t.Fatalf("expected env table header, got %s", result)
		}
		if !strings.Contains(result, "PASSWORD") || !strings.Contains(result, "***") || !strings.Contains(result, "yes") {
			t.Fatalf("expected masked password in env table, got %s", result)
		}
	}

func TestDumpNestedConfig(t *testing.T) {
	type ServerCfg struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	type AppCfg struct {
		Name   string    `json:"name"`
		Server ServerCfg `json:"server"`
	}

	cfg, err := config.NewConfig(nil, &AppCfg{
		Name:   "myapp",
		Server: ServerCfg{Host: "localhost", Port: 8080},
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// JSON config format — exercises setNestedValue traversal
	jsonResult, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "json",
		Content: "config",
	})
	if err != nil {
		t.Fatalf("JSON dump failed: %v", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResult), &data); err != nil {
		t.Fatalf("failed to parse JSON result: %v", err)
	}
	server, ok := data["server"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected server key to be a nested map, got %T", data["server"])
	}
	if server["host"] != "localhost" {
		t.Errorf("expected server.host=localhost, got %v", server["host"])
	}
	if server["port"].(float64) != 8080 {
		t.Errorf("expected server.port=8080, got %v", server["port"])
	}

	// YAML config format — exercises formatYAML nested keys
	yamlResult, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "yaml",
		Content: "config",
	})
	if err != nil {
		t.Fatalf("YAML dump failed: %v", err)
	}
	if !strings.Contains(yamlResult, "host: localhost") {
		t.Errorf("expected host in YAML output, got: %s", yamlResult)
	}
	if !strings.Contains(yamlResult, "port: 8080") {
		t.Errorf("expected port in YAML output, got: %s", yamlResult)
	}

	// text config format — exercises nonprimitiveToString(struct) → skip
	textResult, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "text",
		Content: "config",
	})
	if err != nil {
		t.Fatalf("text dump failed: %v", err)
	}
	if !strings.Contains(textResult, "name=myapp") {
		t.Errorf("expected name in text output, got: %s", textResult)
	}

	// table config format — exercises table struct-value skip
	tableResult, err := cfg.DumpWithOptions(&config.DumpOptions{
		Format:  "table",
		Content: "config",
	})
	if err != nil {
		t.Fatalf("table dump failed: %v", err)
	}
	if !strings.Contains(tableResult, "CONFIG_KEY") {
		t.Errorf("expected table header in output, got: %s", tableResult)
	}
}

func TestDump_DefaultTableConfigMasked(t *testing.T) {
	cfg, err := config.NewConfig(nil, &TestConfigDump{
		Host:     "localhost",
		Port:     8080,
		Password: "supersecret",
		Debug:    false,
	})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	result, err := cfg.Dump()
	if err != nil {
		t.Fatalf("Dump failed: %v", err)
	}

	// Check table output contains expected header and values
	if !strings.Contains(result, "CONFIG_KEY") {
		t.Errorf("Expected CONFIG_KEY header in table output, got: %s", result)
	}
	if !strings.Contains(result, "localhost") {
		t.Errorf("Expected localhost in table output, got: %s", result)
	}
	if !strings.Contains(result, "8080") {
		t.Errorf("Expected 8080 in table output, got: %s", result)
	}
	if !strings.Contains(result, "***") {
		t.Errorf("Expected password to be masked as *** in table output, got: %s", result)
	}
}
