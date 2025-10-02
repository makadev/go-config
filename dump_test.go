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
	if opts.Format != "json" {
		t.Errorf("Expected Format=json, got %s", opts.Format)
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
func TestDump_DefaultYamlConfigMasked(t *testing.T) {
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

	// Check YAML output contains expected values
	if !strings.Contains(result, "host: localhost") {
		t.Errorf("Expected host=localhost in YAML, got: %s", result)
	}
	if !strings.Contains(result, "port: 8080") {
		t.Errorf("Expected port=8080 in YAML, got: %s", result)
	}
	if !strings.Contains(result, "password: '***'") {
		t.Errorf("Expected password to be masked as '***', got: %s", result)
	}
	if !strings.Contains(result, "debug: false") {
		t.Errorf("Expected debug=false in YAML, got: %s", result)
	}
}
