---
title: ghr scale
weight: 4
---

Scale to an exact runner count.

## Synopsis

```
ghr scale COUNT
```

## Description

Adjusts the number of managed runners to exactly COUNT. If the current count is lower, new runners are created. If the current count is higher, excess runners are removed (highest-numbered first).

`ghr scale 0` is equivalent to `ghr down --all`.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `COUNT` | **Yes** | Target number of runners. Must be >= 0. |

## Examples

Scale to exactly 10 runners:

```bash
ghr scale 10
```

Scale down to zero (remove all):

```bash
ghr scale 0
```

## Related Commands

- [`ghr up`](../up) -- add runners without affecting existing ones
- [`ghr down`](../down) -- remove a specific number of runners
- [`ghr list`](../list) -- see the current runner count
