package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary atlas.hcl file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "atlas.hcl")

	configContent := `
env "local" {
  url = "postgres://user:pass@localhost:5432/testdb"
  migration {
    dir = "file://migrations"
    revisions_schema = "custom_revisions"
  }
}

env "prod" {
  url = "postgres://user:pass@prod:5432/proddb"
  migration {
    dir = "file://migrations"
  }
}
`

	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test loading config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test local environment
	localEnv, err := config.GetEnv("local")
	if err != nil {
		t.Errorf("Failed to get local env: %v", err)
	}
	if localEnv.URL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("Expected local URL to be 'postgres://user:pass@localhost:5432/testdb', got '%s'", localEnv.URL)
	}
	if localEnv.RevisionsSchema != "custom_revisions" {
		t.Errorf("Expected revisions_schema to be 'custom_revisions', got '%s'", localEnv.RevisionsSchema)
	}

	// Test prod environment
	prodEnv, err := config.GetEnv("prod")
	if err != nil {
		t.Errorf("Failed to get prod env: %v", err)
	}
	if prodEnv.URL != "postgres://user:pass@prod:5432/proddb" {
		t.Errorf("Expected prod URL to be 'postgres://user:pass@prod:5432/proddb', got '%s'", prodEnv.URL)
	}

	// Test non-existent environment
	_, err = config.GetEnv("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent environment")
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/atlas.hcl")
	if err == nil {
		t.Error("Expected error for non-existent config file")
	}
}

func TestGetEnvList(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "atlas.hcl")

	configContent := `
env "dev" {
  url = "postgres://localhost/dev"
}

env "staging" {
  url = "postgres://localhost/staging"
}

env "prod" {
  url = "postgres://localhost/prod"
}
`

	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Envs) != 3 {
		t.Errorf("Expected 3 environments, got %d", len(config.Envs))
	}

	expectedEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}
	for name := range config.Envs {
		if !expectedEnvs[name] {
			t.Errorf("Unexpected environment: %s", name)
		}
	}
}
