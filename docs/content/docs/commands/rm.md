---
title: ghr rm
weight: 10
---

Remove runner containers.

## Synopsis

```
ghr rm [NAME_OR_NUMBER | --all]
```

## Description

Removes one or more managed runner containers. Running containers are force-removed. Use [`ghr stop`](../stop) first if you want a graceful shutdown.

Either a runner identifier or `--all` must be provided.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `NAME_OR_NUMBER` | No | Runner number, name, or container ID prefix. |

## Flags

| Flag | Description |
|------|-------------|
| `--all` | Remove all managed runners |

## Examples

Remove runner #3:

```bash
ghr rm 3
```

Remove all managed runners:

```bash
ghr rm --all
```

## Related Commands

- [`ghr stop`](../stop) -- stop without removing
- [`ghr down`](../down) -- stop and remove in one step
- [`ghr list`](../list) -- see which runners exist
