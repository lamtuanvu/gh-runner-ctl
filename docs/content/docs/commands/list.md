---
title: ghr list
weight: 5
---

List managed runner containers.

## Synopsis

```
ghr list [flags]
```

**Alias:** `ghr ls`

## Description

Lists all Docker containers managed by ghr, showing their runner number, name, container ID, and Docker status.

With the `--github` flag, ghr also queries the GitHub API to show each runner's online/offline status and whether it is currently busy executing a job.

## Flags

| Flag | Description |
|------|-------------|
| `--github` | Also show GitHub API runner status (online/offline, busy) |

## Examples

List Docker container status:

```bash
ghr list
```

```
NUM  NAME          CONTAINER     STATUS
---  ----          ---------     ------
1    ghr-runner-1  c5904092bff3  Up 2 minutes
2    ghr-runner-2  628e83bedc4b  Up 2 minutes
```

Include GitHub API status:

```bash
ghr list --github
```

```
NUM  NAME          CONTAINER     DOCKER STATUS  GITHUB   BUSY
---  ----          ---------     -------------  ------   ----
1    ghr-runner-1  c5904092bff3  Up 2 minutes   online   no
2    ghr-runner-2  628e83bedc4b  Up 2 minutes   online   yes
```

## Related Commands

- [`ghr status`](../status) -- summary view with runner counts
- [`ghr logs`](../logs) -- view logs for a specific runner
