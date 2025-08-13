// Package config provides a flexible configuration loading system for Go applications.
// It supports loading configuration from files (YAML/JSON) and environment variables
// with struct tag annotations for easy mapping.
package config

import (
	"fmt"
)

type Config[T any] struct {
	Options  *Options
	Metadata map[string]*FieldInfo
	Data     *T
}

func NewConfig[T any](opts *Options) (*Config[T], error) {
	empty := new(T)
	metadata, err := GetFieldInfoMap(empty)
	if err != nil {
		return nil, fmt.Errorf("failed to list fields: %w", err)
	}
	if opts == nil {
		opts = NewOptions()
	}
	return &Config[T]{
		Options:  opts,
		Metadata: metadata,
		Data:     empty,
	}, nil
}

func (c *Config[T]) Load() error {
	// Validate input
	if err := validateInput(c.Data); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Set default config paths if not provided
	if len(c.Options.ConfigPaths) == 0 {
		c.Options.ConfigPaths = DefaultConfigPaths
	}

	// Apply defaults
	if err := applyDefaults(c.Data, c.Metadata); err != nil {
		return fmt.Errorf("failed to apply defaults: %w", err)
	}

	// Load from configuration files
	if !c.Options.SkipFiles {
		if err := loadFromFiles(c.Data, c.Options.ConfigPaths); err != nil {
			return fmt.Errorf("failed to load from files: %w", err)
		}
	}

	// Load from environment variables
	if !c.Options.SkipEnv {
		if err := loadFromEnv(c.Data, c.Metadata, c.Options.EnvPrefix); err != nil {
			return fmt.Errorf("failed to load from environment: %w", err)
		}
	}

	// Validate required fields
	if err := validateRequired(c.Data, c.Metadata); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func (c *Config[T]) GetFieldValue(fieldPath string) (interface{}, error) {
	return GetFieldValue(c.Data, fieldPath)
}

func (c *Config[T]) SetFieldValue(fieldPath string, value interface{}) error {
	return SetFieldValue(c.Data, fieldPath, value)
}

func (c *Config[T]) GetFieldInfo(fieldPath string) (*FieldInfo, error) {
	info, ok := c.Metadata[fieldPath]
	if !ok {
		return nil, fmt.Errorf("field not found: %s", fieldPath)
	}
	return info, nil
}

func (c *Config[T]) GetFieldInfoMap() map[string]*FieldInfo {
	return c.Metadata
}

func (c *Config[T]) RedactedCopy() (*Config[T], error) {
	secretData, err := RedactedCopy(c.Data, c.Metadata, c.Options.SecretWith)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret copy: %w", err)
	}
	typedSecretData := secretData.(*T)
	return &Config[T]{Options: c.Options, Metadata: c.Metadata, Data: typedSecretData}, nil
}

func (c *Config[T]) Dump(format string) (string, error) {
	return Dump(c.Data, c.Metadata, format, c.Options.Secret, c.Options.SecretWith, c.Options.EnvPrefix)
}
