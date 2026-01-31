package docker

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// Client wraps the Docker Engine SDK.
type Client struct {
	cli *client.Client
}

// NewClient creates a Docker client using the given socket path.
// If host is the default "/var/run/docker.sock", it also checks the active
// Docker CLI context (e.g. OrbStack, colima) for the real endpoint.
func NewClient(host string) (*Client, error) {
	opts := []client.Opt{client.FromEnv, client.WithAPIVersionNegotiation()}

	// If DOCKER_HOST is already set, FromEnv handles it.
	// Otherwise, try to resolve from the active docker context.
	if os.Getenv("DOCKER_HOST") == "" {
		if endpoint := resolveDockerContext(); endpoint != "" {
			opts = append(opts, client.WithHost(endpoint))
		} else if host != "" {
			opts = append(opts, client.WithHost("unix://"+host))
		}
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("creating docker client: %w", err)
	}
	return &Client{cli: cli}, nil
}

// resolveDockerContext reads ~/.docker/config.json to find the current context,
// then reads the context metadata to extract the Docker endpoint.
func resolveDockerContext() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configPath := filepath.Join(home, ".docker", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	var dockerConfig struct {
		CurrentContext string `json:"currentContext"`
	}
	if err := json.Unmarshal(data, &dockerConfig); err != nil || dockerConfig.CurrentContext == "" {
		return ""
	}
	if dockerConfig.CurrentContext == "default" {
		return ""
	}

	// Docker CLI hashes the context name with SHA256 for the directory name
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(dockerConfig.CurrentContext)))
	metaPath := filepath.Join(home, ".docker", "contexts", "meta", hash, "meta.json")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		return ""
	}

	var meta struct {
		Endpoints map[string]struct {
			Host string `json:"Host"`
		} `json:"Endpoints"`
	}
	if err := json.Unmarshal(metaData, &meta); err != nil {
		return ""
	}
	if ep, ok := meta.Endpoints["docker"]; ok && ep.Host != "" {
		return ep.Host
	}
	return ""
}

// Close closes the underlying Docker client.
func (c *Client) Close() error {
	return c.cli.Close()
}

// PullImage pulls the specified image.
func (c *Client) PullImage(ctx context.Context, ref string) error {
	reader, err := c.cli.ImagePull(ctx, ref, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("pulling image %s: %w", ref, err)
	}
	defer reader.Close()
	// Consume output to complete the pull.
	_, _ = io.Copy(io.Discard, reader)
	return nil
}

// ContainerLogs returns a reader for the container's logs.
func (c *Client) ContainerLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error) {
	return c.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
		Tail:       "100",
	})
}
