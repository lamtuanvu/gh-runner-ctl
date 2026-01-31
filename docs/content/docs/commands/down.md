---
title: ghr down
weight: 3
---

Stop and remove runners.

## Synopsis

```
ghr down [COUNT | --all]
```

## Description

Stops and removes COUNT runners, starting with the **highest-numbered** runners first. Use `--all` to remove every managed runner.

Either COUNT or `--all` must be provided.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `COUNT` | No | Number of runners to remove (highest-numbered first). Must be >= 1. |

## Flags

| Flag | Description |
|------|-------------|
| `--all` | Remove all managed runners |

## Examples

Remove 3 runners (highest-numbered first):

```bash
ghr down 3
```

Remove all managed runners:

```bash
ghr down --all
```

## Related Commands

- [`ghr up`](../up) -- create new runners
- [`ghr scale`](../scale) -- scale to an exact count
- [`ghr rm`](../rm) -- remove already-stopped runners
