---
title: ghr logs
weight: 6
---

Show runner container logs.

## Synopsis

```
ghr logs NAME_OR_NUMBER [flags]
```

## Description

Displays the logs from a runner container. The runner can be identified by its number, name, or container ID prefix.

By default, shows the last 100 lines of logs. Use `--follow` to stream new log output in real-time.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `NAME_OR_NUMBER` | **Yes** | Runner number (e.g., `1`), name (e.g., `ghr-runner-3`), or container ID prefix. |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--follow` | `-f` | Follow (stream) log output |

## Examples

View logs for runner #1:

```bash
ghr logs 1
```

View logs by runner name:

```bash
ghr logs ghr-runner-3
```

Follow logs in real-time:

```bash
ghr logs 1 -f
```

## Related Commands

- [`ghr list`](../list) -- find runner numbers and names
- [`ghr status`](../status) -- overview of all runners
