// Package config provides a flexible configuration loading system for Go applications.
// It supports loading configuration from files (YAML/JSON) and environment variables
// with struct tag annotations for easy mapping.
package config

import (
	"fmt"
)

type Config[T any] struct {
	Options  *Options
	Metadata *ConfigMetadata
	Data     *T
}

func NewConfig[T any](opts *Options, initData *T) (*Config[T], error) {
	if err := CheckConfigStruct(initData); err != nil {
		return nil, fmt.Errorf("invalid config struct: %w", err)
	}
	if opts == nil {
		opts = NewOptions()
	}

	cfg := &Config[T]{
		Options:  opts,
		Metadata: nil,
		Data:     initData,
	}
	if err := cfg.initMetadata(); err != nil {
		return nil, fmt.Errorf("failed to list fields: %w", err)
	}

	return cfg, nil
}
