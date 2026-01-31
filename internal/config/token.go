package config

import (
	"fmt"
	"os"
	"strings"
)

// LoadDotenv loads variables from a .env file into the process environment.
// Variables already set in the environment are not overwritten.
func LoadDotenv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
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
		val = strings.TrimSpace(val)
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
