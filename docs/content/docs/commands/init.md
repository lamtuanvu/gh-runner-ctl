---
title: ghr init
weight: 1
---

Interactive config setup. Creates `~/.ghr/config.yaml`.

## Synopsis

```
ghr init [flags]
```

## Description

Runs an interactive wizard that prompts for all configuration values and writes the result to `~/.ghr/config.yaml`. If the config file already exists, ghr asks for confirmation before overwriting.

If a `.env` file is found in `~/.ghr/` or the current directory, ghr offers to import settings from it automatically. The `--import-env` flag forces this import without prompting.

## Flags

| Flag | Description |
|------|-------------|
| `--import-env` | Force import settings from an existing `.env` file |

## Interactive Prompts

The wizard prompts for the following values (defaults shown in brackets):

1. **Scope** -- `org` or `repo` (default: `org`)
2. **Organization** -- GitHub org name (if scope is `org`)
3. **Repository owner / name** -- repo coordinates (if scope is `repo`)
4. **Token** -- GitHub token or `env:VAR` reference (default: `env:GH_TOKEN`)
5. **Image** -- Docker image (default: `myoung34/github-runner:latest`)
6. **Labels** -- comma-separated list (default: `local,dev`)
7. **Runner group** -- GitHub runner group (default: `Default`)
8. **Name prefix** -- container name prefix (default: `ghr`)

## Examples

Interactive setup:

```bash
ghr init
```

Force import from `.env`:

```bash
ghr init --import-env
```

## Related Commands

- [`ghr status`](../status) -- verify config after init
- [`ghr up`](../up) -- launch runners after init
