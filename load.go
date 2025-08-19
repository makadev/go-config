package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

func (c *Config[T]) Load() error {
	// Set default config paths if not provided
	if len(c.Options.ConfigPaths) == 0 {
		c.Options.ConfigPaths = DefaultConfigPaths
	}

	// Load from configuration files
	if !c.Options.SkipFiles {
		if err := c.loadFromFiles(); err != nil {
			return fmt.Errorf("failed to load from files: %w", err)
		}
	}

	// Load from environment variables
	if !c.Options.SkipEnv {
		if err := c.loadFromEnv(); err != nil {
			return fmt.Errorf("failed to load from environment: %w", err)
		}
	}

	return nil
}

// loadFromFiles loads configuration from the first available file in the provided paths
func (c *Config[T]) loadFromFiles() error {
	for _, path := range c.Options.ConfigPaths {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				continue // Try next file
			}
			return fmt.Errorf("failed to stat config file %s: %w", path, err)
		}

		// File exists, try to load it
		if err := c.loadConfigFile(path); err != nil {
			return fmt.Errorf("failed to load config file %s: %w", path, err)
		}

		// Successfully loaded, don't try other files
		return nil
	}

	// No config file found - this is not an error, we'll use defaults and env vars
	return nil
}

// loadConfigFile loads a single configuration file based on its extension
func (c *Config[T]) loadConfigFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		return c.loadYAMLData(data)
	case ".json":
		return c.loadJSONData(data)
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
}

// loadYAMLData unmarshals YAML data into the config struct
func (c *Config[T]) loadYAMLData(data []byte) error {
	if err := yaml.Unmarshal(data, c.Data); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return nil
}

// loadJSONData unmarshals JSON data into the config struct
func (c *Config[T]) loadJSONData(data []byte) error {
	if err := json.Unmarshal(data, c.Data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}
