---
title: Installation
weight: 1
---

ghr is distributed as a single binary. Choose the installation method that fits your environment.

## Homebrew (macOS / Linux)

The recommended method for macOS and Linux systems with Homebrew:

```bash
brew install lamtuanvu/tap/ghr
```

This installs the latest release and makes `ghr` available on your `PATH`.

## Debian / Ubuntu

Download the `.deb` package from the [latest release](https://github.com/lamtuanvu/gh-runner-ctl/releases/latest):

```bash
curl -LO https://github.com/lamtuanvu/gh-runner-ctl/releases/latest/download/ghr_<VERSION>_linux_amd64.deb
sudo dpkg -i ghr_<VERSION>_linux_amd64.deb
```

Replace `<VERSION>` with the actual version number (e.g., `0.1.0`).

## Install Script

A portable install script that detects your OS and architecture:

```bash
curl -sSL https://raw.githubusercontent.com/lamtuanvu/gh-runner-ctl/main/install.sh | sh
```

To install a specific version or to a custom directory:

```bash
curl -sSL https://raw.githubusercontent.com/lamtuanvu/gh-runner-ctl/main/install.sh | sh -s -- -v v0.1.0 -d ~/.local/bin
```

| Flag | Description |
|------|-------------|
| `-v VERSION` | Install a specific version (e.g., `-v v0.1.0`) |
| `-d DIRECTORY` | Install to a custom directory (default: `/usr/local/bin`) |

## Build from Source

Requires **Go 1.21+** and **Git**.

```bash
git clone https://github.com/lamtuanvu/gh-runner-ctl.git
cd gh-runner-ctl
make build        # outputs to ./bin/ghr
make install      # installs to $GOPATH/bin
```

The `make build` target injects the version from the latest git tag.

## Verify Installation

After installing, verify that ghr is available:

```bash
ghr version
```

## Prerequisites

ghr requires a running Docker daemon. It works with:

- **Docker Desktop** (macOS, Windows, Linux)
- **OrbStack** (macOS)
- **colima** (macOS, Linux)
- Any Docker-compatible runtime accessible via the Docker socket

ghr auto-detects the active Docker context, so no manual socket configuration is needed in most setups. See [Docker Context Detection](../../architecture/docker-context) for details.
