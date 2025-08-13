package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load loads configuration into the provided struct pointer.
// It follows the priority order: defaults -> config files -> environment variables.
//
// Example:
//
//	type Config struct {
//	    Host string `config:"host" env:"HOST" default:"localhost"`
//	    Port int    `config:"port" env:"PORT" default:"8080"`
//	}
//
//	var cfg Config
//	err := config.Load(&cfg)
func Load(configStruct interface{}) error {
	return LoadWithOptions(configStruct, &Options{})
}

// LoadWithOptions loads configuration with custom options
func LoadWithOptions(configStruct interface{}, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}

	// Validate input
	if err := validateInput(configStruct); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Set default config paths if not provided
	if len(opts.ConfigPaths) == 0 {
		opts.ConfigPaths = DefaultConfigPaths
	}

	// Parse struct metadata
	metadata, err := parseStructMetadata(reflect.TypeOf(configStruct).Elem(), "")
	if err != nil {
		return fmt.Errorf("failed to parse struct metadata: %w", err)
	}

	// Apply defaults
	if err := applyDefaults(configStruct, metadata); err != nil {
		return fmt.Errorf("failed to apply defaults: %w", err)
	}

	// Load from configuration files
	if !opts.SkipFiles {
		if err := loadFromFiles(configStruct, opts.ConfigPaths); err != nil {
			return fmt.Errorf("failed to load from files: %w", err)
		}
	}

	// Load from environment variables
	if !opts.SkipEnv {
		if err := loadFromEnv(configStruct, metadata, opts.EnvPrefix); err != nil {
			return fmt.Errorf("failed to load from environment: %w", err)
		}
	}

	// Validate required fields
	if err := validateRequired(configStruct, metadata); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// loadFromFiles loads configuration from the first available file in the provided paths
func loadFromFiles(configStruct interface{}, configPaths []string) error {
	for _, path := range configPaths {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				continue // Try next file
			}
			return fmt.Errorf("failed to stat config file %s: %w", path, err)
		}

		// File exists, try to load it
		if err := loadConfigFile(configStruct, path); err != nil {
			return fmt.Errorf("failed to load config file %s: %w", path, err)
		}

		// Successfully loaded, don't try other files
		return nil
	}

	// No config file found - this is not an error, we'll use defaults and env vars
	return nil
}

// loadConfigFile loads a single configuration file based on its extension
func loadConfigFile(configStruct interface{}, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		return loadYAMLData(configStruct, data)
	case ".json":
		return loadJSONData(configStruct, data)
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
}

// loadYAMLData unmarshals YAML data into the config struct
func loadYAMLData(configStruct interface{}, data []byte) error {
	if err := yaml.Unmarshal(data, configStruct); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return nil
}

// loadJSONData unmarshals JSON data into the config struct
func loadJSONData(configStruct interface{}, data []byte) error {
	if err := json.Unmarshal(data, configStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// FindConfigFile searches for the first available config file in the given paths
func FindConfigFile(paths []string) (string, error) {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return path, nil
		}
	}
	return "", fmt.Errorf("no config file found in paths: %v", paths)
}

// GetConfigFormat determines the format of a config file based on its extension
func GetConfigFormat(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		return "yaml", nil
	case ".json":
		return "json", nil
	default:
		return "", fmt.Errorf("unsupported file extension: %s", ext)
	}
}
