package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Scope   string     `yaml:"scope"`
	Org     string     `yaml:"org,omitempty"`
	Repo    RepoConfig `yaml:"repo,omitempty"`
	Token   string     `yaml:"token"`
	Runners RunnerConf `yaml:"runners"`
	Docker  DockerConf `yaml:"docker"`
}

type RepoConfig struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
}

type RunnerConf struct {
	Count      int               `yaml:"count"`
	Image      string            `yaml:"image"`
	Labels     []string          `yaml:"labels"`
	Group      string            `yaml:"group"`
	NamePrefix string            `yaml:"name_prefix"`
	Ephemeral  bool              `yaml:"ephemeral"`
	ExtraEnv   map[string]string `yaml:"extra_env,omitempty"`
}

type DockerConf struct {
	Socket            string `yaml:"socket"`
	MountDockerSocket bool   `yaml:"mount_docker_socket"`
	RestartPolicy     string `yaml:"restart_policy"`
	WorkDirBase       string `yaml:"work_dir_base,omitempty"`
}

// Dir returns the ghr config directory (~/.ghr).
func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".ghr"
	}
	return filepath.Join(home, ".ghr")
}

// DefaultConfigPath returns the default config file path (~/.ghr/config.yaml).
func DefaultConfigPath() string {
	return filepath.Join(Dir(), "config.yaml")
}

// DotenvPath returns the path to the .env file inside the config dir (~/.ghr/.env).
func DotenvPath() string {
	return filepath.Join(Dir(), ".env")
}

func Default() *Config {
	return &Config{
		Scope: "org",
		Token: "env:GH_TOKEN",
		Runners: RunnerConf{
			Count:      10,
			Image:      "myoung34/github-runner:latest",
			Labels:     []string{"local", "dev"},
			Group:      "Default",
			NamePrefix: "ghr",
			Ephemeral:  true,
		},
		Docker: DockerConf{
			Socket:            "/var/run/docker.sock",
			MountDockerSocket: true,
			RestartPolicy:     "unless-stopped",
		},
	}
}

// searchPaths returns config file paths in priority order.
func searchPaths() []string {
	return []string{
		DefaultConfigPath(),
	}
}

// Load reads config from the given path, or searches default locations.
func Load(path string) (*Config, string, error) {
	if path != "" {
		cfg, err := loadFrom(path)
		return cfg, path, err
	}
	for _, p := range searchPaths() {
		if _, err := os.Stat(p); err == nil {
			cfg, err := loadFrom(p)
			return cfg, p, err
		}
	}
	return nil, "", fmt.Errorf("no config file found; run `ghr init` to create one")
}

func loadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return cfg, nil
}

// Save writes the config to the given path.
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Validate checks required fields.
func Validate(cfg *Config) error {
	if cfg.Scope != "org" && cfg.Scope != "repo" {
		return fmt.Errorf("scope must be 'org' or 'repo', got %q", cfg.Scope)
	}
	if cfg.Scope == "org" && cfg.Org == "" {
		return fmt.Errorf("org name is required when scope is 'org'")
	}
	if cfg.Scope == "repo" {
		if cfg.Repo.Owner == "" || cfg.Repo.Name == "" {
			return fmt.Errorf("repo owner and name are required when scope is 'repo'")
		}
	}
	if cfg.Token == "" {
		return fmt.Errorf("token is required")
	}
	if cfg.Runners.Image == "" {
		return fmt.Errorf("runner image is required")
	}
	if cfg.Runners.NamePrefix == "" {
		return fmt.Errorf("runner name_prefix is required")
	}
	return nil
}
