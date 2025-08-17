package config

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldInfo contains metadata about a struct field
type FieldInfo struct {
	// FieldPath is the dot-separated path to the field using struct field names
	// f.e. "Server.Port"
	FieldPath string

	// ConfigKey is the full key used to reference the field inside the config using configured config key names
	// f.e. "server.port"
	ConfigKey string

	// ConfigName is the key used in configuration files
	// f.e. "port"
	ConfigName string

	// EnvVar is the environment variable name
	EnvVar string

	// Secret indicates if this field should be treated as secret
	Secret bool

	// StructField contains the original reflect.StructField
	StructField reflect.StructField
}

type ConfigMetadata struct {
	FieldPathMap map[string]*FieldInfo
	EnvMap       map[string]*FieldInfo
	KeyMap       map[string]*FieldInfo
}

func (c *Config[T]) initMetadata() error {
	metadata := &ConfigMetadata{
		FieldPathMap: make(map[string]*FieldInfo),
		EnvMap:       make(map[string]*FieldInfo),
		KeyMap:       make(map[string]*FieldInfo),
	}

	if err := getStructMetadata(c.Options, metadata, reflect.TypeOf(c.Data).Elem(), "", ""); err != nil {
		return err
	}
	c.Metadata = metadata
	return nil
}

func getConfigName(opts *Options, field reflect.StructField) string {
	for _, v := range opts.ConfigTags {
		configName := field.Tag.Get(v)
		if configName != "" {
			return configName
		}
	}
	return strings.ToLower(field.Name)
}

func getStructMetadata(opts *Options, metadata *ConfigMetadata, t reflect.Type, path_prefix string, key_prefix string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get config name
		configName := getConfigName(opts, field)

		// Create fieldpath and config key
		fieldPath := path_prefix + field.Name
		configKey := key_prefix + configName

		// Get environment variable name if exists and env loading is enabled
		envVar := ""
		if !opts.SkipEnv {
			envVar = field.Tag.Get("env")
			// Set environment variable name if not provided and AutoEnv is enabled
			if opts.AutoEnv && envVar == "" {
				envVar = strings.ToUpper(strings.ReplaceAll(configKey, ".", "_"))
			}
		}

		secret := field.Tag.Get("secret") == "true"

		// Init FieldInfo
		fieldInfo := &FieldInfo{
			FieldPath:   fieldPath,
			ConfigKey:   configKey,
			ConfigName:  configName,
			EnvVar:      envVar,
			Secret:      secret,
			StructField: field,
		}

		// Store FieldInfo by field path
		metadata.FieldPathMap[fieldPath] = fieldInfo

		// Store FieldInfo by env var and check for dups
		if envVar != "" {
			if _, ok := metadata.EnvMap[envVar]; ok {
				return fmt.Errorf("duplicate env var %q found", envVar)
			}
			metadata.EnvMap[envVar] = fieldInfo
		}

		// Store FieldInfo by config key and check for dups
		if _, ok := metadata.KeyMap[configKey]; ok {
			return fmt.Errorf("duplicate config key %q found", configKey)
		}
		metadata.KeyMap[configKey] = fieldInfo

		// recursive handle nested, embedded and referenced structs
		switch field.Type.Kind() {
		case reflect.Struct:
			// Handle embedded or nested structs
			if err := getStructMetadata(opts, metadata, field.Type, fieldPath+".", configKey+"."); err != nil {
				return err
			}
			continue
		case reflect.Ptr:
			if field.Type.Elem().Kind() == reflect.Struct {
				// Handle structs
				if err := getStructMetadata(opts, metadata, field.Type.Elem(), fieldPath+".", configKey+"."); err != nil {
					return err
				}
			}
			continue
		}
	}

	return nil
}
