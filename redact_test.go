package config_test

import (
	"testing"

	"github.com/makadev/go-config"
)

type RedactNested struct {
	Password string `secret:"true"`
	Token    string
}

type RedactTestConfig struct {
	Username   string
	Password   string   `secret:"true"`
	Keys       []string `secret:"true"`
	Seeds      []int64  `secret:"true"`
	Nested     RedactNested
	Ptr        *RedactNested
	Ptr2       *RedactNested
	SPtr       *string           `secret:"true"`
	SecureBool bool              `secret:"true"`
	Meta       map[string]string `secret:"true"`
	Count      int
}

func TestRedactedCopy_Basic(t *testing.T) {
	orig := &RedactTestConfig{
		Username: "user1",
		Password: "secret",
		Keys:     []string{"key1", "key2"},
		Seeds:    []int64{1, 2},
		Nested:   RedactNested{Password: "nestedsecret", Token: "tok"},
		Ptr:      &RedactNested{Password: "ptrsecret", Token: "ptrtok"},
		Ptr2:     nil,
		Meta:     map[string]string{"k": "v"},
		Count:    42,
	}
	metadata, err := config.GetFieldInfoMap(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	redacted, err := config.RedactedCopy(orig, metadata, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cpy := redacted.(*RedactTestConfig)
	if cpy.Password != "***" {
		t.Errorf("Password not redacted: got %v", cpy.Password)
	}
	if len(cpy.Keys) != 1 || cpy.Keys[0] != "***" {
		t.Errorf("Keys not redacted: got %v", cpy.Keys)
	}
	if len(cpy.Seeds) > 0 {
		t.Errorf("Seeds not redacted: got %v", cpy.Seeds)
	}
	if cpy.Nested.Password != "***" {
		t.Errorf("Nested.Password not redacted: got %v", cpy.Nested.Password)
	}
	if cpy.Nested.Token != "tok" {
		t.Errorf("Nested.Token not redacted: got %v", cpy.Nested.Token)
	}
	if cpy.Ptr == nil {
		t.Errorf("Ptr not copied: got %v", cpy.Ptr)
	}
	if cpy.Ptr.Password != "***" {
		t.Errorf("Ptr.Password not redacted: got %v", cpy.Ptr)
	}
	if cpy.Ptr2 != nil {
		t.Errorf("Ptr2 initialized despite being nil: got %v", cpy.Ptr2)
	}
	if cpy.SPtr != nil && *cpy.SPtr != "***" {
		t.Errorf("SPtr not redacted: got %v", cpy.SPtr)
	}
	if cpy.SecureBool != false {
		t.Errorf("SecureBool not redacted: got %v", cpy.SecureBool)
	}
	if len(cpy.Meta) != 0 {
		t.Errorf("Meta not redacted: got %v", cpy.Meta)
	}
	if cpy.Username != orig.Username {
		t.Errorf("Username changed: got %v", cpy.Username)
	}
	if cpy.Count != orig.Count {
		t.Errorf("Count changed: got %v", cpy.Count)
	}
}

func TestRedactedCopy_DefaultRedactWith(t *testing.T) {
	orig := &RedactTestConfig{Password: "secret"}
	metadata, err := config.GetFieldInfoMap(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	redacted, err := config.RedactedCopy(orig, metadata, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cpy := redacted.(*RedactTestConfig)
	if cpy.Password != "..." {
		t.Errorf("Default redacted value not used: got %v", cpy.Password)
	}
}

func TestRedactedCopy_NoSecrets(t *testing.T) {
	orig := &RedactTestConfig{Username: "user1", Count: 5}
	metadata, err := config.GetFieldInfoMap(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	redacted, err := config.RedactedCopy(orig, metadata, "xxx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cpy := redacted.(*RedactTestConfig)
	if cpy.Username != orig.Username || cpy.Count != orig.Count {
		t.Errorf("Fields changed unexpectedly: got %+v", cpy)
	}
}

func TestRedactedCopy_InvalidInput(t *testing.T) {
	_, err := config.RedactedCopy(nil, nil, "")
	if err == nil {
		t.Error("expected error for nil input")
	}
	_, err = config.RedactedCopy(RedactTestConfig{}, nil, "")
	if err == nil {
		t.Error("expected error for non-pointer input")
	}
}
