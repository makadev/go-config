package config_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/makadev/go-config"
)

// TestConfig is a test configuration struct
type TestConfig struct {
	StringField   string            `config:"string_field" env:"STRING_FIELD" default:"default_string"`
	IntField      int               `config:"int_field" env:"INT_FIELD" default:"42"`
	BoolField     bool              `config:"bool_field" env:"BOOL_FIELD" default:"true"`
	DurationField time.Duration     `config:"duration_field" env:"DURATION_FIELD" default:"5m"`
	SliceField    []string          `config:"slice_field" env:"SLICE_FIELD"`
	MapField      map[string]string `config:"map_field" env:"MAP_FIELD"`
	RequiredField string            `config:"required_field" env:"REQUIRED_FIELD" required:"true"`
	SecretField   string            `config:"secret_field" env:"SECRET_FIELD" secret:"true"`

	Nested NestedConfig `config:"nested"`
}

type NestedConfig struct {
	Host string `config:"host" env:"NESTED_HOST" default:"localhost"`
	Port int    `config:"port" env:"NESTED_PORT" default:"8080"`
}

func TestGenericConfig_Defaults(t *testing.T) {
	opts := &config.Options{}
	genericCfg, err := config.NewConfig[TestConfig](opts)
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}
	genericCfg.Data.RequiredField = "set"
	err = genericCfg.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	// Check defaults
	if genericCfg.Data.StringField != "default_string" {
		t.Errorf("Expected StringField to be 'default_string', got '%s'", genericCfg.Data.StringField)
	}
	if genericCfg.Data.IntField != 42 {
		t.Errorf("Expected IntField to be 42, got %d", genericCfg.Data.IntField)
	}
	if genericCfg.Data.BoolField != true {
		t.Errorf("Expected BoolField to be true, got %v", genericCfg.Data.BoolField)
	}
	if genericCfg.Data.DurationField != 5*time.Minute {
		t.Errorf("Expected DurationField to be 5m, got %v", genericCfg.Data.DurationField)
	}
	if genericCfg.Data.Nested.Host != "localhost" {
		t.Errorf("Expected Nested.Host to be 'localhost', got '%s'", genericCfg.Data.Nested.Host)
	}
	if genericCfg.Data.Nested.Port != 8080 {
		t.Errorf("Expected Nested.Port to be 8080, got %d", genericCfg.Data.Nested.Port)
	}
	if genericCfg.Data.RequiredField != "set" {
		t.Errorf("Expected RequiredField to be 'set', got '%s'", genericCfg.Data.RequiredField)
	}
}

func TestGenericConfig_EnvironmentVariables(t *testing.T) {
	os.Setenv("STRING_FIELD", "env_string")
	os.Setenv("INT_FIELD", "100")
	os.Setenv("BOOL_FIELD", "false")
	os.Setenv("DURATION_FIELD", "10m")
	os.Setenv("SLICE_FIELD", "a,b,c")
	os.Setenv("MAP_FIELD", "key1=value1,key2=value2")
	os.Setenv("NESTED_HOST", "env.example.com")
	os.Setenv("NESTED_PORT", "9090")
	os.Setenv("REQUIRED_FIELD", "env_required")

	defer func() {
		os.Unsetenv("STRING_FIELD")
		os.Unsetenv("INT_FIELD")
		os.Unsetenv("BOOL_FIELD")
		os.Unsetenv("DURATION_FIELD")
		os.Unsetenv("SLICE_FIELD")
		os.Unsetenv("MAP_FIELD")
		os.Unsetenv("NESTED_HOST")
		os.Unsetenv("NESTED_PORT")
		os.Unsetenv("REQUIRED_FIELD")
	}()

	opts := &config.Options{}
	genericCfg, err := config.NewConfig[TestConfig](opts)
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}
	err = genericCfg.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if genericCfg.Data.StringField != "env_string" {
		t.Errorf("Expected StringField to be 'env_string', got '%s'", genericCfg.Data.StringField)
	}
	if genericCfg.Data.IntField != 100 {
		t.Errorf("Expected IntField to be 100, got %d", genericCfg.Data.IntField)
	}
	if genericCfg.Data.BoolField != false {
		t.Errorf("Expected BoolField to be false, got %v", genericCfg.Data.BoolField)
	}
	if genericCfg.Data.DurationField != 10*time.Minute {
		t.Errorf("Expected DurationField to be 10m, got %v", genericCfg.Data.DurationField)
	}
	expectedSlice := []string{"a", "b", "c"}
	if !reflect.DeepEqual(genericCfg.Data.SliceField, expectedSlice) {
		t.Errorf("Expected SliceField to be %v, got %v", expectedSlice, genericCfg.Data.SliceField)
	}
	expectedMap := map[string]string{"key1": "value1", "key2": "value2"}
	if !reflect.DeepEqual(genericCfg.Data.MapField, expectedMap) {
		t.Errorf("Expected MapField to be %v, got %v", expectedMap, genericCfg.Data.MapField)
	}
	if genericCfg.Data.Nested.Host != "env.example.com" {
		t.Errorf("Expected Nested.Host to be 'env.example.com', got '%s'", genericCfg.Data.Nested.Host)
	}
	if genericCfg.Data.Nested.Port != 9090 {
		t.Errorf("Expected Nested.Port to be 9090, got %d", genericCfg.Data.Nested.Port)
	}
}

// Test für RedactedCopy für generischen Config-Typ
func TestGenericConfig_RedactedCopy(t *testing.T) {
	opts := &config.Options{Secret: true, SecretWith: "[REDACTED]"}
	genericCfg, err := config.NewConfig[TestConfig](opts)
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}
	genericCfg.Data.RequiredField = "required"
	genericCfg.Data.StringField = "visible"
	genericCfg.Data.SecretField = "sensitive_data"
	err = genericCfg.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	secret, err := genericCfg.RedactedCopy()
	if err != nil {
		t.Fatalf("RedactedCopy failed: %v", err)
	}
	if secret.Data.SecretField != "[REDACTED]" {
		t.Errorf("Expected SecretField to be redacted, got '%s'", secret.Data.SecretField)
	}
	if secret.Data.StringField != "visible" {
		t.Errorf("Expected StringField to be unchanged, got '%s'", secret.Data.StringField)
	}
	dump, err := genericCfg.Dump("json")
	if err != nil {
		t.Fatalf("Dump failed: %v", err)
	}
	if len(dump) == 0 {
		t.Error("Dump output is empty")
	}
}

// Benchmark tests
func BenchmarkInit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := config.NewConfig[TestConfig](nil)
		if err != nil {
			b.Fatalf("NewConfig failed: %v", err)
		}
	}
}

func BenchmarkLoad(b *testing.B) {
	cfg, err := config.NewConfig[TestConfig](nil)
	if err != nil {
		b.Fatalf("NewConfig failed: %v", err)
	}
	cfg.Data.RequiredField = "benchmark"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := cfg.Load()
		if err != nil {
			b.Fatalf("Load failed: %v", err)
		}
	}
}
