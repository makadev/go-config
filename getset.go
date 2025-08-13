package config

import (
	"fmt"
	"reflect"
	"strings"
)

// GetFieldValue retrieves the value of a field by its path
func GetFieldValue(configStruct interface{}, fieldPath string) (interface{}, error) {
	structValue := reflect.ValueOf(configStruct)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}

	field, err := getFieldByPath(structValue, fieldPath, false)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

// SetFieldValue sets the value of a field by its path
func SetFieldValue(configStruct interface{}, fieldPath string, value interface{}) error {
	structValue := reflect.ValueOf(configStruct)
	if structValue.Kind() != reflect.Ptr {
		return fmt.Errorf("configStruct must be a pointer")
	}
	structValue = structValue.Elem()

	field, err := getFieldByPath(structValue, fieldPath, true)
	if err != nil {
		return err
	}

	valueReflect := reflect.ValueOf(value)
	if !valueReflect.Type().ConvertibleTo(field.Type()) {
		return fmt.Errorf("cannot convert %T to %s", value, field.Type())
	}

	field.Set(valueReflect.Convert(field.Type()))
	return nil
}

// getFieldByPath traverses a struct to find a field by its dot-separated path
func getFieldByPath(structValue reflect.Value, fieldPath string, checkWritable bool) (reflect.Value, error) {
	parts := strings.Split(fieldPath, ".")
	current := structValue

	for _, part := range parts {
		if current.Kind() == reflect.Ptr {
			if current.IsNil() {
				// Initialize nil pointer
				current.Set(reflect.New(current.Type().Elem()))
			}
			current = current.Elem()
		}

		if current.Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf("expected struct, got %s at path %s", current.Kind(), fieldPath)
		}

		fieldVal := current.FieldByName(part)
		if !fieldVal.IsValid() {
			return reflect.Value{}, fmt.Errorf("field %s not found", part)
		}
		if !fieldVal.CanInterface() {
			return reflect.Value{}, fmt.Errorf("field %s cannot be accessed", part)
		}
		if checkWritable && !fieldVal.CanSet() {
			return reflect.Value{}, fmt.Errorf("field %s cannot be set", part)
		}

		current = fieldVal
	}

	return current, nil
}
