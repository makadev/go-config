package config_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/makadev/go-config"
)

type Env_TestConfig struct {
	StringField   string            `env:"STRING_VAR"`
	IntField      int               `env:"INT_VAR"`
	BoolField     bool              `env:"BOOL_VAR"`
	FloatField    float64           `env:"FLOAT_VAR"`
	DurationField time.Duration     `env:"DURATION_VAR"`
	SliceField    []string          `env:"SLICE_VAR"`
	MapField      map[string]string `env:"MAP_VAR"`

	// Pointer fields
	StringPtr *string `env:"STRING_PTR_VAR"`
	IntPtr    *int    `env:"INT_PTR_VAR"`

	// Nested struct without env tag (should be ignored)
	NestedField struct {
		Value string
	}

	// Field without env tag (should be ignored)
	NoEnvField string
}

func Test_loadFromEnv(t *testing.T) {

	ptrString := "ptr_value"
	ptrInt := 100

	tests := []struct {
		name      string
		envVars   map[string]string
		envPrefix string
		expected  Env_TestConfig
	}{
		{
			name: "basic env loading",
			envVars: map[string]string{
				"STRING_VAR":     "test_string",
				"INT_VAR":        "42",
				"BOOL_VAR":       "true",
				"FLOAT_VAR":      "3.14",
				"DURATION_VAR":   "5m",
				"SLICE_VAR":      "a,b,c",
				"MAP_VAR":        "key1=val1,key2=val2",
				"STRING_PTR_VAR": "ptr_value",
				"INT_PTR_VAR":    "100",
			},
			expected: Env_TestConfig{
				StringField:   "test_string",
				IntField:      42,
				BoolField:     true,
				FloatField:    3.14,
				DurationField: 5 * time.Minute,
				SliceField:    []string{"a", "b", "c"},
				MapField:      map[string]string{"key1": "val1", "key2": "val2"},
				StringPtr:     &ptrString,
				IntPtr:        &ptrInt,
			},
		},
		{
			name: "with env prefix",
			envVars: map[string]string{
				"MYAPP_STRING_VAR": "prefixed_string",
				"MYAPP_INT_VAR":    "99",
			},
			envPrefix: "MYAPP_",
			expected: Env_TestConfig{
				StringField: "prefixed_string",
				IntField:    99,
			},
		},
		{
			name: "partial env vars",
			envVars: map[string]string{
				"STRING_VAR": "only_string",
			},
			expected: Env_TestConfig{
				StringField: "only_string",
			},
		},
		{
			name:     "no env vars",
			envVars:  map[string]string{},
			expected: Env_TestConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			// Create config
			configData := &Env_TestConfig{}
			opts := config.NewOptions()
			if tt.envPrefix != "" {
				opts.EnvPrefix = tt.envPrefix
			}

			cfg, err := config.NewConfig(opts, configData)
			if err != nil {
				t.Fatalf("failed to create config: %v", err)
			}

			// Load from environment (through Load() method with SkipFiles option)
			cfg.Options.SkipFiles = true
			err = cfg.Load()
			if err != nil {
				t.Fatalf("failed to load from env: %v", err)
			}

			// Compare results
			if configData.StringField != tt.expected.StringField {
				t.Errorf("StringField: expected %q, got %q", tt.expected.StringField, configData.StringField)
			}
			if configData.IntField != tt.expected.IntField {
				t.Errorf("IntField: expected %d, got %d", tt.expected.IntField, configData.IntField)
			}
			if configData.BoolField != tt.expected.BoolField {
				t.Errorf("BoolField: expected %t, got %t", tt.expected.BoolField, configData.BoolField)
			}
			if configData.FloatField != tt.expected.FloatField {
				t.Errorf("FloatField: expected %f, got %f", tt.expected.FloatField, configData.FloatField)
			}
			if configData.DurationField != tt.expected.DurationField {
				t.Errorf("DurationField: expected %v, got %v", tt.expected.DurationField, configData.DurationField)
			}
			if !reflect.DeepEqual(configData.SliceField, tt.expected.SliceField) {
				t.Errorf("SliceField: expected %v, got %v", tt.expected.SliceField, configData.SliceField)
			}
			if !reflect.DeepEqual(configData.MapField, tt.expected.MapField) {
				t.Errorf("MapField: expected %v, got %v", tt.expected.MapField, configData.MapField)
			}

			// Check pointer fields
			if !reflect.DeepEqual(configData.StringPtr, tt.expected.StringPtr) {
				t.Errorf("StringPtr: expected %v, got %v", ptrValue(tt.expected.StringPtr), ptrValue(configData.StringPtr))
			}
			if !reflect.DeepEqual(configData.IntPtr, tt.expected.IntPtr) {
				t.Errorf("IntPtr: expected %v, got %v", ptrValue(tt.expected.IntPtr), ptrValue(configData.IntPtr))
			}
		})
	}
}

func ptrValue(ptr interface{}) interface{} {
	if ptr == nil {
		return nil
	}
	v := reflect.ValueOf(ptr)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		return v.Elem().Interface()
	}
	return nil
}

func Test_loadFromEnv_errors(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "invalid int",
			envVars: map[string]string{
				"INT_VAR": "not_an_int",
			},
			wantErr: true,
		},
		{
			name: "invalid bool",
			envVars: map[string]string{
				"BOOL_VAR": "maybe",
			},
			wantErr: true,
		},
		{
			name: "invalid float",
			envVars: map[string]string{
				"FLOAT_VAR": "not_a_float",
			},
			wantErr: true,
		},
		{
			name: "invalid duration",
			envVars: map[string]string{
				"DURATION_VAR": "not_a_duration",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			// Create config
			configData := &Env_TestConfig{}
			cfg, err := config.NewConfig(nil, configData)
			if err != nil {
				t.Fatalf("failed to create config: %v", err)
			}

			// Load from environment (through Load() method with SkipFiles option)
			cfg.Options.SkipFiles = true
			err = cfg.Load()
			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func Test_loadFromEnv_autoenv(t *testing.T) {
	type NestedConfig struct {
		Host string
		Port int
	}

	type AppConfig struct {
		Name   string
		Env    string
		Nested NestedConfig
	}

	// Set test environment variables
	envVars := map[string]string{
		"APP_NAME":        "TestApp",
		"APP_ENV":         "development",
		"APP_NESTED_HOST": "127.0.0.1",
		"APP_NESTED_PORT": "8080",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	// Create config
	configData := &AppConfig{}
	opts := config.NewOptions()
	opts.EnvPrefix = "APP_"
	opts.AutoEnv = true
	cfg, err := config.NewConfig(opts, configData)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// Load from environment
	cfg.Options.SkipFiles = true
	err = cfg.Load()
	if err != nil {
		t.Fatalf("failed to load from env: %v", err)
	}

	// Compare results
	if configData.Name != "TestApp" {
		t.Errorf("Name: expected %q, got %q", "TestApp", configData.Name)
	}
	if configData.Env != "development" {
		t.Errorf("Env: expected %q, got %q", "development", configData.Env)
	}
	if configData.Nested.Host != "127.0.0.1" {
		t.Errorf("Nested.Host: expected %q, got %q", "127.0.0.1", configData.Nested.Host)
	}
	if configData.Nested.Port != 8080 {
		t.Errorf("Nested.Port: expected %d, got %d", 8080, configData.Nested.Port)
	}
}

// Additional test cases for different data types through the main interface
func Test_loadFromEnv_different_types(t *testing.T) {
	type TypeTestConfig struct {
		// Numeric types
		Int8Field    int8    `env:"INT8_VAR"`
		Int16Field   int16   `env:"INT16_VAR"`
		Int32Field   int32   `env:"INT32_VAR"`
		Int64Field   int64   `env:"INT64_VAR"`
		UintField    uint    `env:"UINT_VAR"`
		Uint8Field   uint8   `env:"UINT8_VAR"`
		Uint16Field  uint16  `env:"UINT16_VAR"`
		Uint32Field  uint32  `env:"UINT32_VAR"`
		Uint64Field  uint64  `env:"UINT64_VAR"`
		Float32Field float32 `env:"FLOAT32_VAR"`
		Float64Field float64 `env:"FLOAT64_VAR"`

		// Complex types
		IntSliceField  []int             `env:"INT_SLICE_VAR"`
		StringMapField map[string]string `env:"STRING_MAP_VAR"`
		IntMapField    map[string]int    `env:"INT_MAP_VAR"`
		BoolMapField   map[string]bool   `env:"BOOL_MAP_VAR"`

		// Boolean variants
		BoolTrueField1  bool `env:"BOOL_TRUE1"`
		BoolTrueField2  bool `env:"BOOL_TRUE2"`
		BoolTrueField3  bool `env:"BOOL_TRUE3"`
		BoolFalseField1 bool `env:"BOOL_FALSE1"`
		BoolFalseField2 bool `env:"BOOL_FALSE2"`
	}

	// Set test environment variables
	envVars := map[string]string{
		"INT8_VAR":       "127",
		"INT16_VAR":      "32767",
		"INT32_VAR":      "2147483647",
		"INT64_VAR":      "9223372036854775807",
		"UINT_VAR":       "100",
		"UINT8_VAR":      "255",
		"UINT16_VAR":     "65535",
		"UINT32_VAR":     "4294967295",
		"UINT64_VAR":     "18446744073709551615",
		"FLOAT32_VAR":    "3.14",
		"FLOAT64_VAR":    "2.718281828",
		"INT_SLICE_VAR":  "1,2,3,4,5",
		"STRING_MAP_VAR": "key1=value1,key2=value2",
		"INT_MAP_VAR":    "count=42,total=100",
		"BOOL_MAP_VAR":   "enabled=true,debug=false",
		"BOOL_TRUE1":     "true",
		"BOOL_TRUE2":     "yes",
		"BOOL_TRUE3":     "1",
		"BOOL_FALSE1":    "false",
		"BOOL_FALSE2":    "no",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	// Create config
	configData := &TypeTestConfig{}
	cfg, err := config.NewConfig(nil, configData)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// Load from environment
	cfg.Options.SkipFiles = true
	err = cfg.Load()
	if err != nil {
		t.Fatalf("failed to load from env: %v", err)
	}

	// Verify numeric types
	if configData.Int8Field != 127 {
		t.Errorf("Int8Field: expected 127, got %d", configData.Int8Field)
	}
	if configData.Int16Field != 32767 {
		t.Errorf("Int16Field: expected 32767, got %d", configData.Int16Field)
	}
	if configData.Int32Field != 2147483647 {
		t.Errorf("Int32Field: expected 2147483647, got %d", configData.Int32Field)
	}
	if configData.Int64Field != 9223372036854775807 {
		t.Errorf("Int64Field: expected 9223372036854775807, got %d", configData.Int64Field)
	}
	if configData.UintField != 100 {
		t.Errorf("UintField: expected 100, got %d", configData.UintField)
	}
	if configData.Uint8Field != 255 {
		t.Errorf("Uint8Field: expected 255, got %d", configData.Uint8Field)
	}
	if configData.Uint16Field != 65535 {
		t.Errorf("Uint16Field: expected 65535, got %d", configData.Uint16Field)
	}
	if configData.Uint32Field != 4294967295 {
		t.Errorf("Uint32Field: expected 4294967295, got %d", configData.Uint32Field)
	}
	if configData.Uint64Field != 18446744073709551615 {
		t.Errorf("Uint64Field: expected 18446744073709551615, got %d", configData.Uint64Field)
	}
	if configData.Float32Field != 3.14 {
		t.Errorf("Float32Field: expected 3.14, got %f", configData.Float32Field)
	}
	if configData.Float64Field != 2.718281828 {
		t.Errorf("Float64Field: expected 2.718281828, got %f", configData.Float64Field)
	}

	// Verify complex types
	expectedIntSlice := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(configData.IntSliceField, expectedIntSlice) {
		t.Errorf("IntSliceField: expected %v, got %v", expectedIntSlice, configData.IntSliceField)
	}

	expectedStringMap := map[string]string{"key1": "value1", "key2": "value2"}
	if !reflect.DeepEqual(configData.StringMapField, expectedStringMap) {
		t.Errorf("StringMapField: expected %v, got %v", expectedStringMap, configData.StringMapField)
	}

	expectedIntMap := map[string]int{"count": 42, "total": 100}
	if !reflect.DeepEqual(configData.IntMapField, expectedIntMap) {
		t.Errorf("IntMapField: expected %v, got %v", expectedIntMap, configData.IntMapField)
	}

	expectedBoolMap := map[string]bool{"enabled": true, "debug": false}
	if !reflect.DeepEqual(configData.BoolMapField, expectedBoolMap) {
		t.Errorf("BoolMapField: expected %v, got %v", expectedBoolMap, configData.BoolMapField)
	}

	// Verify boolean variants
	if !configData.BoolTrueField1 {
		t.Error("BoolTrueField1: expected true, got false")
	}
	if !configData.BoolTrueField2 {
		t.Error("BoolTrueField2: expected true, got false")
	}
	if !configData.BoolTrueField3 {
		t.Error("BoolTrueField3: expected true, got false")
	}
	if configData.BoolFalseField1 {
		t.Error("BoolFalseField1: expected false, got true")
	}
	if configData.BoolFalseField2 {
		t.Error("BoolFalseField2: expected false, got true")
	}
}

func Test_loadFromEnv_Map(t *testing.T) {
	type TypeTestConfig struct {
		StringMapField   map[string]string `env:"APP_STRING_MAP_VAR"`
		IntMapField      map[string]int    `env:"APP_INT_MAP_VAR"`
		BoolMapField     map[string]bool   `env:"APP_BOOL_MAP_VAR"`
		EmptyStringField map[string]string `env:"APP_EMPTY_STRING_MAP_VAR"`
		EmptyMapField    map[string]string `env:"APP_EMPTY_MAP_VAR"`
	}
	envVars := map[string]string{
		"APP_STRING_MAP_VAR":       "key1=value1,key2=value2",
		"APP_INT_MAP_VAR":          "count=42,total=100",
		"APP_BOOL_MAP_VAR":         "enabled=true,debug=false",
		"APP_EMPTY_STRING_MAP_VAR": "test=,test2=",
		"APP_EMPTY_MAP_VAR":        "",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	// Create config
	configData := &TypeTestConfig{}
	cfg, err := config.NewConfig(nil, configData)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// Load from environment
	cfg.Options.SkipFiles = true
	err = cfg.Load()
	if err != nil {
		t.Fatalf("failed to load from env: %v", err)
	}

	// Verify complex types
	expectedStringMap := map[string]string{"key1": "value1", "key2": "value2"}
	if !reflect.DeepEqual(configData.StringMapField, expectedStringMap) {
		t.Errorf("StringMapField: expected %v, got %v", expectedStringMap, configData.StringMapField)
	}

	expectedIntMap := map[string]int{"count": 42, "total": 100}
	if !reflect.DeepEqual(configData.IntMapField, expectedIntMap) {
		t.Errorf("IntMapField: expected %v, got %v", expectedIntMap, configData.IntMapField)
	}

	expectedBoolMap := map[string]bool{"enabled": true, "debug": false}
	if !reflect.DeepEqual(configData.BoolMapField, expectedBoolMap) {
		t.Errorf("BoolMapField: expected %v, got %v", expectedBoolMap, configData.BoolMapField)
	}

	expectedEmptyStringMap := map[string]string{"test": "", "test2": ""}
	if !reflect.DeepEqual(configData.EmptyStringField, expectedEmptyStringMap) {
		t.Errorf("EmptyStringField: expected %v, got %v", expectedEmptyStringMap, configData.EmptyStringField)
	}

	expectedEmptyMap := map[string]string{}
	if !reflect.DeepEqual(configData.EmptyMapField, expectedEmptyMap) {
		t.Errorf("EmptyMapField: expected %v, got %v", expectedEmptyMap, configData.EmptyMapField)
	}
}

func Test_loadFromEnv_KeyNoString(t *testing.T) {
	type TypeTestConfig struct {
		StringMapField map[int]string `env:"APP_STRING_MAP_VAR"`
	}
	envVars := map[string]string{
		"APP_STRING_MAP_VAR": "1=value1,2=value2",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	// Create config
	configData := &TypeTestConfig{}
	cfg, err := config.NewConfig(nil, configData)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// Load from environment
	cfg.Options.SkipFiles = true
	err = cfg.Load()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_loadFromEnv_Map_error(t *testing.T) {
	type TypeTestConfig struct {
		EmptyIntMapField map[string]int `env:"APP_INT_MAP_VAR"`
	}
	envVars := map[string]string{
		"APP_INT_MAP_VAR": "count=",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	// Create config
	configData := &TypeTestConfig{}
	cfg, err := config.NewConfig(nil, configData)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// Load from environment
	cfg.Options.SkipFiles = true
	err = cfg.Load()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func Test_loadFromEnv_Slice(t *testing.T) {
	type TypeTestConfig struct {
		StringSliceField []string `env:"APP_STRING_SLICE_VAR"`
		IntSliceField    []int    `env:"APP_INT_SLICE_VAR"`
		BoolSliceField   []bool   `env:"APP_BOOL_SLICE_VAR"`
		EmptyStringField []string `env:"APP_EMPTY_STRING_SLICE_VAR"`
	}
	envVars := map[string]string{
		"APP_STRING_SLICE_VAR":       "value1,value2,value3",
		"APP_INT_SLICE_VAR":          "1,2,3",
		"APP_BOOL_SLICE_VAR":         "true,false,true",
		"APP_EMPTY_STRING_SLICE_VAR": "",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	// Create config
	configData := &TypeTestConfig{}
	cfg, err := config.NewConfig(nil, configData)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// Load from environment
	cfg.Options.SkipFiles = true
	err = cfg.Load()
	if err != nil {
		t.Fatalf("failed to load from env: %v", err)
	}

	// Verify slice types
	expectedStringSlice := []string{"value1", "value2", "value3"}
	if !reflect.DeepEqual(configData.StringSliceField, expectedStringSlice) {
		t.Errorf("StringSliceField: expected %v, got %v", expectedStringSlice, configData.StringSliceField)
	}

	expectedIntSlice := []int{1, 2, 3}
	if !reflect.DeepEqual(configData.IntSliceField, expectedIntSlice) {
		t.Errorf("IntSliceField: expected %v, got %v", expectedIntSlice, configData.IntSliceField)
	}

	expectedBoolSlice := []bool{true, false, true}
	if !reflect.DeepEqual(configData.BoolSliceField, expectedBoolSlice) {
		t.Errorf("BoolSliceField: expected %v, got %v", expectedBoolSlice, configData.BoolSliceField)
	}

	expectedEmptySlice := []string{}
	if !reflect.DeepEqual(configData.EmptyStringField, expectedEmptySlice) {
		t.Errorf("EmptyStringField: expected %v, got %v", expectedEmptySlice, configData.EmptyStringField)
	}
}

func Test_loadFromEnv_NestedStruct(t *testing.T) {
	type NestedConfig struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	}

	type TypeTestConfig struct {
		NestedField  NestedConfig
		NestedField2 NestedConfig `env:"NESTED_"`
	}

	envVars := map[string]string{
		"APP_HOST":        "example.com",
		"APP_PORT":        "9999",
		"APP_NESTED_HOST": "localhost",
		"APP_NESTED_PORT": "8080",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	// Create config
	opts := config.NewOptions()
	opts.EnvPrefix = "APP_"
	configData := &TypeTestConfig{}
	cfg, err := config.NewConfig(opts, configData)
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// Load from environment
	cfg.Options.SkipFiles = true
	err = cfg.Load()
	if err != nil {
		t.Fatalf("failed to load from env: %v", err)
	}

	// Verify nested struct fields
	expectedNested := NestedConfig{
		Host: "example.com",
		Port: 9999,
	}
	if !reflect.DeepEqual(configData.NestedField, expectedNested) {
		t.Errorf("NestedField: expected %v, got %v", expectedNested, configData.NestedField)
	}

	expectedNested2 := NestedConfig{
		Host: "localhost",
		Port: 8080,
	}
	if !reflect.DeepEqual(configData.NestedField2, expectedNested2) {
		t.Errorf("NestedField2: expected %v, got %v", expectedNested2, configData.NestedField2)
	}
}
