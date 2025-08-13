package config

import (
	"fmt"
	"reflect"
)

// RedactedCopy creates a copy of the config struct with secret fields.
//
// DO NOT USE with cyclic references
func RedactedCopy(configStruct interface{}, metadata map[string]*FieldInfo, redactWith string) (interface{}, error) {
	if err := validateInput(configStruct); err != nil {
		return nil, err
	}
	structValue := reflect.ValueOf(configStruct).Elem()

	// Create a new struct of the same type
	redactedStruct := reflect.New(structValue.Type()).Interface()

	// Set default redacted value if not provided
	if redactWith == "" {
		redactWith = "..."
	}

	// Copy all non-redacted values
	if err := copyStructValuesWithPath(structValue, reflect.ValueOf(redactedStruct).Elem(), metadata, redactWith, ""); err != nil {
		return nil, fmt.Errorf("failed to copy struct values: %w", err)
	}

	return redactedStruct, nil
}

// copyStructValuesWithPath recursively copies struct values with path tracking
func copyStructValuesWithPath(src, dst reflect.Value, metadata map[string]*FieldInfo, redactedValue string, pathPrefix string) error {
	structType := src.Type()

	for i := 0; i < src.NumField(); i++ {
		field := structType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldPath := field.Name
		if pathPrefix != "" {
			fieldPath = pathPrefix + "." + field.Name
		}

		srcField := src.Field(i)
		dstField := dst.Field(i)

		// Handle nested structs
		if field.Type.Kind() == reflect.Struct {
			if err := copyStructValuesWithPath(srcField, dstField, metadata, redactedValue, fieldPath); err != nil {
				return err
			}
			continue
		}

		// Handle pointer to struct
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if !srcField.IsNil() {
				dstField.Set(reflect.New(field.Type.Elem()))
				if err := copyStructValuesWithPath(srcField.Elem(), dstField.Elem(), metadata, redactedValue, fieldPath); err != nil {
					return err
				}
			}
			continue
		}

		// Check if field should be secret
		if info, exists := metadata[fieldPath]; exists && info.Secret {
			// Set secret value
			if err := setRedactedValue(dstField, redactedValue, field.Type); err != nil {
				return fmt.Errorf("failed to set secret value for field %s: %w", fieldPath, err)
			}
		} else {
			// Copy original value
			if srcField.CanInterface() && dstField.CanSet() {
				dstField.Set(srcField)
			}
		}
	}

	return nil
}

// setRedactedValue sets an appropriate redacted value based on the field type
func setRedactedValue(field reflect.Value, redactedValue string, fieldType reflect.Type) error {
	switch fieldType.Kind() {
	case reflect.String:
		field.SetString(redactedValue)
	case reflect.Slice:
		if fieldType.Elem().Kind() == reflect.String {
			// For string slices, create a single-element slice with redacted value
			slice := reflect.MakeSlice(fieldType, 1, 1)
			slice.Index(0).SetString(redactedValue)
			field.Set(slice)
		} else {
			// For other slice types, set empty slice
			field.Set(reflect.MakeSlice(fieldType, 0, 0))
		}
	case reflect.Map:
		// For maps, set empty map
		field.Set(reflect.MakeMap(fieldType))
	case reflect.Ptr:
		// For pointers, handle the underlying type
		if field.IsNil() {
			field.Set(reflect.New(fieldType.Elem()))
		}
		return setRedactedValue(field.Elem(), redactedValue, fieldType.Elem())
	default:
		// For other types (int, bool, etc.), try to set zero value
		field.Set(reflect.Zero(fieldType))
	}

	return nil
}
