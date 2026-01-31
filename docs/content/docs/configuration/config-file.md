---
title: Config File Reference
weight: 1
---

ghr stores its configuration at `~/.ghr/config.yaml`. This page documents every field.

## Full Example

```yaml
scope: org
org: my-org
repo:
  owner: ""
  name: ""
token: env:GH_TOKEN
runners:
  count: 10
  image: myoung34/github-runner:latest
  labels:
    - local
    - dev
  group: Default
  name_prefix: ghr
  ephemeral: true
  extra_env: {}
docker:
  socket: /var/run/docker.sock
  mount_docker_socket: true
  restart_policy: unless-stopped
  work_dir_base: ""
```

## Top-level Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `scope` | `string` | `"org"` | Runner scope. Must be `"org"` or `"repo"`. |
| `org` | `string` | `""` | GitHub organization name. **Required** when `scope` is `"org"`. |
| `repo` | `object` | -- | Repository details. **Required** when `scope` is `"repo"`. |
| `repo.owner` | `string` | `""` | Repository owner (user or org). |
| `repo.name` | `string` | `""` | Repository name. |
| `token` | `string` | `"env:GH_TOKEN"` | GitHub token. Supports `env:VAR` syntax. See [Token Setup](../token-setup). |
| `runners` | `object` | -- | Runner configuration. |
| `docker` | `object` | -- | Docker configuration. |

## Runner Configuration (`runners`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `count` | `int` | `10` | Default number of runners for `ghr up` (when no argument is given). |
| `image` | `string` | `"myoung34/github-runner:latest"` | Docker image for runner containers. **Required**. |
| `labels` | `[]string` | `["local", "dev"]` | Labels attached to each runner. These appear in the GitHub UI and can be used in `runs-on`. |
| `group` | `string` | `"Default"` | GitHub runner group name. |
| `name_prefix` | `string` | `"ghr"` | Prefix for container and runner names. Containers are named `{prefix}-runner-{N}`. **Required**. |
| `ephemeral` | `bool` | `true` | If true, runners de-register after completing one job. See [Ephemeral Runners](../../guides/ephemeral-runners). |
| `extra_env` | `map[string]string` | `{}` | Additional environment variables passed to the runner container. |

## Docker Configuration (`docker`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `socket` | `string` | `"/var/run/docker.sock"` | Path to the Docker socket. Usually auto-detected from the Docker context. |
| `mount_docker_socket` | `bool` | `true` | Mount the Docker socket inside the runner container. Required for workflows that use Docker actions. |
| `restart_policy` | `string` | `"unless-stopped"` | Docker restart policy for runner containers. Common values: `"no"`, `"always"`, `"unless-stopped"`, `"on-failure"`. |
| `work_dir_base` | `string` | `""` | Base directory for runner work directories. If empty, Docker named volumes are used instead of bind mounts. |

## Config Directory

All ghr files live in `~/.ghr/`:

```
~/.ghr/
  config.yaml   # main configuration
  .env          # environment variables (GH_TOKEN, etc.)
```

## Overriding the Config Path

Use the `--config` flag on any command to use a different config file:

```bash
ghr --config /path/to/config.yaml list
```

## Validation Rules

ghr validates the config on every command that needs it. The following rules are enforced:

- `scope` must be `"org"` or `"repo"`
- `org` is required when `scope` is `"org"`
- `repo.owner` and `repo.name` are required when `scope` is `"repo"`
- `token` must not be empty
- `runners.image` must not be empty
- `runners.name_prefix` must not be empty

Commands that do not require a config (such as `init`, `completion`, and `version`) skip validation.
