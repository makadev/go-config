package config

import (
	"fmt"
	"reflect"
	"strings"
)

// validateRequired checks that all required fields have been set
func validateRequired(configStruct interface{}, metadata map[string]*FieldInfo) error {
	structValue := reflect.ValueOf(configStruct).Elem()

	var errors []string

	// Check fields marked as required in struct tags
	for fieldPath, info := range metadata {
		if !info.Required {
			continue
		}

		field, err := getFieldByPath(structValue, fieldPath, false)
		if err != nil {
			errors = append(errors, fmt.Sprintf("required field %s: %v", fieldPath, err))
			continue
		}

		if isZeroValue(field) {
			errors = append(errors, fmt.Sprintf("required field %s is not set", fieldPath))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// validateInput ensures the provided interface is a pointer to a struct
func validateInput(configStruct interface{}) error {
	if configStruct == nil {
		return fmt.Errorf("config struct cannot be nil")
	}

	rv := reflect.ValueOf(configStruct)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("config struct must be a pointer")
	}

	if rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config struct must be a pointer to a struct")
	}

	if !rv.Elem().CanSet() {
		return fmt.Errorf("config struct must be settable")
	}

	return nil
}

// Validate performs basic validation on a config struct
func Validate(configStruct interface{}, metadata map[string]*FieldInfo) error {
	if err := validateInput(configStruct); err != nil {
		return err
	}

	// Validate required fields
	return validateRequired(configStruct, metadata)
}
