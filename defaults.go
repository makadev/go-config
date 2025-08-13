package config

import (
	"fmt"
	"reflect"
)

// applyDefaults sets default values for fields that have them defined
func applyDefaults(configStruct interface{}, metadata map[string]*FieldInfo) error {
	structValue := reflect.ValueOf(configStruct).Elem()

	for fieldPath, info := range metadata {
		if info.DefaultValue == "" {
			continue // No default value specified
		}

		// Get the field to set
		field, err := getFieldByPath(structValue, fieldPath, true)
		if err != nil {
			return fmt.Errorf("failed to get field %s: %w", fieldPath, err)
		}

		// Check if field is already set (non-zero value)
		if !isZeroValue(field) {
			continue // Field already has a value, don't override with default
		}

		// Set the default value
		if err := setFieldFromString(field, info.DefaultValue, info.Type); err != nil {
			return fmt.Errorf("failed to set default value for field %s: %w", fieldPath, err)
		}
	}

	return nil
}

// isZeroValue checks if a reflect.Value contains the zero value for its type
func isZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0.0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil() || v.Len() == 0
	case reflect.Array:
		// For arrays, check if all elements are zero
		for i := 0; i < v.Len(); i++ {
			if !isZeroValue(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		// For structs, check if all fields are zero
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() && !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		// Use reflect's zero comparison as fallback
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}
}
