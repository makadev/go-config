package config_test

import (
	"reflect"
	"testing"

	"github.com/makadev/go-config"
)

type GetSet_TestStruct struct {
	StringField  string  `config:"string_field" env:"STRING_FIELD"`
	IntField     int     `config:"int_field" env:"INT_FIELD"`
	BoolField    bool    `config:"bool_field" env:"BOOL_FIELD"`
	FloatField   float64 `config:"float_field" env:"FLOAT_FIELD"`
	PtrField     *string `config:"ptr_field" env:"PTR_FIELD"`
	NestedStruct struct {
		NestedString string `config:"nested_string" env:"NESTED_STRING"`
		NestedInt    int    `config:"nested_int" env:"NESTED_INT"`
	} `config:"nested"`
	PtrNested *struct {
		PtrNestedString string `config:"ptr_nested_string" env:"PTR_NESTED_STRING"`
	} `config:"ptr_nested"`
}

func Test_GetFieldValue_Success(t *testing.T) {
	testValue := "test"
	initData := &GetSet_TestStruct{
		StringField: "hello",
		IntField:    42,
		BoolField:   true,
		FloatField:  3.14,
		PtrField:    &testValue,
		NestedStruct: struct {
			NestedString string `config:"nested_string" env:"NESTED_STRING"`
			NestedInt    int    `config:"nested_int" env:"NESTED_INT"`
		}{
			NestedString: "nested",
			NestedInt:    100,
		},
	}

	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	tests := []struct {
		name      string
		fieldPath string
		expected  interface{}
	}{
		{"string field", "StringField", "hello"},
		{"int field", "IntField", 42},
		{"bool field", "BoolField", true},
		{"float field", "FloatField", 3.14},
		{"pointer field", "PtrField", &testValue},
		{"nested string", "NestedStruct.NestedString", "nested"},
		{"nested int", "NestedStruct.NestedInt", 100},
	}

	// test getting multiple values expecting success
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := cfg.GetFieldValue(tt.fieldPath)
			if err != nil {
				t.Fatalf("GetFieldValue failed: %v", err)
			}
			if value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, value)
			}
		})
	}
}

func Test_GetFieldValue_NonExistentField(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// accessing a non existent key/path should return an error
	_, err = cfg.GetFieldValue("NonExistentField")
	if err == nil {
		t.Fatal("expected error for non-existent field")
	}
	if !contains(err.Error(), "field NonExistentField not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_GetFieldValue_InvalidPath(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// accessing a non existent key/path should return an error
	_, err = cfg.GetFieldValue("StringField.NonExistent")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
	if !contains(err.Error(), "expected struct") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_SetFieldValue_Success(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	tests := []struct {
		name      string
		fieldPath string
		value     interface{}
		expected  interface{}
	}{
		{"string field", "StringField", "new value", "new value"},
		{"int field", "IntField", 123, 123},
		{"bool field", "BoolField", true, true},
		{"float field", "FloatField", 2.71, 2.71},
		{"nested string", "NestedStruct.NestedString", "nested value", "nested value"},
		{"nested int", "NestedStruct.NestedInt", 999, 999},
	}

	// test setting multiple values expecting success
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cfg.SetFieldValue(tt.fieldPath, tt.value)
			if err != nil {
				t.Fatalf("SetFieldValue failed: %v", err)
			}

			value, err := cfg.GetFieldValue(tt.fieldPath)
			if err != nil {
				t.Fatalf("GetFieldValue failed: %v", err)
			}
			if value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, value)
			}
		})
	}
}

func Test_SetFieldValue_TypeConversion(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// SetFieldValue should initialize the accessed path and convert int8 to int
	err = cfg.SetFieldValue("IntField", int8(-42))
	if err != nil {
		t.Fatalf("SetFieldValue failed: %v", err)
	}

	// GetFieldValue should now read the set value
	value, err := cfg.GetFieldValue("IntField")
	if err != nil {
		t.Fatalf("GetFieldValue failed: %v", err)
	}
	i, ok := value.(int)
	if !ok {
		t.Fatalf("expected int, got %T", value)
	}
	if i != -42 {
		t.Errorf("expected -42, got %v", i)
	}
}

func Test_SetFieldValue_NonExistentField(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// accessing a non existent key/path should return an error
	err = cfg.SetFieldValue("NonExistentField", "value")
	if err == nil {
		t.Fatal("expected error for non-existent field")
	}
	if !contains(err.Error(), "field NonExistentField not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_SetFieldValue_IncompatibleType(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// trying to set a string value to an int field should return an error
	err = cfg.SetFieldValue("IntField", "not a number")
	if err == nil {
		t.Fatal("expected error for incompatible type")
	}
	if !contains(err.Error(), "cannot convert") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_SetFieldValue_NonWritable(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// trying to set a string value to an int field should return an error
	err = cfg.SetFieldValue("IntField", "not a number")
	if err == nil {
		t.Fatal("expected error for incompatible type")
	}
	if !contains(err.Error(), "cannot convert") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_SetFieldValue_NilPointer(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// trying to set a value on a nil pointer should initialize the pointed type and set the value
	err = cfg.SetFieldValue("PtrNested.PtrNestedString", "test value")
	if err != nil {
		t.Fatalf("SetFieldValue failed: %v", err)
	}

	// GetFieldValue should now read the set value
	value, err := cfg.GetFieldValue("PtrNested.PtrNestedString")
	if err != nil {
		t.Fatalf("GetFieldValue failed: %v", err)
	}
	if value != "test value" {
		t.Errorf("expected 'test value', got %v", value)
	}
}

func Test_GetConfigValue_Success(t *testing.T) {
	testValue := "test"
	initData := &GetSet_TestStruct{
		StringField: "hello",
		IntField:    42,
		BoolField:   true,
		FloatField:  3.14,
		PtrField:    &testValue,
		NestedStruct: struct {
			NestedString string `config:"nested_string" env:"NESTED_STRING"`
			NestedInt    int    `config:"nested_int" env:"NESTED_INT"`
		}{
			NestedString: "nested",
			NestedInt:    100,
		},
	}

	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	tests := []struct {
		name      string
		configKey string
		expected  interface{}
	}{
		{"string field", "string_field", "hello"},
		{"int field", "int_field", 42},
		{"bool field", "bool_field", true},
		{"float field", "float_field", 3.14},
		{"pointer field", "ptr_field", &testValue},
		{"nested string", "nested.nested_string", "nested"},
		{"nested int", "nested.nested_int", 100},
	}

	// test getting multiple values expecting success
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := cfg.GetConfigValue(tt.configKey)
			if err != nil {
				t.Fatalf("GetConfigValue failed: %v", err)
			}
			if value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, value)
			}
		})
	}
}

func Test_GetConfigValue_NonExistentKey(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// accessing a non existent key should return an error
	_, err = cfg.GetConfigValue("non_existent_key")
	if err == nil {
		t.Fatal("expected error for non-existent config key")
	}
	if !contains(err.Error(), "config key \"non_existent_key\" not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_SetConfigValue_Success(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	tests := []struct {
		name      string
		configKey string
		value     interface{}
		expected  interface{}
	}{
		{"string field", "string_field", "new value", "new value"},
		{"int field", "int_field", 123, 123},
		{"bool field", "bool_field", true, true},
		{"float field", "float_field", 2.71, 2.71},
		{"nested string", "nested.nested_string", "nested value", "nested value"},
		{"nested int", "nested.nested_int", 999, 999},
	}

	// test setting multiple values expecting success
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cfg.SetConfigValue(tt.configKey, tt.value)
			if err != nil {
				t.Fatalf("SetConfigValue failed: %v", err)
			}

			value, err := cfg.GetConfigValue(tt.configKey)
			if err != nil {
				t.Fatalf("GetConfigValue failed: %v", err)
			}
			if value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, value)
			}
		})
	}
}

func Test_SetConfigValue_TypeConversion(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// SetConfigValue should initialize the accessed path and convert int8 to int
	err = cfg.SetConfigValue("int_field", int8(-42))
	if err != nil {
		t.Fatalf("SetConfigValue failed: %v", err)
	}

	// GetConfigValue should now read the set value and get the correct type
	value, err := cfg.GetConfigValue("int_field")
	if err != nil {
		t.Fatalf("GetConfigValue failed: %v", err)
	}
	i, ok := value.(int)
	if !ok {
		t.Fatalf("expected int, got %T", value)
	}
	if i != -42 {
		t.Errorf("expected -42, got %v", i)
	}
}

func Test_SetConfigValue_IncompatibleType(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// trying to set a string value to an int field should return an error
	err = cfg.SetConfigValue("int_field", "not a number")
	if err == nil {
		t.Fatal("expected error for incompatible type")
	}
	if !contains(err.Error(), "cannot convert") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_SetConfigValue_NonExistentKey(t *testing.T) {
	initData := &GetSet_TestStruct{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// accessing a non existent key should return an error
	err = cfg.SetConfigValue("non_existent_key", "value")
	if err == nil {
		t.Fatal("expected error for non-existent config key")
	}
	if !contains(err.Error(), "config key \"non_existent_key\" not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_getFieldByPath_NilPointerInitializationWithSet(t *testing.T) {
	type DeepNested struct {
		Value string `config:"value"`
	}
	type Nested struct {
		Deep *DeepNested `config:"deep"`
	}
	type TestStruct2 struct {
		Nested *Nested `config:"nested"`
	}

	initData := &TestStruct2{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// SetFieldValue should initialize the accessed path
	err = cfg.SetFieldValue("Nested.Deep.Value", "test")
	if err != nil {
		t.Fatalf("SetFieldValue failed: %v", err)
	}

	if cfg.Data.Nested == nil {
		t.Fatal("expected Nested to be initialized")
	}
	if cfg.Data.Nested.Deep == nil {
		t.Fatal("expected Deep to be initialized")
	}
	if cfg.Data.Nested.Deep.Value != "test" {
		t.Errorf("expected 'test', got %v", cfg.Data.Nested.Deep.Value)
	}
}

func Test_getFieldByPath_NilPointerInitializationWithGet(t *testing.T) {
	type DeepNested struct {
		Value string `config:"value"`
	}
	type Nested struct {
		Deep *DeepNested `config:"deep"`
	}
	type TestStruct2 struct {
		Nested *Nested `config:"nested"`
	}

	initData := &TestStruct2{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// GetFieldValue should initialize the accessed path (and return at least the expected Zero Type)
	value, err := cfg.GetFieldValue("Nested.Deep.Value")
	if err != nil {
		t.Fatalf("GetFieldValue failed: %v", err)
	}
	if value != "" {
		t.Errorf("expected <zero value>, got %v", value)
	}

	if cfg.Data.Nested == nil {
		t.Fatal("expected Nested to be initialized")
	}
	if cfg.Data.Nested.Deep == nil {
		t.Fatal("expected Deep to be initialized")
	}
	if cfg.Data.Nested.Deep.Value != "" {
		t.Errorf("expected <zero value>, got %v", cfg.Data.Nested.Deep.Value)
	}
}

func Test_getFieldByPath_NilPointerInitializationWithSetGet(t *testing.T) {
	type DeepNested struct {
		Value string `config:"value"`
	}
	type Nested struct {
		Deep *DeepNested `config:"deep"`
	}
	type TestStruct2 struct {
		Nested *Nested `config:"nested"`
	}

	initData := &TestStruct2{}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// GetFieldValue should initialize the accessed path (and return at least the expected Zero Type)
	value, err := cfg.GetFieldValue("Nested")
	if err != nil {
		t.Fatalf("GetFieldValue failed: %v", err)
	}

	if reflect.ValueOf(value).IsNil() {
		t.Errorf("expected initialized struct, got nil")
	}

	// SetFieldValue should initialize the full accessed path
	err = cfg.SetFieldValue("Nested.Deep.Value", "test")
	if err != nil {
		t.Fatalf("SetFieldValue failed: %v", err)
	}

	// GetFieldValue should now read the set value
	value, err = cfg.GetFieldValue("Nested.Deep.Value")
	if err != nil {
		t.Fatalf("GetFieldValue failed: %v", err)
	}
	if value != "test" {
		t.Errorf("expected 'test', got %v", value)
	}

	if cfg.Data.Nested == nil {
		t.Fatal("expected Nested to be initialized")
	}
	if cfg.Data.Nested.Deep == nil {
		t.Fatal("expected Deep to be initialized")
	}
	if cfg.Data.Nested.Deep.Value != "test" {
		t.Errorf("expected 'test', got %v", cfg.Data.Nested.Deep.Value)
	}
}

func Test_getFieldByPath_InvalidStructType(t *testing.T) {
	type TestStruct3 struct {
		StringField string `config:"string_field"`
	}

	initData := &TestStruct3{StringField: "test"}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// access a leaf type as if it was a node type (struct/embedded) should return an error
	_, err = cfg.GetFieldValue("StringField.InvalidPath")
	if err == nil {
		t.Fatal("expected error for invalid path on non-struct field")
	}
	if !contains(err.Error(), "expected struct, got string") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func Test_getFieldByPath_UnexportedField(t *testing.T) {
	type TestStruct4 struct {
		ExportedField   string `config:"exported"`
		unexportedField string `config:"unexported"`
	}
	initData := &TestStruct4{
		unexportedField: "unexported value",
	}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// no access via GetFieldValue since the exported field can not be interfaced
	_, err = cfg.GetFieldValue("unexportedField")
	if err == nil {
		t.Fatal("expected error for unexported field")
	}
	if !contains(err.Error(), "field unexportedField cannot be accessed") {
		t.Errorf("unexpected error message: %v", err)
	}
	// direct access to unexported field still working
	if cfg.Data.unexportedField != "unexported value" {
		t.Errorf("expected 'unexported value', got %v", cfg.Data.unexportedField)
	}
}
