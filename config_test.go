package config_test

import (
	"testing"

	"github.com/makadev/go-config"
)

type NewConfig_TestConfig struct {
	AppName string
	AppUrl  string
	Server  struct {
		Host string
		Port int
	}
}

func Test_NewConfig(t *testing.T) {
	initData := &NewConfig_TestConfig{
		AppName: "TestApp",
		AppUrl:  "http://127.0.0.1",
		Server: struct {
			Host string
			Port int
		}{
			Host: "127.0.0.1",
			Port: 5000,
		},
	}
	cfg, err := config.NewConfig(nil, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if cfg == nil {
		t.Fatal("failed to create config")
	}

	if cfg.Data != initData {
		t.Errorf("expected config data to be %+v, got %+v", initData, cfg.Data)
	}
}

func Test_NewConfigWithOpts(t *testing.T) {
	initData := &NewConfig_TestConfig{
		AppName: "TestApp",
		AppUrl:  "http://127.0.0.1",
		Server: struct {
			Host string
			Port int
		}{
			Host: "127.0.0.1",
			Port: 5000,
		},
	}
	opts := config.NewOptions()
	cfg, err := config.NewConfig(opts, initData)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if cfg == nil {
		t.Fatal("failed to create config")
	}

	if cfg.Options != opts {
		t.Errorf("expected config options to be %+v, got %+v", opts, cfg.Options)
	}
}

func Test_NewConfig_WithoutStruct(t *testing.T) {
	lst := []string{"sadsa"}
	_, err := config.NewConfig(nil, &lst)
	if err == nil {
		t.Fatalf("expected error but got none")
	}
}

func Test_NewConfig_NilData(t *testing.T) {
	var defaultConfig *config.Config[NewConfig_TestConfig] = nil
	_, err := config.NewConfig(nil, defaultConfig)
	if err == nil {
		t.Fatal("expected error for uninitialized data")
	}
	if !contains(err.Error(), "invalid config struct: config struct cannot be nil") {
		t.Errorf("unexpected error message: %v", err)
	}
}
