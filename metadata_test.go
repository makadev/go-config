package config_test

import (
	"reflect"
	"testing"

	"github.com/makadev/go-config"
)

type MetadataNestedConfig struct {
	Enabled bool   `config:"enabled" env:"NESTED_ENABLED" default:"false" required:"true"`
	Secret  string `config:"secret" env:"NESTED_SECRET" secret:"true"`
}

type MetadataTestConfig struct {
	ServerPort int                  `config:"server.port" env:"SERVER_PORT" default:"8080"`
	Database   string               `config:"database" env:"DB_NAME" required:"true"`
	Nested     MetadataNestedConfig `config:"nested"`
	PtrNested  *MetadataNestedConfig
	Unexported string // should be skipped
}

func TestGetFieldInfoMap_BasicStruct(t *testing.T) {
	cfg := &MetadataTestConfig{}
	infoMap, err := config.GetFieldInfoMap(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check some expected fields
	tests := []struct {
		fieldPath    string
		configKey    string
		envVar       string
		defaultValue string
		required     bool
		secret       bool
		typ          reflect.Kind
	}{
		{"ServerPort", "server.port", "SERVER_PORT", "8080", false, false, reflect.Int},
		{"Database", "database", "DB_NAME", "", true, false, reflect.String},
		{"Nested.Enabled", "enabled", "NESTED_ENABLED", "false", true, false, reflect.Bool},
		{"Nested.Secret", "secret", "NESTED_SECRET", "", false, true, reflect.String},
		{"PtrNested.Enabled", "enabled", "NESTED_ENABLED", "false", true, false, reflect.Bool},
		{"PtrNested.Secret", "secret", "NESTED_SECRET", "", false, true, reflect.String},
	}

	for _, tt := range tests {
		info, ok := infoMap[tt.fieldPath]
		if !ok {
			t.Errorf("missing field info for %s", tt.fieldPath)
			continue
		}
		if info.ConfigKey != tt.configKey {
			t.Errorf("ConfigKey mismatch for %s: got %s, want %s", tt.fieldPath, info.ConfigKey, tt.configKey)
		}
		if info.EnvVar != tt.envVar {
			t.Errorf("EnvVar mismatch for %s: got %s, want %s", tt.fieldPath, info.EnvVar, tt.envVar)
		}
		if info.DefaultValue != tt.defaultValue {
			t.Errorf("DefaultValue mismatch for %s: got %s, want %s", tt.fieldPath, info.DefaultValue, tt.defaultValue)
		}
		if info.Required != tt.required {
			t.Errorf("Required mismatch for %s: got %v, want %v", tt.fieldPath, info.Required, tt.required)
		}
		if info.Secret != tt.secret {
			t.Errorf("Secret mismatch for %s: got %v, want %v", tt.fieldPath, info.Secret, tt.secret)
		}
		if info.Type.Kind() != tt.typ {
			t.Errorf("Type mismatch for %s: got %v, want %v", tt.fieldPath, info.Type.Kind(), tt.typ)
		}
	}
}

func TestGetFieldInfoMap_UnexportedFields(t *testing.T) {
	cfg := &MetadataNestedConfig{}
	infoMap, err := config.GetFieldInfoMap(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := infoMap["Unexported"]; ok {
		t.Errorf("Unexported field should not be present in infoMap")
	}
}

func TestGetFieldInfoMap_InvalidInput(t *testing.T) {
	// Passing a non-pointer should fail
	_, err := config.GetFieldInfoMap(MetadataNestedConfig{})
	if err == nil {
		t.Errorf("expected error for non-pointer input, got nil")
	}
}
