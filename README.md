# ghr - GitHub Self-Hosted Runner Manager

A CLI tool for managing GitHub Actions self-hosted runners as Docker containers. Replaces verbose `docker-compose.yml` files with a single command-line interface.

**[Full Documentation](https://lamtuanvu.github.io/gh-runner-ctl/)**

## Quick Start

```bash
# Install
brew install lamtuanvu/tap/ghr

# One-time setup
export GH_TOKEN=ghp_xxxxxxxxxxxx
ghr init

# Launch 10 runners
ghr up 10

# Check status
ghr list

# Scale to 20
ghr scale 20

# Tear down
ghr down --all
```

## Installation

| Method | Command |
|--------|---------|
| **Homebrew** | `brew install lamtuanvu/tap/ghr` |
| **Debian/Ubuntu** | [Download .deb](https://github.com/lamtuanvu/gh-runner-ctl/releases/latest) |
| **Install script** | `curl -sSL https://raw.githubusercontent.com/lamtuanvu/gh-runner-ctl/main/install.sh \| sh` |
| **From source** | `git clone` + `make build` ([details](https://lamtuanvu.github.io/gh-runner-ctl/docs/getting-started/installation/)) |

## Commands

| Command | Description |
|---------|-------------|
| `ghr init` | Interactive config setup |
| `ghr up [COUNT]` | Create and start runners |
| `ghr down [COUNT \| --all]` | Stop and remove runners |
| `ghr scale COUNT` | Scale to exactly COUNT runners |
| `ghr list [--github]` | List managed runners |
| `ghr logs NAME_OR_NUMBER [-f]` | Show runner logs |
| `ghr status` | Show config summary |
| `ghr stop / start / rm` | Lifecycle management |
| `ghr completion` | Shell completions |
| `ghr version` | Print version |

See the [command reference](https://lamtuanvu.github.io/gh-runner-ctl/docs/commands/) for full details.

## How It Works

ghr manages runners as Docker containers using **labels** for stateless tracking -- no state files, no databases. It auto-detects your Docker context (Docker Desktop, OrbStack, colima) and works alongside other containers without interference.

Learn more in the [architecture docs](https://lamtuanvu.github.io/gh-runner-ctl/docs/architecture/).

## Development

```bash
make build    # build to ./bin/ghr
make test     # run tests
make vet      # run go vet
make install  # install to $GOPATH/bin
```

## License

MIT
