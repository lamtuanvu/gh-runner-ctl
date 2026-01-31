---
title: ghr status
weight: 7
---

Show config summary and runner counts.

## Synopsis

```
ghr status
```

## Description

Displays a summary of the current ghr configuration and the state of all managed runners, including counts of running and stopped containers.

## Examples

```bash
ghr status
```

```
Config:   ~/.ghr/config.yaml
Scope:    org
Org:      my-org
Image:    myoung34/github-runner:latest
Labels:   local, dev
Runners:  10 total (10 running, 0 stopped)
```

## Related Commands

- [`ghr list`](../list) -- detailed per-runner listing
- [`ghr init`](../init) -- create or update configuration
