package config_test

import (
	"testing"

	"github.com/makadev/go-config"
)

type Load_ServerConfig struct {
	Host string `yaml:"host" json:"host" config:"host" env:"HOST"`
	Port int    `yaml:"port" json:"port" config:"port" env:"PORT"`
}

type Load_ConfigData struct {
	AppName string            `yaml:"app_name" json:"app_name" config:"app_name" env:"APP_NAME"`
	AppUrl  string            `yaml:"app_url" json:"app_url" config:"app_url" env:"APP_URL"`
	Server  Load_ServerConfig `yaml:"server" json:"server" config:"server" env:"SERVER"`
}

func Test_LoadConfig_YAML(t *testing.T) {
	opts := config.NewOptions()
	opts.ConfigPaths = []string{"testdata/load.yaml"}
	cfg, err := config.NewConfig(opts, &Load_ConfigData{
		AppName: "TestApp",
		AppUrl:  "http://127.0.0.1",
		Server: Load_ServerConfig{
			Host: "127.0.0.1",
			Port: 5000,
		},
	})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if err := cfg.Load(); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
}

func Test_LoadConfig_JSON(t *testing.T) {
	opts := config.NewOptions()
	opts.ConfigPaths = []string{"testdata/load.json"}
	cfg, err := config.NewConfig(opts, &Load_ConfigData{
		AppName: "TestApp",
		AppUrl:  "http://127.0.0.1",
		Server: Load_ServerConfig{
			Host: "127.0.0.1",
			Port: 5000,
		},
	})

	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if err := cfg.Load(); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
}
