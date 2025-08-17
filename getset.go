package config

import (
	"fmt"
	"reflect"
	"strings"
)

func (c *Config[T]) GetFieldValue(fieldPath string) (interface{}, error) {
	field, err := c.getFieldByPath(fieldPath, false)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

func (c *Config[T]) SetFieldValue(fieldPath string, value interface{}) error {
	field, err := c.getFieldByPath(fieldPath, true)
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

func (c *Config[T]) GetConfigValue(configKey string) (interface{}, error) {
	field, err := c.getFieldByKey(configKey, false)
	if err != nil {
		return nil, err
	}
	return field.Interface(), nil
}

func (c *Config[T]) SetConfigValue(configKey string, value interface{}) error {
	field, err := c.getFieldByKey(configKey, true)
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

// getFieldByKey retrieves the reflect.Value of a config key by its path
func (c *Config[T]) getFieldByKey(configKey string, checkWritable bool) (reflect.Value, error) {
	fieldInfo, ok := c.Metadata.KeyMap[configKey]
	if !ok {
		return reflect.Value{}, fmt.Errorf("config key %q not found", configKey)
	}
	return c.getFieldByPath(fieldInfo.FieldPath, checkWritable)
}

// getFieldByPath traverses a struct to find a field by its dot-separated path
func (c *Config[T]) getFieldByPath(fieldPath string, checkWritable bool) (reflect.Value, error) {
	parts := strings.Split(fieldPath, ".")
	current := reflect.ValueOf(c.Data).Elem()

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
