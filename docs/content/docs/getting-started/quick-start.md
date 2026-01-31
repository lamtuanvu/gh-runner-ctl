---
title: Quick Start
weight: 2
---

This guide walks you through setting up ghr from scratch and launching your first self-hosted runners.

## 1. Prepare a GitHub Token

You need a GitHub Personal Access Token (classic) with the appropriate scopes:

- **Organization runners**: `admin:org` scope
- **Repository runners**: `repo` scope

Create a token at [github.com/settings/tokens](https://github.com/settings/tokens) and export it:

```bash
export GH_TOKEN=ghp_xxxxxxxxxxxx
```

See [Token Setup](../../configuration/token-setup) for details on token formats and security.

## 2. Initialize Configuration

Run the interactive setup wizard:

```bash
ghr init
```

ghr will prompt you for:

1. **Scope** -- `org` or `repo`
2. **Organization or repository** -- the target for runner registration
3. **Token** -- defaults to `env:GH_TOKEN` (reads from the environment variable)
4. **Image** -- Docker image for runners (default: `myoung34/github-runner:latest`)
5. **Labels** -- comma-separated labels attached to each runner
6. **Runner group** -- GitHub runner group name (default: `Default`)
7. **Name prefix** -- prefix for container/runner names (default: `ghr`)

The config is saved to `~/.ghr/config.yaml`.

{{< callout type="info" >}}
If a `.env` file exists in `~/.ghr/` or the current directory, ghr offers to import settings from it automatically. Use `ghr init --import-env` to force import.
{{< /callout >}}

## 3. Launch Runners

Start 5 runners:

```bash
ghr up 5
```

This creates 5 Docker containers, each registering as a self-hosted runner with GitHub. Runner numbers fill the lowest available gaps, so you get runners 1 through 5.

## 4. Check Status

List your running containers:

```bash
ghr list
```

```
NUM  NAME          CONTAINER     STATUS
---  ----          ---------     ------
1    ghr-runner-1  c5904092bff3  Up 2 minutes
2    ghr-runner-2  628e83bedc4b  Up 2 minutes
3    ghr-runner-3  a1b2c3d4e5f6  Up 2 minutes
4    ghr-runner-4  f6e5d4c3b2a1  Up 2 minutes
5    ghr-runner-5  1a2b3c4d5e6f  Up 2 minutes
```

To also see GitHub API status (online/offline, busy):

```bash
ghr list --github
```

View a summary of your setup:

```bash
ghr status
```

```
Config:   ~/.ghr/config.yaml
Scope:    org
Org:      my-org
Image:    myoung34/github-runner:latest
Labels:   local, dev
Runners:  5 total (5 running, 0 stopped)
```

## 5. Scale Up or Down

Scale to exactly 10 runners:

```bash
ghr scale 10
```

This adds 5 more runners (since 5 already exist). To scale back down:

```bash
ghr scale 3
```

This removes the 7 highest-numbered runners.

## 6. View Logs

Check the logs from a specific runner:

```bash
ghr logs 1           # by runner number
ghr logs 1 -f        # follow logs in real-time
```

## 7. Clean Up

When you're done, tear down all runners:

```bash
ghr down --all
```

This stops and removes all ghr-managed containers. Other Docker containers are not affected.

## Next Steps

- [Configuration Reference](../../configuration/config-file) -- all config fields with types and defaults
- [Commands Reference](../../commands/) -- full documentation for every command
- [Org vs Repo Runners](../../guides/org-vs-repo) -- choosing the right scope
- [Ephemeral Runners](../../guides/ephemeral-runners) -- single-use runner setup
