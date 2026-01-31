---
title: ghr up
weight: 2
---

Create and start new runners.

## Synopsis

```
ghr up [COUNT]
```

## Description

Creates and starts COUNT new runner containers. This operation is **additive** -- it does not affect existing runners. New runner numbers fill the lowest available gaps. For example, if runners 1, 2, and 5 exist, `ghr up 2` creates runners 3 and 4.

If COUNT is omitted, the value from `runners.count` in the config file is used (default: 10).

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `COUNT` | No | Number of runners to create. Must be >= 1. Defaults to `runners.count` from config. |

## Examples

Create 5 runners:

```bash
ghr up 5
```

Create the default number of runners (from config):

```bash
ghr up
```

## Related Commands

- [`ghr down`](../down) -- stop and remove runners
- [`ghr scale`](../scale) -- scale to an exact count
- [`ghr list`](../list) -- verify created runners
