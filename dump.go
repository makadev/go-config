package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// StructToFlatConfigMap converts a struct to a flat map with config keys.
func StructToFlatConfigMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	structToFlatConfigMapHelper(reflect.ValueOf(v), "", result)
	return result
}

func structToFlatConfigMapHelper(val reflect.Value, prefix string, result map[string]interface{}) {
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		configKey := field.Tag.Get("config")
		if configKey == "" {
			configKey = field.Name
		}
		key := configKey
		if prefix != "" {
			key = prefix + "." + configKey
		}
		fieldVal := val.Field(i)
		switch fieldVal.Kind() {
		case reflect.Struct:
			structToFlatConfigMapHelper(fieldVal, key, result)
		case reflect.Ptr:
			if !fieldVal.IsNil() {
				structToFlatConfigMapHelper(fieldVal, key, result)
			}
		case reflect.Slice:
			for j := 0; j < fieldVal.Len(); j++ {
				elem := fieldVal.Index(j)
				elemKey := key + fmt.Sprintf("[%d]", j)
				if elem.Kind() == reflect.Struct || (elem.Kind() == reflect.Ptr && !elem.IsNil() && elem.Elem().Kind() == reflect.Struct) {
					structToFlatConfigMapHelper(elem, elemKey, result)
				} else {
					result[elemKey] = elem.Interface()
				}
			}
		case reflect.Map:
			for _, mapKey := range fieldVal.MapKeys() {
				valElem := fieldVal.MapIndex(mapKey)
				elemKey := key + "[" + mapKey.String() + "]"
				if valElem.Kind() == reflect.Struct || (valElem.Kind() == reflect.Ptr && !valElem.IsNil() && valElem.Elem().Kind() == reflect.Struct) {
					structToFlatConfigMapHelper(valElem, elemKey, result)
				} else {
					result[elemKey] = valElem.Interface()
				}
			}
		default:
			result[key] = fieldVal.Interface()
		}
	}
}

// StructToConfigMap converts a struct to a map with config keys.
func StructToConfigMap(v interface{}) map[string]interface{} {
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	typ := val.Type()
	result := make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		configKey := field.Tag.Get("config")
		if configKey == "" {
			configKey = field.Name
		}
		fieldVal := val.Field(i)
		switch fieldVal.Kind() {
		case reflect.Struct:
			result[configKey] = StructToConfigMap(fieldVal.Interface())
		case reflect.Ptr:
			if !fieldVal.IsNil() {
				result[configKey] = StructToConfigMap(fieldVal.Interface())
			}
		case reflect.Slice:
			slice := make([]interface{}, fieldVal.Len())
			for j := 0; j < fieldVal.Len(); j++ {
				elem := fieldVal.Index(j)
				if elem.Kind() == reflect.Struct || (elem.Kind() == reflect.Ptr && !elem.IsNil() && elem.Elem().Kind() == reflect.Struct) {
					slice[j] = StructToConfigMap(elem.Interface())
				} else {
					slice[j] = elem.Interface()
				}
			}
			result[configKey] = slice
		case reflect.Map:
			mapResult := make(map[string]interface{})
			for _, key := range fieldVal.MapKeys() {
				valElem := fieldVal.MapIndex(key)
				if valElem.Kind() == reflect.Struct || (valElem.Kind() == reflect.Ptr && !valElem.IsNil() && valElem.Elem().Kind() == reflect.Struct) {
					mapResult[key.String()] = StructToConfigMap(valElem.Interface())
				} else {
					mapResult[key.String()] = valElem.Interface()
				}
			}
			result[configKey] = mapResult
		default:
			result[configKey] = fieldVal.Interface()
		}
	}
	return result
}

// Dump converts a config struct to a string in the specified format.
func Dump(configStruct interface{}, metadata map[string]*FieldInfo, format string, redactSecrets bool, redactWith string, envPrefix string) (string, error) {
	var dumpData interface{} = configStruct

	// Apply secret masking if requested
	if redactSecrets {
		secretData, err := RedactedCopy(configStruct, metadata, redactWith)
		if err != nil {
			return "", fmt.Errorf("failed to mask secret fields: %w", err)
		}
		dumpData = secretData
	}

	// Handle the format
	switch strings.ToLower(format) {
	case "yaml", "yml":
		configMap := StructToConfigMap(dumpData)
		return dumpYAML(configMap)
	case "json":
		configMap := StructToConfigMap(dumpData)
		return dumpJSON(configMap)
	case "env", ".env":
		return dumpEnv(dumpData, envPrefix)
	case "flat":
		return dumpFlat(dumpData)
	default:
		return "", fmt.Errorf("unsupported format: %s (supported: yaml, json)", format)
	}
}

func dumpYAML(data interface{}) (string, error) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return string(yamlData), nil
}

func dumpJSON(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonData), nil
}

// dumpEnv dumps only the fields with an `env` tag in a .env compatible format
func dumpEnv(configStruct interface{}, envPrefix string) (string, error) {
	fields, err := GetFieldInfoMap(configStruct)
	if err != nil {
		return "", fmt.Errorf("failed to list fields: %w", err)
	}
	var envOutput strings.Builder

	for _, field := range fields {
		if envKey := field.EnvVar; envKey != "" {
			val, err := GetFieldValue(configStruct, field.FieldPath)
			if err != nil {
				return "", fmt.Errorf("failed to get field value for %s: %w", field.FieldPath, err)
			}
			envOutput.WriteString(fmt.Sprintf("%s%s=%v\n", envPrefix, envKey, val))
		}
	}
	return envOutput.String(), nil
}

func dumpFlat(configStruct interface{}) (string, error) {
	fields, err := GetFieldInfoMap(configStruct)
	if err != nil {
		return "", fmt.Errorf("failed to list fields: %w", err)
	}

	var flatOutput strings.Builder
	for _, field := range fields {
		val, err := GetFieldValue(configStruct, field.FieldPath)
		if err != nil {
			return "", fmt.Errorf("failed to get field value for %s: %w", field.FieldPath, err)
		}
		flatOutput.WriteString(fmt.Sprintf("%s=%v\n", field.EnvVar, val))
	}
	return flatOutput.String(), nil
}
