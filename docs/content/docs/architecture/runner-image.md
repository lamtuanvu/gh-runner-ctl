---
title: Runner Image
weight: 4
---

ghr creates runner containers from a configurable Docker image. The default image is [`myoung34/github-runner`](https://github.com/myoung34/docker-github-actions-runner), a widely used community image for self-hosted runners.

## Environment Variable Mapping

Each runner container receives the following environment variables, mapped from ghr's configuration:

| Environment Variable | Source | Description |
|---------------------|--------|-------------|
| `RUNNER_SCOPE` | `config.scope` | `org` or `repo` |
| `ORG_NAME` | `config.org` | Organization name (when scope is `org`) |
| `REPO_URL` | Derived from `config.repo` | Full repo URL `https://github.com/{owner}/{name}` (when scope is `repo`) |
| `RUNNER_NAME` | Derived | Container name (e.g., `ghr-runner-3`) |
| `RUNNER_LABELS` | `config.runners.labels` | Comma-separated labels |
| `RUNNER_GROUP` | `config.runners.group` | Runner group name |
| `ACCESS_TOKEN` | Resolved `config.token` | GitHub PAT (resolved from `env:VAR` if applicable) |
| `EPHEMERAL` | `config.runners.ephemeral` | Set to `true` when ephemeral mode is enabled |

Any additional variables from `runners.extra_env` are also passed to the container.

## Container Mounts

ghr configures two types of mounts on each runner container:

### Docker Socket

When `docker.mount_docker_socket` is `true` (the default), the Docker socket is bind-mounted into the container. This is required for workflows that use Docker-based actions (e.g., `docker/build-push-action`).

### Work Directory

The work directory mount depends on the `docker.work_dir_base` config:

| `work_dir_base` | Behavior |
|-----------------|----------|
| Empty (default) | A Docker **named volume** is created for each runner |
| A path | A **bind mount** at `{work_dir_base}/runner-{N}` is used |

Named volumes are managed by Docker and are the recommended option for most setups.

## Custom Images

You can use any Docker image that follows the same environment variable conventions as `myoung34/github-runner`. Set the image in your config:

```yaml
runners:
  image: ghcr.io/my-org/custom-runner:latest
```

The image must:

1. Accept the environment variables listed above
2. Register itself as a GitHub Actions runner on startup
3. Handle graceful shutdown on `SIGTERM`

## Restart Policy

The `docker.restart_policy` config controls what happens when a runner container exits:

| Policy | Behavior |
|--------|----------|
| `"no"` | Do not restart |
| `"always"` | Always restart, regardless of exit status |
| `"unless-stopped"` | Restart unless explicitly stopped (default) |
| `"on-failure"` | Restart only on non-zero exit status |

For ephemeral runners, consider using `"no"` since the runner is intended to exit after one job.
