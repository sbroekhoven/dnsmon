package config

import (
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	conf, err := LoadConfiguration("test_config.yaml")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if conf.DomainsFile == "" {
		t.Errorf("Expected DomainsFile to be set, got empty string")
	}

	// Add more assertions as needed
}
