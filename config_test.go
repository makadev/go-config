package config_test

import (
	"sync"
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

func newTestConfig(t *testing.T) *config.Config[NewConfig_TestConfig] {
	t.Helper()
	cfg, err := config.NewConfig(nil, &NewConfig_TestConfig{AppName: "test"})
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	return cfg
}

func Test_Config_LockUnlock(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.Lock()
	cfg.Data.AppName = "locked-write"
	cfg.Unlock()
	if cfg.Data.AppName != "locked-write" {
		t.Errorf("expected AppName to be 'locked-write', got %q", cfg.Data.AppName)
	}
}

func Test_Config_RLockRUnlock(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.Data.AppName = "read-value"
	cfg.RLock()
	name := cfg.Data.AppName
	cfg.RUnlock()
	if name != "read-value" {
		t.Errorf("expected AppName to be 'read-value', got %q", name)
	}
}

func Test_Config_WithLock(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.WithLock(func() {
		cfg.Data.AppName = "with-lock"
	})
	if cfg.Data.AppName != "with-lock" {
		t.Errorf("expected AppName to be 'with-lock', got %q", cfg.Data.AppName)
	}
}

func Test_Config_WithRLock(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.Data.AppName = "with-rlock"
	var got string
	cfg.WithRLock(func() {
		got = cfg.Data.AppName
	})
	if got != "with-rlock" {
		t.Errorf("expected AppName to be 'with-rlock', got %q", got)
	}
}

func Test_Config_ConcurrentLocking(t *testing.T) {
	cfg := newTestConfig(t)
	cfg.Data.AppName = "initial"

	const goroutines = 20
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if id%2 == 0 {
				cfg.WithLock(func() {
					cfg.Data.AppName = "writer"
				})
			} else {
				cfg.WithRLock(func() {
					_ = cfg.Data.AppName
				})
			}
		}(i)
	}
	wg.Wait()
}
