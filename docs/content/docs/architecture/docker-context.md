---
title: Docker Context Detection
weight: 3
---

ghr automatically detects your active Docker CLI context so it works out of the box with various Docker environments.

## Supported Environments

ghr works with any Docker-compatible runtime:

- **Docker Desktop** (macOS, Windows, Linux)
- **OrbStack** (macOS)
- **colima** (macOS, Linux)
- **Rancher Desktop**
- **Podman** (with Docker-compatible socket)
- Any environment providing a Docker socket

## How Detection Works

ghr resolves the Docker socket endpoint in the following order:

1. **`DOCKER_HOST` environment variable** -- if set, ghr uses this directly
2. **`docker.socket` config field** -- if set in `config.yaml`
3. **Active Docker CLI context** -- ghr reads `~/.docker/config.json` to find the `currentContext`, then looks up the endpoint metadata for that context

### Context Resolution

Docker CLI contexts store metadata in `~/.docker/contexts/`. ghr reads the active context by:

1. Reading `currentContext` from `~/.docker/config.json`
2. Hashing the context name with SHA256 to locate its metadata directory
3. Extracting the Docker endpoint from the context metadata file

This is the same resolution logic that `docker` CLI uses, so ghr connects to the same Docker daemon.

## Manual Configuration

If auto-detection doesn't work for your setup, you can set the socket path explicitly:

```yaml
docker:
  socket: /var/run/docker.sock
```

Or use the `DOCKER_HOST` environment variable:

```bash
export DOCKER_HOST=unix:///var/run/docker.sock
ghr list
```

## Common Socket Paths

| Runtime | Typical Socket Path |
|---------|-------------------|
| Docker Desktop (macOS) | `/var/run/docker.sock` or `~/.docker/run/docker.sock` |
| Docker Desktop (Linux) | `/var/run/docker.sock` |
| OrbStack | `~/.orbstack/run/docker.sock` |
| colima | `~/.colima/default/docker.sock` |
| Rancher Desktop | `~/.rd/docker.sock` |
