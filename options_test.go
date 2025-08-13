package config_test

import (
	"reflect"
	"testing"

	"github.com/makadev/go-config"
)

func TestNewOptions_DefaultValues(t *testing.T) {
	opts := config.NewOptions()

	if !reflect.DeepEqual(opts.ConfigPaths, config.DefaultConfigPaths) {
		t.Errorf("ConfigPaths = %v, want %v", opts.ConfigPaths, config.DefaultConfigPaths)
	}
	if opts.EnvPrefix != "" {
		t.Errorf("EnvPrefix = %q, want \"\"", opts.EnvPrefix)
	}
	if opts.SkipEnv != false {
		t.Errorf("SkipEnv = %v, want false", opts.SkipEnv)
	}
	if opts.SkipFiles != false {
		t.Errorf("SkipFiles = %v, want false", opts.SkipFiles)
	}
	if opts.Secret != true {
		t.Errorf("Secret = %v, want true", opts.Secret)
	}
	if opts.SecretWith != "[REDACTED]" {
		t.Errorf("SecretWith = %q, want \"[REDACTED]\"", opts.SecretWith)
	}
}

func TestDefaultConfigPaths(t *testing.T) {
	expected := []string{
		"config.local.yaml",
		"config.local.json",
		"config.yaml",
		"config.json",
	}
	if !reflect.DeepEqual(config.DefaultConfigPaths, expected) {
		t.Errorf("DefaultConfigPaths = %v, want %v", config.DefaultConfigPaths, expected)
	}
}
