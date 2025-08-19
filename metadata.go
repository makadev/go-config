package config

import (
	"fmt"
	"reflect"
	"regexp"
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

	if err := getStructMetadata(c.Options, metadata, reflect.TypeOf(c.Data).Elem(), "", "", c.Options.EnvPrefix); err != nil {
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

var (
	EnvVarRE     = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]+$`)
	ConfigNameRE = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

func getStructMetadata(opts *Options, metadata *ConfigMetadata, t reflect.Type, path_prefix string, key_prefix string, env_prefix string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get config name
		configName := getConfigName(opts, field)

		if !ConfigNameRE.MatchString(configName) {
			return fmt.Errorf("config name %q must match regex %q", configName, ConfigNameRE.String())
		}

		// Create fieldpath and config key
		fieldPath := path_prefix + field.Name
		configKey := key_prefix + configName
		newEnvPrefix := env_prefix

		// Get environment variable name if exists and env loading is enabled
		envVar := ""
		if !opts.SkipEnv {
			envVar = field.Tag.Get("env")
			// Set environment variable name if not provided and AutoEnv is enabled
			if opts.AutoEnv {
				// when autoenv is enabled, the "env" field overrides the automatic building
				// otherwise the "env" field is build based on the full configkey like server.port -> "SERVER_PORT"
				if envVar == "" {
					envVar = env_prefix + strings.ToUpper(strings.ReplaceAll(configKey, ".", "_"))
				}
			} else {
				// if autoenv is not enabled we use "env" field for building the env var
				// and structs are specially handled as additional prefix
				// f.e. &struct {
				// 	Server struct {
				// 		Host string `env:"HOST"`
				// 	} `env:"SERVER"`
				// }{} -> Server.Host will be filled from "SERVER_HOST"
				if envVar != "" {
					envVar = env_prefix + envVar
					newEnvPrefix = envVar
				}
			}
		}

		if envVar != "" && !EnvVarRE.MatchString(envVar) {
			return fmt.Errorf("environment variable %q must match regex %q", envVar, EnvVarRE.String())
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
			if err := getStructMetadata(opts, metadata, field.Type, fieldPath+".", configKey+".", newEnvPrefix); err != nil {
				return err
			}
			continue
		case reflect.Ptr:
			// Handle ptr to struct
			if field.Type.Elem().Kind() == reflect.Struct {
				if err := getStructMetadata(opts, metadata, field.Type.Elem(), fieldPath+".", configKey+".", newEnvPrefix); err != nil {
					return err
				}
			}
			continue
		}
	}

	return nil
}
