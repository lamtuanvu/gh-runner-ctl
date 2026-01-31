package config

import (
	"fmt"
	"os"
	"strings"
)

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
