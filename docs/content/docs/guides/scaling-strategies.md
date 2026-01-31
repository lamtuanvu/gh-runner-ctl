---
title: Scaling Strategies
weight: 3
---

ghr provides several commands for managing runner count. This guide covers common scaling patterns.

## Fixed Pool

The simplest strategy: maintain a fixed number of runners at all times.

```bash
ghr up 10
```

Use `ghr scale` to adjust the pool size when needs change:

```bash
ghr scale 20    # grow to 20
ghr scale 5     # shrink to 5
```

`ghr scale` is idempotent -- running it multiple times with the same number has no effect.

## Scale-to-Zero

Tear down all runners when not needed (e.g., outside business hours or after a CI burst):

```bash
ghr down --all
# or equivalently
ghr scale 0
```

Bring them back when needed:

```bash
ghr up 10
```

## Additive Scaling

`ghr up` is additive -- it creates new runners without affecting existing ones. This is useful for temporarily adding capacity:

```bash
# Start with 5 runners
ghr up 5

# Burst: add 15 more for a large job queue
ghr up 15

# After the burst, scale back to 5
ghr scale 5
```

## Runner Number Assignment

ghr assigns the **lowest available** numbers when creating runners. When removing runners, it removes the **highest-numbered** first. This keeps numbering compact.

Example:

```bash
ghr up 5          # creates runners 1-5
ghr down 2        # removes runners 5, 4
ghr up 1          # creates runner 4 (fills the gap)
```

## Monitoring Runner State

Use `ghr list --github` to see which runners are busy:

```bash
ghr list --github
```

```
NUM  NAME          CONTAINER     DOCKER STATUS  GITHUB   BUSY
---  ----          ---------     -------------  ------   ----
1    ghr-runner-1  c5904092bff3  Up 10 minutes  online   yes
2    ghr-runner-2  628e83bedc4b  Up 10 minutes  online   no
3    ghr-runner-3  a1b2c3d4e5f6  Up 10 minutes  online   yes
```

If most runners show `busy: yes`, it may be time to scale up.

## Labels for Workload Routing

Use labels to route different workflow types to different runner pools:

```yaml
# Pool 1: general CI
runners:
  labels:
    - ci
    - general
  name_prefix: ci
```

```yaml
# Pool 2: GPU workloads (separate config)
runners:
  labels:
    - gpu
    - ml
  name_prefix: gpu
```

In your workflow:

```yaml
jobs:
  build:
    runs-on: [self-hosted, ci]
  train:
    runs-on: [self-hosted, gpu]
```

{{< callout type="info" >}}
Each label pool requires its own config file. Use `--config` to manage multiple pools:
```bash
ghr --config ~/.ghr/ci.yaml up 10
ghr --config ~/.ghr/gpu.yaml up 2
```
{{< /callout >}}

## Tips

- **Start small** -- begin with a few runners and scale up based on queue depth
- **Use ephemeral mode** for clean-room builds where security matters
- **Monitor with `ghr status`** -- check running vs stopped counts regularly
- **Use `ghr scale`** over `ghr up`/`ghr down` when you want a specific count -- it's easier to reason about
