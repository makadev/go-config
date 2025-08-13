package config_test

import (
	"testing"

	"github.com/makadev/go-config"
)

type GetSetNested struct {
	Value int
}

type GetSetTestConfig struct {
	Name   string
	Age    int
	Nested GetSetNested
}

func TestGetFieldValue_SimpleField(t *testing.T) {
	cfg := GetSetTestConfig{Name: "Alice", Age: 30}
	val, err := config.GetFieldValue(cfg, "Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "Alice" {
		t.Errorf("expected 'Alice', got %v", val)
	}
}

func TestGetFieldValue_NestedField(t *testing.T) {
	cfg := GetSetTestConfig{Nested: GetSetNested{Value: 42}}
	val, err := config.GetFieldValue(cfg, "Nested.Value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestGetFieldValue_PointerStruct(t *testing.T) {
	cfg := &GetSetTestConfig{Name: "Bob"}
	val, err := config.GetFieldValue(cfg, "Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "Bob" {
		t.Errorf("expected 'Bob', got %v", val)
	}
}

func TestGetFieldValue_FieldNotFound(t *testing.T) {
	cfg := GetSetTestConfig{}
	_, err := config.GetFieldValue(cfg, "Unknown")
	if err == nil {
		t.Error("expected error for unknown field, got nil")
	}
}
func TestSetFieldValue_SimpleField(t *testing.T) {
	cfg := &GetSetTestConfig{Name: "Alice", Age: 30}
	err := config.SetFieldValue(cfg, "Name", "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "Bob" {
		t.Errorf("expected 'Bob', got %v", cfg.Name)
	}
}

func TestSetFieldValue_NestedField(t *testing.T) {
	cfg := &GetSetTestConfig{Nested: GetSetNested{Value: 10}}
	err := config.SetFieldValue(cfg, "Nested.Value", 99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Nested.Value != 99 {
		t.Errorf("expected 99, got %v", cfg.Nested.Value)
	}
}

func TestSetFieldValue_PointerRequired(t *testing.T) {
	cfg := GetSetTestConfig{Name: "Alice"}
	err := config.SetFieldValue(cfg, "Name", "Charlie")
	if err == nil {
		t.Error("expected error for non-pointer configStruct, got nil")
	}
}

func TestSetFieldValue_FieldNotFound(t *testing.T) {
	cfg := &GetSetTestConfig{}
	err := config.SetFieldValue(cfg, "Unknown", "value")
	if err == nil {
		t.Error("expected error for unknown field, got nil")
	}
}

func TestSetFieldValue_TypeMismatch(t *testing.T) {
	cfg := &GetSetTestConfig{Age: 25}
	err := config.SetFieldValue(cfg, "Age", "not-an-int")
	if err == nil {
		t.Error("expected error for type mismatch, got nil")
	}
}
