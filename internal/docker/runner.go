package docker

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/lamtuanvu/gh-runner-ctl/internal/config"
)

// RunnerContainer holds info about a managed runner container.
type RunnerContainer struct {
	ID     string
	Name   string
	Num    int
	State  string // running, exited, created, etc.
	Status string // human-readable status from Docker
	Labels map[string]string
}

// CreateRunner creates and starts a new runner container.
func (c *Client) CreateRunner(ctx context.Context, cfg *config.Config, num int, token string) (string, error) {
	name := fmt.Sprintf("%s-runner-%d", cfg.Runners.NamePrefix, num)

	env := []string{
		"RUNNER_SCOPE=" + cfg.Scope,
		"RUNNER_NAME=" + name,
		"RUNNER_LABELS=" + strings.Join(cfg.Runners.Labels, ","),
		"RUNNER_GROUP=" + cfg.Runners.Group,
		"ACCESS_TOKEN=" + token,
	}
	if cfg.Scope == "org" {
		env = append(env, "ORG_NAME="+cfg.Org)
	} else {
		env = append(env, "REPO_URL=https://github.com/"+cfg.Repo.Owner+"/"+cfg.Repo.Name)
	}
	if cfg.Runners.Ephemeral {
		env = append(env, "EPHEMERAL=true")
	}
	for k, v := range cfg.Runners.ExtraEnv {
		env = append(env, k+"="+v)
	}

	labels := ManagedLabels(cfg.Scope, cfg.Org, cfg.Repo.Owner, cfg.Repo.Name, num)

	var mounts []mount.Mount
	if cfg.Docker.MountDockerSocket {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: cfg.Docker.Socket,
			Target: "/var/run/docker.sock",
		})
	}

	// Work directory: bind mount if work_dir_base is set, otherwise named volume.
	if cfg.Docker.WorkDirBase != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: cfg.Docker.WorkDirBase + "/" + name,
			Target: "/home/runner/_work",
		})
	} else {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: name + "-work",
			Target: "/home/runner/_work",
		})
	}

	restartPolicy := container.RestartPolicy{Name: container.RestartPolicyMode(cfg.Docker.RestartPolicy)}

	resp, err := c.cli.ContainerCreate(ctx,
		&container.Config{
			Image:  cfg.Runners.Image,
			Env:    env,
			Labels: labels,
		},
		&container.HostConfig{
			Mounts:        mounts,
			RestartPolicy: restartPolicy,
		},
		nil, nil, name,
	)
	if err != nil {
		return "", fmt.Errorf("creating container %s: %w", name, err)
	}

	if err := c.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("starting container %s: %w", name, err)
	}

	return resp.ID, nil
}

// ListManagedContainers returns all ghr-managed containers (including stopped).
func (c *Client) ListManagedContainers(ctx context.Context) ([]RunnerContainer, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: ManagedFilter(),
	})
	if err != nil {
		return nil, fmt.Errorf("listing containers: %w", err)
	}

	var runners []RunnerContainer
	for _, ctr := range containers {
		num, _ := strconv.Atoi(ctr.Labels[LabelRunnerNum])
		name := ""
		if len(ctr.Names) > 0 {
			name = strings.TrimPrefix(ctr.Names[0], "/")
		}
		runners = append(runners, RunnerContainer{
			ID:     ctr.ID[:12],
			Name:   name,
			Num:    num,
			State:  ctr.State,
			Status: ctr.Status,
			Labels: ctr.Labels,
		})
	}
	return runners, nil
}

// StopRunner stops a container by ID or name.
func (c *Client) StopRunner(ctx context.Context, idOrName string) error {
	timeout := 30
	return c.cli.ContainerStop(ctx, idOrName, container.StopOptions{Timeout: &timeout})
}

// StartRunner starts a stopped container by ID or name.
func (c *Client) StartRunner(ctx context.Context, idOrName string) error {
	return c.cli.ContainerStart(ctx, idOrName, container.StartOptions{})
}

// RemoveRunner removes a container by ID or name, forcing removal if running.
func (c *Client) RemoveRunner(ctx context.Context, idOrName string) error {
	return c.cli.ContainerRemove(ctx, idOrName, container.RemoveOptions{Force: true})
}

// InspectRunner returns full container info.
func (c *Client) InspectRunner(ctx context.Context, idOrName string) (types.ContainerJSON, error) {
	return c.cli.ContainerInspect(ctx, idOrName)
}
