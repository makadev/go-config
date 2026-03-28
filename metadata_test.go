package config_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/makadev/go-config"
)

// contains checks if substr is in s.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func Test_MetadataInit(t *testing.T) {
	cfg, err := config.NewConfig(nil, &struct {
		somefiled string `config:"somefiled" env:"SOMEFILED"`
	}{})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if cfg.Metadata == nil {
		t.Fatal("expected Metadata to be initialized")
	}
	// FieldPathMap should be initialized map[string]*FieldInfo
	if cfg.Metadata.FieldPathMap == nil ||
		reflect.TypeOf(cfg.Metadata.FieldPathMap).Kind() != reflect.Map ||
		reflect.TypeOf(cfg.Metadata.FieldPathMap).Elem().Kind() != reflect.Ptr ||
		reflect.TypeOf(cfg.Metadata.FieldPathMap).Elem().Elem().Kind() != reflect.Struct ||
		reflect.TypeOf(cfg.Metadata.FieldPathMap).Elem().Elem().Name() != "FieldInfo" {
		t.Fatal("expected FieldPathMap to be initialized map[string]*FieldInfo")
	}
}

func Test_MetadataInit_Simple(t *testing.T) {
	initData := &struct {
		// simple fields
		Field1 string `config:"field1" env:"FIELD1"`
		Field2 int    `config:"field2" env:"FIELD2"`
		S1     struct {
			Field3 bool `config:"field3" env:"FIELD3"`
		} `config:"s1"`
		SliceField     []string `config:"slice_field" env:"SLICE_FIELD"`
		SliceOfStructs []struct {
			// this field won't be handled since it's inside a slice
			Field4 float64 `config:"field4"`
		} `config:"slice_of_structs"`
		ArrayField   [2]int            `config:"array_field" env:"ARRAY_FIELD"`
		MapField     map[string]string `config:"map_field" env:"MAP_FIELD"`
		MapOfStructs map[string]struct {
			// this field won't be handled since it's inside a map
			Field5 int `config:"field5"`
		} `config:"map_of_structs" env:"MAP_FIELD2"`
		// this field won't be handled since it's not exported
		ignoredField string `config:"ignored_field" env:"IGNORED_FIELD"`
	}{}

	cfg, err := config.NewConfig(nil, initData)

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if cfg.Metadata == nil {
		t.Fatal("expected Metadata to be initialized")
	}
	// FieldPathMap should be initialized map[string]*FieldInfo
	if cfg.Metadata.FieldPathMap == nil ||
		reflect.TypeOf(cfg.Metadata.FieldPathMap).Kind() != reflect.Map ||
		reflect.TypeOf(cfg.Metadata.FieldPathMap).Elem().Kind() != reflect.Ptr ||
		reflect.TypeOf(cfg.Metadata.FieldPathMap).Elem().Elem().Kind() != reflect.Struct ||
		reflect.TypeOf(cfg.Metadata.FieldPathMap).Elem().Elem().Name() != "FieldInfo" {
		t.Fatal("expected FieldPathMap to be initialized map[string]*FieldInfo")
	}

	// simple fields

	if cfg.Metadata.FieldPathMap["Field1"] == nil &&
		cfg.Metadata.FieldPathMap["Field1"].FieldPath != "Field1" {
		t.Fatal("expected Field1 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["Field1"].ConfigKey != "field1" &&
		cfg.Metadata.FieldPathMap["Field1"].ConfigName != "field1" {
		t.Fatal("expected Field1.ConfigKey to be 'field1'")
	}

	if cfg.Metadata.FieldPathMap["Field2"] == nil {
		t.Fatal("expected Field2 to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["S1.Field3"] == nil &&
		cfg.Metadata.FieldPathMap["S1.Field3"].FieldPath != "S1.Field3" {
		t.Fatal("expected S1.Field3 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["S1.Field3"].ConfigKey != "s1.field3" &&
		cfg.Metadata.FieldPathMap["S1.Field3"].ConfigName != "field3" {
		t.Fatal("expected S1.Field3.ConfigKey to be 's1.field3'")
	}

	if cfg.Metadata.FieldPathMap["SliceField"] == nil {
		t.Fatal("expected SliceField to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["SliceOfStructs"] == nil {
		t.Fatal("expected SliceOfStructs to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["SliceOfStructs.Field4"] != nil {
		t.Fatal("expected SliceOfStructs.Field4 to not be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["ArrayField"] == nil {
		t.Fatal("expected ArrayField to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["MapField"] == nil {
		t.Fatal("expected MapField to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["MapOfStructs"] == nil {
		t.Fatal("expected MapOfStructs to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["MapOfStructs.Field5"] != nil {
		t.Fatal("expected MapOfStructs.Field5 to not be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["ignoredField"] != nil {
		t.Fatal("expected ignoredField to not be initialized in FieldPathMap")
	}
}

func Test_MetadataInit_Pointers(t *testing.T) {
	initData := &struct {
		// and now with pointers
		PtrField1 *string `config:"ptr_field1" env:"PFIELD1"`
		PtrField2 *int    `config:"ptr_field2" env:"PFIELD2"`
		PtrS1     *struct {
			Field3 bool `config:"field3" env:"PFIELD3"`
		} `config:"ptr_s1"`
		PtrSliceField     *[]string `config:"ptr_slice_field" env:"PSLICE_FIELD"`
		PtrSliceOfStructs *[]struct {
			// this field won't be handled since it's inside a slice
			Field4 float64 `config:"field4"`
		} `config:"ptr_slice_of_structs"`
		PtrArrayField   *[2]int             `config:"ptr_array_field" env:"PARARRAY_FIELD"`
		PtrMapField     *map[string]string  `config:"ptr_map_field" env:"PMAP_FIELD"`
		PtrMapPtrField  *map[string]*string `config:"ptr_map_ptr_field" env:"PMAP_PTR_FIELD"`
		PtrMapOfStructs *map[string]struct {
			// this field won't be handled since it's inside a map
			Field5 int `config:"field5"`
		} `config:"ptr_map_of_structs"`
		ptrIgnoredField *string `config:"ptr_ignored_field" env:"PIGNORED_FIELD"`
	}{}

	cfg, err := config.NewConfig(nil, initData)

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// more pointers

	if cfg.Metadata.FieldPathMap["PtrField1"] == nil {
		t.Fatal("expected PtrField1 to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrField2"] == nil {
		t.Fatal("expected PtrField2 to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrS1.Field3"] == nil {
		t.Fatal("expected PtrS1.Field3 to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrSliceField"] == nil {
		t.Fatal("expected PtrSliceField to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrSliceOfStructs"] == nil {
		t.Fatal("expected PtrSliceOfStructs to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrSliceOfStructs.Field4"] != nil {
		t.Fatal("expected PtrSliceOfStructs.Field4 to not be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrArrayField"] == nil {
		t.Fatal("expected PtrArrayField to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrMapField"] == nil {
		t.Fatal("expected PtrMapField to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrMapPtrField"] == nil {
		t.Fatal("expected PtrMapPtrField to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrMapOfStructs"] == nil {
		t.Fatal("expected PtrMapOfStructs to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrMapOfStructs.Field5"] != nil {
		t.Fatal("expected PtrMapOfStructs.Field5 to not be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["ptrIgnoredField"] != nil {
		t.Fatal("expected ptrIgnoredField to not be initialized in FieldPathMap")
	}

}

func Test_MetadataInit_Deep(t *testing.T) {
	initData := &struct {
		// pointers to pointers to pointers ...
		DeepPtrSliceOfStructs ****[]struct {
			// this field won't be handled since it's pointer chained
			Field4 float64 `config:"field4"`
		} `config:"deep_ptr_slice_of_structs"`

		// deep type
		PtrMapOfArrayOfStructs *map[string][]struct {
			Field5 int `config:"field5"`
		} `config:"ptr_map_of_structs2"`
	}{}

	cfg, err := config.NewConfig(nil, initData)

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	// pointers to pointers to pointers ...

	if cfg.Metadata.FieldPathMap["DeepPtrSliceOfStructs"] == nil {
		t.Fatal("expected DeepPtrSliceOfStructs to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["DeepPtrSliceOfStructs.Field4"] != nil {
		t.Fatal("expected DeepPtrSliceOfStructs.Field4 to not be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrMapOfArrayOfStructs"] == nil {
		t.Fatal("expected PtrMapOfArrayOfStructs to be initialized in FieldPathMap")
	}

	if cfg.Metadata.FieldPathMap["PtrMapOfArrayOfStructs.Field5"] != nil {
		t.Fatal("expected PtrMapOfArrayOfStructs.Field5 to not be initialized in FieldPathMap")
	}
}

func Test_MetadataInit_AutoEnv(t *testing.T) {
	initData := &struct {
		Field1 string  `config:"field1"`
		Field2 int     `config:"field2"`
		Field3 bool    `config:"field3"`
		Field4 float64 `config:"field4" env:"SPECIAL_NAME"`
	}{}

	opts := config.NewOptions()
	opts.AutoEnv = true
	cfg, err := config.NewConfig(opts, initData)

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if cfg.Metadata.FieldPathMap["Field1"] == nil {
		t.Fatal("expected Field1 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["Field1"].EnvVar != "FIELD1" {
		t.Fatalf("expected Field1 to have an environment variable %v, got %v", "FIELD1", cfg.Metadata.FieldPathMap["Field1"].EnvVar)
	}

	if cfg.Metadata.FieldPathMap["Field2"] == nil {
		t.Fatal("expected Field2 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["Field2"].EnvVar != "FIELD2" {
		t.Fatalf("expected Field2 to have an environment variable %v, got %v", "FIELD2", cfg.Metadata.FieldPathMap["Field2"].EnvVar)
	}

	if cfg.Metadata.FieldPathMap["Field3"] == nil {
		t.Fatal("expected Field3 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["Field3"].EnvVar != "FIELD3" {
		t.Fatalf("expected Field3 to have an environment variable %v, got %v", "FIELD3", cfg.Metadata.FieldPathMap["Field3"].EnvVar)
	}

	if cfg.Metadata.FieldPathMap["Field4"] == nil {
		t.Fatal("expected Field4 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["Field4"].EnvVar != "SPECIAL_NAME" {
		t.Fatalf("expected Field4 to have an environment variable %v, got %v", "SPECIAL_NAME", cfg.Metadata.FieldPathMap["Field4"].EnvVar)
	}
}

func Test_MetadataInit_SkipEnv(t *testing.T) {
	initData := &struct {
		Field1 string  `config:"field1"`
		Field2 float64 `config:"field2" env:"SPECIAL_NAME"`
	}{}

	opts := config.NewOptions()
	opts.SkipEnv = true
	cfg, err := config.NewConfig(opts, initData)

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if cfg.Metadata.FieldPathMap["Field1"] == nil {
		t.Fatal("expected Field1 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["Field1"].EnvVar != "" {
		t.Fatalf("expected Field1 to have an environment variable %v, got %v", "", cfg.Metadata.FieldPathMap["Field1"].EnvVar)
	}

	if cfg.Metadata.FieldPathMap["Field2"] == nil {
		t.Fatal("expected Field2 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["Field2"].EnvVar != "" {
		t.Fatalf("expected Field2 to have an environment variable %v, got %v", "", cfg.Metadata.FieldPathMap["Field2"].EnvVar)
	}
}

func Test_MetadataInit_DuplicateConfigKey(t *testing.T) {
	initData := &struct {
		Field1 string `config:"field1"`
		Field2 string `config:"field1"`
	}{}

	_, err := config.NewConfig(nil, initData)

	if err == nil {
		t.Fatalf("expected error due to duplicate config key, got nil")
	}
	expectedErr := "duplicate config key"
	if err != nil && !contains(err.Error(), expectedErr) {
		t.Fatalf("expected error to contain %q, got %v", expectedErr, err)
	}
}

func Test_MetadataInit_DuplicateEnvVar(t *testing.T) {
	initData := &struct {
		Field1 string `config:"field1" env:"FIELD1"`
		Field2 string `config:"field2" env:"FIELD1"`
	}{}

	_, err := config.NewConfig(nil, initData)

	if err == nil {
		t.Fatalf("expected error due to duplicate env var, got nil")
	}
	expectedErr := "duplicate env var"
	if err != nil && !contains(err.Error(), expectedErr) {
		t.Fatalf("expected error to contain %q, got %v", expectedErr, err)
	}
}

func Test_MetadataInit_InvalidConfigName(t *testing.T) {
	// A config tag with an invalid character (space) must be rejected.
	_, err := config.NewConfig(nil, &struct {
		Field string `config:"in valid"`
	}{})
	if err == nil {
		t.Fatal("expected error for invalid config name")
	}
	if !contains(err.Error(), "must match regex") {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_MetadataInit_InvalidEnvVar(t *testing.T) {
	// An env tag starting with a digit is not a valid identifier.
	_, err := config.NewConfig(nil, &struct {
		Field string `config:"field" env:"1INVALID"`
	}{})
	if err == nil {
		t.Fatal("expected error for invalid env var name")
	}
	if !contains(err.Error(), "must match regex") {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_MetadataInit_NestedStructPropagatesError(t *testing.T) {
	// Verify that an invalid config name inside a nested struct propagates up
	// (covers the `return err` inside the reflect.Struct recursion case).
	_, err := config.NewConfig(nil, &struct {
		Inner struct {
			BadField string `config:"bad field"`
		} `config:"inner"`
	}{})
	if err == nil {
		t.Fatal("expected error propagated from nested struct")
	}
	if !contains(err.Error(), "must match regex") {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_MetadataInit_NestedPtrStructPropagatesError(t *testing.T) {
	// Verify that an invalid config name inside a pointer-to-struct propagates up
	// (covers the `return err` inside the reflect.Ptr recursion case).
	_, err := config.NewConfig(nil, &struct {
		Inner *struct {
			BadField string `config:"bad field"`
		} `config:"inner"`
	}{})
	if err == nil {
		t.Fatal("expected error propagated from pointer-to-struct")
	}
	if !contains(err.Error(), "must match regex") {
		t.Errorf("unexpected error: %v", err)
	}
}

func Test_MetadataInit_FallbackNaming(t *testing.T) {
	initData := &struct {
		Field1 string
		Field2 string
	}{}

	cfg, err := config.NewConfig(nil, initData)

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if cfg.Metadata.FieldPathMap["Field1"] == nil &&
		cfg.Metadata.FieldPathMap["Field1"].FieldPath != "Field1" {
		t.Fatal("expected Field1 to be initialized in FieldPathMap")
	}
	if cfg.Metadata.FieldPathMap["Field1"].ConfigKey != "field1" &&
		cfg.Metadata.FieldPathMap["Field1"].ConfigName != "field1" {
		t.Fatal("expected Field1.ConfigKey to be 'field1'")
	}
}
