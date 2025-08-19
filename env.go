package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// loadFromEnv loads values from environment variables into the config struct
func (c *Config[T]) loadFromEnv() error {
	for envVar, info := range c.Metadata.EnvMap {
		envKey := envVar

		envValue, ok := os.LookupEnv(envKey)
		if !ok {
			continue // Environment variable not set
		}

		// Get the field to set
		field, err := c.getFieldByPath(info.FieldPath, true)
		if err != nil {
			return fmt.Errorf("failed to get field %s: %w", info.FieldPath, err)
		}

		// Convert and set the value
		if err := setFieldFromString(field, envValue, info.StructField.Type); err != nil {
			return fmt.Errorf("failed to set field %s from env var %s: %w", info.FieldPath, envKey, err)
		}
	}

	return nil
}

// setFieldFromString converts a string value to the appropriate type and sets it on the field
func setFieldFromString(field reflect.Value, value string, fieldType reflect.Type) error {
	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(fieldType.Elem()))
		}
		return setFieldFromString(field.Elem(), value, fieldType.Elem())
	}

	switch fieldType.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		boolVal, err := ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value %q: %w", value, err)
		}
		field.SetBool(boolVal)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldType == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("invalid duration %q: %w", value, err)
			}
			field.Set(reflect.ValueOf(duration))
		} else {
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer value %q: %w", value, err)
			}
			if field.OverflowInt(intVal) {
				return fmt.Errorf("integer value %q overflows %s", value, fieldType.Kind())
			}
			field.SetInt(intVal)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value %q: %w", value, err)
		}
		if field.OverflowUint(uintVal) {
			return fmt.Errorf("unsigned integer value %q overflows %s", value, fieldType.Kind())
		}
		field.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value %q: %w", value, err)
		}
		if field.OverflowFloat(floatVal) {
			return fmt.Errorf("float value %q overflows %s", value, fieldType.Kind())
		}
		field.SetFloat(floatVal)

	case reflect.Slice:
		return setSliceFromString(field, value, fieldType)

	case reflect.Map:
		return setMapFromString(field, value, fieldType)

	default:
		return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}

	return nil
}

// setSliceFromString parses a comma-separated string into a slice
func setSliceFromString(field reflect.Value, value string, fieldType reflect.Type) error {
	if value == "" {
		field.Set(reflect.MakeSlice(fieldType, 0, 0))
		return nil
	}

	parts := strings.Split(value, ",")
	elemType := fieldType.Elem()
	slice := reflect.MakeSlice(fieldType, len(parts), len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		elemField := slice.Index(i)
		if err := setFieldFromString(elemField, part, elemType); err != nil {
			return fmt.Errorf("failed to set slice element %d: %w", i, err)
		}
	}

	field.Set(slice)
	return nil
}

// setMapFromString parses a comma-separated key=value string into a map
func setMapFromString(field reflect.Value, value string, fieldType reflect.Type) error {
	if value == "" {
		field.Set(reflect.MakeMap(fieldType))
		return nil
	}

	if fieldType.Key().Kind() != reflect.String {
		return fmt.Errorf("only string keys are supported for map environment variables")
	}

	pairs := strings.Split(value, ",")
	mapVal := reflect.MakeMap(fieldType)

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid key=value pair: %s", pair)
		}

		key := reflect.ValueOf(strings.TrimSpace(kv[0]))
		val := reflect.New(fieldType.Elem()).Elem()

		if err := setFieldFromString(val, strings.TrimSpace(kv[1]), fieldType.Elem()); err != nil {
			return fmt.Errorf("failed to set map value for key %s: %w", kv[0], err)
		}

		mapVal.SetMapIndex(key, val)
	}

	field.Set(mapVal)
	return nil
}
