package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseDotenv reads a .env file and returns the key-value pairs.
func ParseDotenv(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vars := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := parts[1]
		// Strip inline comments (e.g. "value  # comment")
		if idx := strings.Index(val, " #"); idx >= 0 {
			val = val[:idx]
		}
		vars[key] = strings.TrimSpace(val)
	}
	return vars, nil
}

// SaveDotenv writes vars into a .env file, merging with any existing content.
// Existing keys not in vars are preserved; keys in vars overwrite existing ones.
// The file is created with 0600 permissions if it does not exist.
func SaveDotenv(path string, vars map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// Read existing file to preserve comments and non-overwritten keys.
	existing, _ := os.ReadFile(path)

	wrote := make(map[string]bool)
	var out []string

	if len(existing) > 0 {
		for _, line := range strings.Split(string(existing), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				out = append(out, line)
				continue
			}
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				if newVal, ok := vars[key]; ok {
					out = append(out, key+"="+newVal)
					wrote[key] = true
					continue
				}
			}
			out = append(out, line)
		}
	}

	// Append new keys that weren't already in the file.
	for k, v := range vars {
		if !wrote[k] {
			out = append(out, k+"="+v)
		}
	}

	return os.WriteFile(path, []byte(strings.Join(out, "\n")+"\n"), 0600)
}

// LoadDotenv loads variables from a .env file into the process environment.
// Variables already set in the environment are not overwritten.
func LoadDotenv(path string) {
	vars, err := ParseDotenv(path)
	if err != nil {
		return
	}
	for key, val := range vars {
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, val)
		}
	}
}

// ResolveToken resolves the token value. If the token starts with "env:",
// it reads the value from the named environment variable.
func ResolveToken(raw string) (string, error) {
	if strings.HasPrefix(raw, "env:") {
		envVar := strings.TrimPrefix(raw, "env:")
		val := os.Getenv(envVar)
		if val == "" {
			return "", fmt.Errorf("environment variable %q is not set", envVar)
		}
		return val, nil
	}
	return raw, nil
}
