package config_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/makadev/go-config"
)

type ThreadsafeConfig struct {
	Value  int `json:"value" config:"value" env:"VALUE"`
	Nested *struct {
		Name string `json:"name" config:"name" env:"NAME"`
	} `json:"nested" config:"nested"`
}

func Test_ConfigMethodLevelThreadSafety(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "threadsafe.json")
	fileData := ThreadsafeConfig{
		Value: 99,
		Nested: &struct {
			Name string `json:"name" config:"name" env:"NAME"`
		}{
			Name: "from-file",
		},
	}

	payload, err := json.Marshal(fileData)
	if err != nil {
		t.Fatalf("failed to marshal config payload: %v", err)
	}
	if err := os.WriteFile(configPath, payload, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	opts := config.NewOptions()
	opts.SkipEnv = true
	opts.ConfigPaths = []string{configPath}

	cfg, err := config.NewConfig(opts, &ThreadsafeConfig{})
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}

	const workers = 24
	const iterations = 100

	errCh := make(chan error, workers*iterations)
	var wg sync.WaitGroup

	for worker := 0; worker < workers; worker++ {
		workerID := worker
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				switch i % 4 {
				case 0:
					if err := cfg.Load(); err != nil {
						errCh <- fmt.Errorf("load failed: %w", err)
					}
				case 1:
					if err := cfg.SetConfigValue("value", workerID*iterations+i); err != nil {
						errCh <- fmt.Errorf("set failed: %w", err)
					}
				case 2:
					if _, err := cfg.GetFieldValue("Nested.Name"); err != nil {
						errCh <- fmt.Errorf("get failed: %w", err)
					}
				case 3:
					if _, err := cfg.DumpWithOptions(&config.DumpOptions{
						Format:      "json",
						Content:     "config",
						MaskSecrets: true,
						MaskWith:    "***",
					}); err != nil {
						errCh <- fmt.Errorf("dump failed: %w", err)
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)
	for callErr := range errCh {
		t.Fatal(callErr)
	}
}
