# ghr - GitHub Self-Hosted Runner Manager

A CLI tool for managing GitHub Actions self-hosted runners as Docker containers. Replaces verbose `docker-compose.yml` files with a single command-line interface.

## Motivation

Managing self-hosted runners with docker-compose means copy-pasting service definitions for every runner. A 10-runner setup produces ~150 lines of nearly identical YAML. Scaling, relabeling, or updating config requires editing multiple places.

`ghr` reduces this to:

```bash
ghr init        # one-time setup
ghr up 10       # launch 10 runners
ghr scale 20    # scale to 20
ghr down --all  # tear down everything
```

## Installation

### Homebrew (macOS / Linux)

```bash
brew install lamtuanvu/tap/ghr
```

### Debian / Ubuntu

Download the `.deb` package from the [latest release](https://github.com/lamtuanvu/gh-runner-ctl/releases/latest):

```bash
curl -LO https://github.com/lamtuanvu/gh-runner-ctl/releases/latest/download/ghr_<VERSION>_linux_amd64.deb
sudo dpkg -i ghr_<VERSION>_linux_amd64.deb
```

### Install script

```bash
curl -sSL https://raw.githubusercontent.com/lamtuanvu/gh-runner-ctl/main/install.sh | sh
```

To install a specific version or to a custom directory:

```bash
curl -sSL https://raw.githubusercontent.com/lamtuanvu/gh-runner-ctl/main/install.sh | sh -s -- -v v0.1.0 -d ~/.local/bin
```

### Build from source

Requires Go 1.21+.

```bash
git clone https://github.com/lamtuanvu/gh-runner-ctl.git
cd gh-runner-ctl
make build        # outputs to ./bin/ghr
make install      # installs to $GOPATH/bin
```

### Verify

```bash
ghr version
```

## Quick Start

```bash
# 1. Set your GitHub token
export GH_TOKEN=ghp_xxxxxxxxxxxx

# 2. Create config (auto-imports from .env if present)
ghr init

# 3. Launch runners
ghr up 5

# 4. Check status
ghr list
ghr status

# 5. Clean up
ghr down --all
```

## Configuration

`ghr init` creates a `.ghr.yaml` file interactively. If an `.env` file exists (from a previous docker-compose setup), it offers to import settings automatically.

### Config file

```yaml
scope: org                              # "org" or "repo"
org: econ-v1                            # required if scope=org
repo:                                   # required if scope=repo
  owner: ""
  name: ""
token: env:GH_TOKEN                     # "env:VAR" reads from env (recommended)
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
  work_dir_base: ""                     # empty = use Docker named volumes
```

### Config search order

1. `./.ghr.yaml` (current directory)
2. `~/.config/ghr/config.yaml` (user config)

Override with `--config path/to/config.yaml` on any command.

### Token configuration

The `token` field supports two formats:

| Format | Example | Description |
|--------|---------|-------------|
| `env:VAR` | `env:GH_TOKEN` | Reads from environment variable (recommended) |
| Literal | `ghp_abc123...` | Stored directly in config (avoid for shared configs) |

The token needs `admin:org` scope for org-level runners, or `repo` scope for repo-level runners.

## Commands

### `ghr init`

Interactive config setup. Detects and imports from existing `.env` files.

```bash
ghr init                # interactive prompts
ghr init --import-env   # force import from .env
```

### `ghr up [COUNT]`

Create and start new runners. Additive — does not affect existing runners.

```bash
ghr up 5    # create 5 runners
ghr up      # create runners.count from config (default 10)
```

Runner numbers fill the lowest available gaps. If runners 1, 2, 5 exist, `ghr up 2` creates runners 3 and 4.

### `ghr down [COUNT | --all]`

Stop and remove runners, highest-numbered first.

```bash
ghr down 3      # remove 3 highest-numbered runners
ghr down --all  # remove all managed runners
```

### `ghr scale COUNT`

Adjust to exactly COUNT runners. Adds or removes as needed.

```bash
ghr scale 10    # ensure exactly 10 runners exist
ghr scale 0     # same as ghr down --all
```

### `ghr list [--github]`

List managed runner containers.

```bash
ghr list            # show Docker container status
ghr list --github   # also show GitHub API online/offline/busy status
```

Example output:

```
NUM  NAME          CONTAINER     STATUS
---  ----          ---------     ------
1    ghr-runner-1  c5904092bff3  Up 2 minutes
2    ghr-runner-2  628e83bedc4b  Up 2 minutes
```

With `--github`:

```
NUM  NAME          CONTAINER     DOCKER STATUS  GITHUB   BUSY
---  ----          ---------     -------------  ------   ----
1    ghr-runner-1  c5904092bff3  Up 2 minutes   online   no
2    ghr-runner-2  628e83bedc4b  Up 2 minutes   online   yes
```

### `ghr logs NAME_OR_NUMBER [-f]`

View runner container logs. Accepts a runner number, name, or container ID.

```bash
ghr logs 1              # logs for runner #1
ghr logs ghr-runner-3   # logs by name
ghr logs 1 -f           # follow (tail) logs
```

### `ghr stop [NAME_OR_NUMBER | --all]`

Stop runners without removing them.

```bash
ghr stop 3        # stop runner #3
ghr stop --all    # stop all managed runners
```

### `ghr start [NAME_OR_NUMBER | --all]`

Start previously stopped runners.

```bash
ghr start 3       # start runner #3
ghr start --all   # start all managed runners
```

### `ghr rm [NAME_OR_NUMBER | --all]`

Remove stopped runner containers.

```bash
ghr rm 3          # remove runner #3
ghr rm --all      # remove all managed runners
```

### `ghr status`

Show config summary and runner counts.

```
Config:   .ghr.yaml
Scope:    org
Org:      econ-v1
Image:    myoung34/github-runner:latest
Labels:   local, dev
Runners:  10 total (10 running, 0 stopped)
```

### `ghr version`

Print the version string.

## How It Works

### Stateless via Docker Labels

`ghr` does not maintain a state file. It identifies containers it manages using Docker labels:

| Label | Example | Purpose |
|-------|---------|---------|
| `dev.ghr.managed` | `true` | Identifies ghr-managed containers |
| `dev.ghr.runner-num` | `3` | Runner number |
| `dev.ghr.scope` | `org` | Scope at creation time |
| `dev.ghr.org` | `econ-v1` | Org name |

This means you can safely use other Docker containers alongside `ghr` — it only touches containers with its own labels.

### Container naming

Containers are named `{prefix}-runner-{N}` (e.g., `ghr-runner-3`). This name is also registered as the GitHub runner name, making it easy to correlate across `ghr list`, `docker ps`, and the GitHub UI.

### Docker context support

`ghr` automatically detects the active Docker CLI context (OrbStack, colima, Docker Desktop, etc.) by reading `~/.docker/config.json`. No manual socket configuration is needed in most setups.

### Runner image

The default image is [`myoung34/github-runner`](https://github.com/myoung34/docker-github-actions-runner), a widely used community image for self-hosted runners. Each container receives environment variables matching the image's expected configuration:

- `RUNNER_SCOPE`, `ORG_NAME` / `REPO_URL`
- `RUNNER_NAME`, `RUNNER_LABELS`, `RUNNER_GROUP`
- `ACCESS_TOKEN`, `EPHEMERAL`

## Project Structure

```
cmd/ghr/main.go               # entry point
internal/
  cli/                         # cobra commands (one file per command)
  config/                      # config loading, saving, validation, token resolution
  docker/                      # Docker client wrapper, labels, container management
  github/                      # GitHub API client, runner status queries
  runner/                      # orchestration (up/down/scale), naming, types
  output/                      # table formatting
```

## Development

```bash
make build    # build to ./bin/ghr
make test     # run tests
make vet      # run go vet
make install  # install to $GOPATH/bin
```

## License

MIT
