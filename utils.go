package config

import (
	"fmt"
	"reflect"
	"strings"
)

// ParseBool parses a string to boolean with support for various formats
func ParseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "t", "yes", "y", "1", "on":
		return true, nil
	case "false", "f", "no", "n", "0", "off", "":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean format")
	}
}

// IsZeroValue checks if a reflect.Value contains the zero value for its type
func IsZeroValue(v reflect.Value) bool {
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
			if !IsZeroValue(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		// For structs, check if all fields are zero
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() && !IsZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		// Use reflect's zero comparison as fallback
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}
}
