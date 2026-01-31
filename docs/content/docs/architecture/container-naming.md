---
title: Container Naming
weight: 2
---

ghr uses a predictable naming convention for runner containers that makes them easy to identify across Docker, GitHub, and ghr itself.

## Naming Format

```
{prefix}-runner-{N}
```

Where:
- **`{prefix}`** is the `runners.name_prefix` from config (default: `ghr`)
- **`{N}`** is the sequential runner number

For example, with the default prefix:

```
ghr-runner-1
ghr-runner-2
ghr-runner-3
```

## GitHub Runner Name

The container name is also registered as the **GitHub runner name**. This means the same name appears in:

- `ghr list` output
- `docker ps` output
- The GitHub Settings > Actions > Runners page

This 1:1 mapping makes it easy to correlate a runner across all three views.

## Number Assignment

When creating new runners, ghr assigns numbers by filling the **lowest available gaps**. For example:

- Runners 1, 2, and 5 exist
- `ghr up 2` creates runners **3** and **4** (not 6 and 7)

When removing runners, ghr removes the **highest-numbered** runners first. This keeps the numbering compact.

## Custom Prefix

Change the prefix in `config.yaml`:

```yaml
runners:
  name_prefix: ci
```

This produces containers named `ci-runner-1`, `ci-runner-2`, etc.

{{< callout type="warning" >}}
Changing the prefix after runners are already created will cause ghr to lose track of the old containers, since it uses the prefix as part of the naming convention. Remove existing runners before changing the prefix.
{{< /callout >}}
