package config

import (
	"os"
	"path/filepath"
	"strings"
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

func TestParseDotenv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	content := `# This is a comment
GH_TOKEN=abc123
GH_ORG=myorg

RUNNER_LABELS=linux,x64 # inline comment
EMPTY_LINE_ABOVE=yes
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	vars, err := ParseDotenv(path)
	if err != nil {
		t.Fatalf("ParseDotenv() error = %v", err)
	}

	tests := map[string]string{
		"GH_TOKEN":         "abc123",
		"GH_ORG":           "myorg",
		"RUNNER_LABELS":    "linux,x64",
		"EMPTY_LINE_ABOVE": "yes",
	}
	for k, want := range tests {
		if got := vars[k]; got != want {
			t.Errorf("ParseDotenv()[%q] = %q, want %q", k, got, want)
		}
	}
	if len(vars) != len(tests) {
		t.Errorf("ParseDotenv() returned %d vars, want %d", len(vars), len(tests))
	}
}

func TestParseDotenv_FileNotFound(t *testing.T) {
	_, err := ParseDotenv("/nonexistent/.env")
	if err == nil {
		t.Error("ParseDotenv() expected error for missing file")
	}
}

func TestSaveDotenv_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", ".env")

	vars := map[string]string{
		"GH_TOKEN": "secret",
		"GH_ORG":   "myorg",
	}
	if err := SaveDotenv(path, vars); err != nil {
		t.Fatalf("SaveDotenv() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o, want 0600", perm)
	}

	got, err := ParseDotenv(path)
	if err != nil {
		t.Fatalf("ParseDotenv() error = %v", err)
	}
	for k, want := range vars {
		if got[k] != want {
			t.Errorf("key %q = %q, want %q", k, got[k], want)
		}
	}
}

func TestSaveDotenv_MergeExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	existing := `# Header comment
OLD_KEY=old_value
SHARED_KEY=original
`
	if err := os.WriteFile(path, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"SHARED_KEY": "updated",
		"NEW_KEY":    "new_value",
	}
	if err := SaveDotenv(path, vars); err != nil {
		t.Fatalf("SaveDotenv() error = %v", err)
	}

	got, err := ParseDotenv(path)
	if err != nil {
		t.Fatalf("ParseDotenv() error = %v", err)
	}
	if got["OLD_KEY"] != "old_value" {
		t.Errorf("OLD_KEY = %q, want %q", got["OLD_KEY"], "old_value")
	}
	if got["SHARED_KEY"] != "updated" {
		t.Errorf("SHARED_KEY = %q, want %q", got["SHARED_KEY"], "updated")
	}
	if got["NEW_KEY"] != "new_value" {
		t.Errorf("NEW_KEY = %q, want %q", got["NEW_KEY"], "new_value")
	}

	// Verify comment is preserved in raw content
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "# Header comment") {
		t.Error("comment was not preserved in merged file")
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
