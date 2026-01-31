package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Scope != "org" {
		t.Errorf("default scope = %q, want %q", cfg.Scope, "org")
	}
	if cfg.Runners.NamePrefix != "ghr" {
		t.Errorf("default name_prefix = %q, want %q", cfg.Runners.NamePrefix, "ghr")
	}
	if cfg.Runners.Image != "myoung34/github-runner:latest" {
		t.Errorf("default image = %q", cfg.Runners.Image)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr bool
	}{
		{"valid org config", func(c *Config) { c.Org = "myorg" }, false},
		{"missing org", func(c *Config) {}, true},
		{"invalid scope", func(c *Config) { c.Scope = "invalid" }, true},
		{"valid repo config", func(c *Config) {
			c.Scope = "repo"
			c.Repo = RepoConfig{Owner: "owner", Name: "repo"}
		}, false},
		{"repo missing name", func(c *Config) {
			c.Scope = "repo"
			c.Repo = RepoConfig{Owner: "owner"}
		}, true},
		{"missing token", func(c *Config) { c.Org = "myorg"; c.Token = "" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			tt.modify(cfg)
			err := Validate(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")

	cfg := Default()
	cfg.Org = "testorg"

	if err := Save(cfg, path); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, loadedPath, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loadedPath != path {
		t.Errorf("Load() path = %q, want %q", loadedPath, path)
	}
	if loaded.Org != "testorg" {
		t.Errorf("loaded org = %q, want %q", loaded.Org, "testorg")
	}
	if loaded.Runners.NamePrefix != "ghr" {
		t.Errorf("loaded name_prefix = %q, want %q", loaded.Runners.NamePrefix, "ghr")
	}
}

func TestResolveToken(t *testing.T) {
	t.Run("env: prefix", func(t *testing.T) {
		os.Setenv("TEST_GHR_TOKEN", "secret123")
		defer os.Unsetenv("TEST_GHR_TOKEN")

		got, err := ResolveToken("env:TEST_GHR_TOKEN")
		if err != nil {
			t.Fatalf("ResolveToken() error = %v", err)
		}
		if got != "secret123" {
			t.Errorf("ResolveToken() = %q, want %q", got, "secret123")
		}
	})

	t.Run("env: prefix missing var", func(t *testing.T) {
		_, err := ResolveToken("env:NONEXISTENT_VAR_12345")
		if err == nil {
			t.Error("ResolveToken() expected error for missing env var")
		}
	})

	t.Run("literal token", func(t *testing.T) {
		got, err := ResolveToken("ghp_abc123")
		if err != nil {
			t.Fatalf("ResolveToken() error = %v", err)
		}
		if got != "ghp_abc123" {
			t.Errorf("ResolveToken() = %q, want %q", got, "ghp_abc123")
		}
	})
}
