package config

import (
	"reflect"
)

// FieldInfo contains metadata about a struct field
type FieldInfo struct {
	// FieldPath is the dot-separated path to the field (e.g., "server.port")
	FieldPath string

	// ConfigKey is the key used in configuration files
	ConfigKey string

	// EnvVar is the environment variable name
	EnvVar string

	// DefaultValue is the default value as a string
	DefaultValue string

	// Required indicates if this field is required
	Required bool

	// Secret indicates if this field should be treated as secret
	Secret bool

	// Type is the reflect.Type of the field
	Type reflect.Type

	// StructField contains the original reflect.StructField
	StructField reflect.StructField
}

// GetFieldInfoMap returns a list of all configurable fields with their metadata
//
// DO NOT USE with cyclic types
func GetFieldInfoMap(configStruct interface{}) (map[string]*FieldInfo, error) {
	if err := validateInput(configStruct); err != nil {
		return nil, err
	}

	return parseStructMetadata(reflect.TypeOf(configStruct).Elem(), "")
}

// parseStructMetadata recursively parses a struct and extracts field metadata
func parseStructMetadata(t reflect.Type, prefix string) (map[string]*FieldInfo, error) {
	metadata := make(map[string]*FieldInfo)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldPath := field.Name
		if prefix != "" {
			fieldPath = prefix + "." + field.Name
		}

		// Handle embedded structs
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			subMetadata, err := parseStructMetadata(field.Type, prefix)
			if err != nil {
				return nil, err
			}
			for k, v := range subMetadata {
				metadata[k] = v
			}
			continue
		}

		// Handle nested structs
		if field.Type.Kind() == reflect.Struct {
			subMetadata, err := parseStructMetadata(field.Type, fieldPath)
			if err != nil {
				return nil, err
			}
			for k, v := range subMetadata {
				metadata[k] = v
			}
			continue
		}

		// Handle pointer to struct
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			subMetadata, err := parseStructMetadata(field.Type.Elem(), fieldPath)
			if err != nil {
				return nil, err
			}
			for k, v := range subMetadata {
				metadata[k] = v
			}
			continue
		}

		// Parse field tags
		configKey := field.Tag.Get("config")
		if configKey == "" {
			configKey = field.Name
		}

		envVar := field.Tag.Get("env")
		defaultValue := field.Tag.Get("default")
		required := field.Tag.Get("required") == "true"
		secret := field.Tag.Get("secret") == "true"

		metadata[fieldPath] = &FieldInfo{
			FieldPath:    fieldPath,
			ConfigKey:    configKey,
			EnvVar:       envVar,
			DefaultValue: defaultValue,
			Required:     required,
			Secret:       secret,
			Type:         field.Type,
			StructField:  field,
		}
	}

	return metadata, nil
}
