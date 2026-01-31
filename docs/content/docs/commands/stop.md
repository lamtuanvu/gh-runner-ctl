---
title: ghr stop
weight: 8
---

Stop runners without removing them.

## Synopsis

```
ghr stop [NAME_OR_NUMBER | --all]
```

## Description

Stops one or more running containers without removing them. Stopped containers retain their configuration and can be restarted with [`ghr start`](../start).

Either a runner identifier or `--all` must be provided.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `NAME_OR_NUMBER` | No | Runner number, name, or container ID prefix. |

## Flags

| Flag | Description |
|------|-------------|
| `--all` | Stop all managed runners |

## Examples

Stop runner #3:

```bash
ghr stop 3
```

Stop all managed runners:

```bash
ghr stop --all
```

## Related Commands

- [`ghr start`](../start) -- restart stopped runners
- [`ghr rm`](../rm) -- remove stopped runners
- [`ghr down`](../down) -- stop and remove in one step
