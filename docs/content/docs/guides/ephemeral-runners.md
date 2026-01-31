---
title: Ephemeral Runners
weight: 2
---

Ephemeral runners are self-hosted runners that automatically de-register after completing a single job. This provides a clean environment for every workflow run.

## Enabling Ephemeral Mode

Ephemeral mode is enabled by default in ghr:

```yaml
runners:
  ephemeral: true
```

When enabled, ghr sets the `EPHEMERAL=true` environment variable on each runner container. The runner image handles deregistration after the first job completes.

## How It Works

1. ghr creates a container with `EPHEMERAL=true`
2. The runner registers with GitHub and picks up one job
3. After the job completes, the runner deregisters from GitHub
4. The container exits

The container remains in a stopped state after exit. Depending on the restart policy, it may or may not restart.

## Restart Policy for Ephemeral Runners

The `docker.restart_policy` setting controls what happens after the runner exits:

| Policy | Behavior with Ephemeral |
|--------|------------------------|
| `"unless-stopped"` (default) | Container restarts, re-registers, and picks up another job. Effectively makes the runner **persistent but clean**. |
| `"no"` | Container stays stopped after one job. Requires manual restart or `ghr up` to create new runners. |
| `"always"` | Same as `unless-stopped` for ephemeral runners. |

### Recommended Configuration

For a pool that stays alive but starts fresh for each job:

```yaml
runners:
  ephemeral: true
docker:
  restart_policy: unless-stopped
```

For true single-use runners that must be explicitly recreated:

```yaml
runners:
  ephemeral: true
docker:
  restart_policy: "no"
```

## Benefits

- **Clean environment**: Each job gets a fresh runner with no leftover state from previous jobs
- **Security**: Prevents data leakage between workflow runs
- **Predictability**: No accumulated disk usage, cached dependencies, or stale processes

## Trade-offs

- **Startup time**: Each job incurs runner registration overhead (typically a few seconds)
- **Restart management**: With `restart_policy: "no"`, you need to create new runners after they're consumed

## Non-Ephemeral Runners

To disable ephemeral mode, set:

```yaml
runners:
  ephemeral: false
```

Non-ephemeral runners persist across jobs. They stay registered with GitHub and continuously pick up new jobs without restarting.
