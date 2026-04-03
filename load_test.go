package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/makadev/go-config"
	"go.yaml.in/yaml/v3"
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

func Test_LoadConfig_UsesFirstExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	secondPath := filepath.Join(tempDir, "second.json")
	thirdPath := filepath.Join(tempDir, "third.json")

	if err := os.WriteFile(secondPath, []byte(`{"app_name":"from-second","server":{"host":"127.0.0.2","port":9001}}`), 0o644); err != nil {
		t.Fatalf("failed to write second config: %v", err)
	}
	if err := os.WriteFile(thirdPath, []byte(`{"app_name":"from-third","server":{"host":"127.0.0.3","port":9002}}`), 0o644); err != nil {
		t.Fatalf("failed to write third config: %v", err)
	}

	opts := config.NewOptions()
	opts.SkipEnv = true
	opts.ConfigPaths = []string{
		filepath.Join(tempDir, "missing.json"),
		secondPath,
		thirdPath,
	}

	data := &Load_ConfigData{AppName: "default"}
	cfg, err := config.NewConfig(opts, data)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if err := cfg.Load(); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if data.AppName != "from-second" {
		t.Fatalf("expected first existing file to win, got %q", data.AppName)
	}
	if data.Server.Host != "127.0.0.2" {
		t.Fatalf("expected host from second file, got %q", data.Server.Host)
	}
	if data.Server.Port != 9001 {
		t.Fatalf("expected port from second file, got %d", data.Server.Port)
	}
	if data.AppName == "from-third" {
		t.Fatal("expected loader to stop after first successful file")
	}
}

func Test_LoadConfig_MissingFilesAreIgnored(t *testing.T) {
	opts := config.NewOptions()
	opts.SkipEnv = true
	opts.ConfigPaths = []string{
		filepath.Join(t.TempDir(), "missing.yaml"),
		filepath.Join(t.TempDir(), "missing.json"),
	}

	data := &Load_ConfigData{
		AppName: "default-name",
		Server:  Load_ServerConfig{Host: "127.0.0.1", Port: 5000},
	}
	cfg, err := config.NewConfig(opts, data)
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	if err := cfg.Load(); err != nil {
		t.Fatalf("expected missing files to be ignored, got %v", err)
	}

	if data.AppName != "default-name" {
		t.Fatalf("expected config to keep default values, got %q", data.AppName)
	}
	if data.Server.Port != 5000 {
		t.Fatalf("expected port to remain unchanged, got %d", data.Server.Port)
	}
}

func Test_LoadConfig_UnsupportedExtension(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.txt")
	if err := os.WriteFile(configPath, []byte("app_name: ignored"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	opts := config.NewOptions()
	opts.SkipEnv = true
	opts.ConfigPaths = []string{configPath}

	cfg, err := config.NewConfig(opts, &Load_ConfigData{})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	err = cfg.Load()
	if err == nil {
		t.Fatal("expected error for unsupported file extension")
	}
	if !strings.Contains(err.Error(), "unsupported file extension: .txt") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_LoadConfig_StatNotDirError(t *testing.T) {
	tempDir := t.TempDir()
	// Create a regular file — stat-ing a path *inside* it returns ENOTDIR,
	// which is not os.IsNotExist, exercising the stat-error branch.
	regularFile := filepath.Join(tempDir, "notadir")
	if err := os.WriteFile(regularFile, []byte("not a dir"), 0o644); err != nil {
		t.Fatalf("failed to create regular file: %v", err)
	}

	opts := config.NewOptions()
	opts.SkipEnv = true
	// The path treats 'notadir' (a file) as if it were a directory.
	opts.ConfigPaths = []string{filepath.Join(regularFile, "config.yaml")}

	cfg, err := config.NewConfig(opts, &Load_ConfigData{})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	err = cfg.Load()
	if err == nil {
		t.Fatal("expected a stat error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to stat config file") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_LoadConfig_InvalidFileContents(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		contents    string
		wantErrPart string
	}{
		{
			name:        "invalid yaml",
			fileName:    "config.yaml",
			contents:    "app_name: [unterminated",
			wantErrPart: "failed to unmarshal YAML",
		},
		{
			name:        "invalid json",
			fileName:    "config.json",
			contents:    `{"app_name":`,
			wantErrPart: "failed to unmarshal JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), tt.fileName)
			if err := os.WriteFile(configPath, []byte(tt.contents), 0o644); err != nil {
				t.Fatalf("failed to write config file: %v", err)
			}

			opts := config.NewOptions()
			opts.SkipEnv = true
			opts.ConfigPaths = []string{configPath}

			cfg, err := config.NewConfig(opts, &Load_ConfigData{})
			if err != nil {
				t.Fatalf("failed to initialize config: %v", err)
			}

			err = cfg.Load()
			if err == nil {
				t.Fatal("expected invalid config content to fail")
			}
			if !strings.Contains(err.Error(), tt.wantErrPart) {
				t.Fatalf("expected error to contain %q, got %v", tt.wantErrPart, err)
			}
		})
	}
}

func Test_LoadDump_RoundTripFixtures(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		format    string
		unmarshal func([]byte, interface{}) error
	}{
		{
			name:      "json",
			path:      "testdata/load.json",
			format:    "json",
			unmarshal: json.Unmarshal,
		},
		{
			name:      "yaml",
			path:      "testdata/load.yaml",
			format:    "yaml",
			unmarshal: yaml.Unmarshal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalBytes, err := os.ReadFile(tt.path)
			if err != nil {
				t.Fatalf("failed to read fixture %s: %v", tt.path, err)
			}

			opts := config.NewOptions()
			opts.SkipEnv = true
			opts.ConfigPaths = []string{tt.path}

			data := &Load_ConfigData{}
			cfg, err := config.NewConfig(opts, data)
			if err != nil {
				t.Fatalf("failed to initialize config: %v", err)
			}

			if err := cfg.Load(); err != nil {
				t.Fatalf("failed to load config: %v", err)
			}

			dumped, err := cfg.DumpWithOptions(&config.DumpOptions{
				Format:      tt.format,
				Content:     "config",
				MaskSecrets: false,
				MaskWith:    "***",
			})
			if err != nil {
				t.Fatalf("failed to dump config in %s format: %v", tt.format, err)
			}

			var originalData map[string]interface{}
			if err := tt.unmarshal(originalBytes, &originalData); err != nil {
				t.Fatalf("failed to parse original %s fixture: %v", tt.format, err)
			}

			var roundTrippedData map[string]interface{}
			if err := tt.unmarshal([]byte(dumped), &roundTrippedData); err != nil {
				t.Fatalf("failed to parse dumped %s output: %v", tt.format, err)
			}

			if !reflect.DeepEqual(originalData, roundTrippedData) {
				t.Fatalf("roundtrip mismatch\noriginal: %#v\ndumped:   %#v", originalData, roundTrippedData)
			}
		})
	}
}
