// Package config provides a flexible configuration loading system for Go applications.
// It supports loading configuration from files (YAML/JSON) and environment variables
// with struct tag annotations for easy mapping.
//
// Thread safety:
// Config method calls (Load, Dump, DumpWithOptions, DumpEnv, Get*/Set*) are
// synchronized internally. Direct concurrent access to exported members
// (Config.Data, Config.Metadata, Config.Options) is not synchronized and must
// be protected by the caller.
package config

import (
	"fmt"
	"reflect"
	"sync"
)

type Config[T any] struct {
	// mu guards all operations done via Config methods (Load, Dump, Get/Set).
	// Direct access to exported fields below is not synchronized by this mutex.
	mu sync.RWMutex

	// Options is exported for convenience, but direct read/write access is not
	// synchronized. Prefer updating options before concurrent use.
	Options  *Options
	// Metadata is exported for advanced introspection. Direct map access is not
	// synchronized and should be externally guarded when accessed concurrently.
	Metadata *ConfigMetadata
	// Data points to the application config struct. Direct field access is not
	// synchronized and should be externally guarded when accessed concurrently.
	Data     *T
}

func NewConfig[T any](opts *Options, initData *T) (*Config[T], error) {
	if err := checkConfigStruct(initData); err != nil {
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

// checkConfigStruct ensures the provided interface is a pointer to a struct
func checkConfigStruct(configStruct interface{}) error {
	if configStruct == nil || reflect.ValueOf(configStruct).IsNil() {
		return fmt.Errorf("config struct cannot be nil")
	}

	rv := reflect.ValueOf(configStruct)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config struct must be a pointer to a struct")
	}

	return nil
}
