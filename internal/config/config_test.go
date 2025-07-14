package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("./config-test.yml")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Services) != 1 {
		t.Errorf("expected 1 service, got %d", len(cfg.Services))
	}
	if cfg.Services[0].Name != "test-api" {
		t.Errorf("expected service name 'test-api', got '%s'", cfg.Services[0].Name)
	}
}
