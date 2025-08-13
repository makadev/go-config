package config_test

import (
	"strings"
	"testing"

	"github.com/makadev/go-config"
)

// Test struct for validation
type ValidateTestConfig struct {
	Name  string `required:"true"`
	Value int
}

type ValidateNoRequiredConfig struct {
	Name  string
	Value int
}

func TestValidate_NilInput(t *testing.T) {
	err := config.Validate(nil, nil)
	if err == nil || err.Error() != "config struct cannot be nil" {
		t.Errorf("expected error for nil input, got: %v", err)
	}
}

func TestValidate_NonPointerInput(t *testing.T) {
	tc := ValidateTestConfig{}
	err := config.Validate(tc, nil)
	if err == nil || err.Error() != "config struct must be a pointer" {
		t.Errorf("expected error for non-pointer input, got: %v", err)
	}
}

func TestValidate_NonStructPointerInput(t *testing.T) {
	x := 42
	err := config.Validate(&x, nil)
	if err == nil || err.Error() != "config struct must be a pointer to a struct" {
		t.Errorf("expected error for pointer to non-struct, got: %v", err)
	}
}

func TestValidate_RequiredFieldNotSet(t *testing.T) {
	tc := &ValidateTestConfig{}
	metadata, err := config.GetFieldInfoMap(tc)
	if err != nil {
		t.Fatalf("failed to get field metadata: %v", err)
	}
	err = config.Validate(tc, metadata)
	if err == nil || !strings.Contains(err.Error(), "required field Name is not set") {
		t.Errorf("expected error for required field not set, got: %v", err)
	}
}

func TestValidate_RequiredFieldSet(t *testing.T) {
	tc := &ValidateTestConfig{Name: "foo"}
	metadata, err := config.GetFieldInfoMap(tc)
	if err != nil {
		t.Fatalf("failed to get field metadata: %v", err)
	}
	err = config.Validate(tc, metadata)
	if err != nil {
		t.Errorf("expected no error when required field is set, got: %v", err)
	}
}

func TestValidate_NoRequiredFields(t *testing.T) {
	tc := &ValidateNoRequiredConfig{}
	metadata, err := config.GetFieldInfoMap(tc)
	if err != nil {
		t.Fatalf("failed to get field metadata: %v", err)
	}
	err = config.Validate(tc, metadata)
	if err != nil {
		t.Errorf("expected no error when no required fields, got: %v", err)
	}
}
