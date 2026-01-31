---
title: ghr start
weight: 9
---

Start stopped runners.

## Synopsis

```
ghr start [NAME_OR_NUMBER | --all]
```

## Description

Starts one or more previously stopped runner containers. The containers must already exist (use [`ghr up`](../up) to create new runners).

Either a runner identifier or `--all` must be provided.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `NAME_OR_NUMBER` | No | Runner number, name, or container ID prefix. |

## Flags

| Flag | Description |
|------|-------------|
| `--all` | Start all stopped managed runners |

## Examples

Start runner #3:

```bash
ghr start 3
```

Start all managed runners:

```bash
ghr start --all
```

## Related Commands

- [`ghr stop`](../stop) -- stop running runners
- [`ghr up`](../up) -- create new runners
- [`ghr list`](../list) -- see which runners are stopped
