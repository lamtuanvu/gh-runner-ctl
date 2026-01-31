---
title: Org vs Repo Runners
weight: 1
---

ghr supports registering runners at two levels: **organization** and **repository**. This guide explains the differences and when to use each.

## Organization Runners

Organization-level runners are shared across all repositories in the organization (subject to runner group restrictions).

```yaml
scope: org
org: my-org
```

**When to use:**
- You want a shared pool of runners for multiple repositories
- You prefer centralized runner management
- You use runner groups to control access per team or repo

**Token scope required:** `admin:org`

### Runner Groups

Organization runners can be assigned to [runner groups](https://docs.github.com/en/actions/hosting-your-own-runners/managing-self-hosted-runners/managing-access-to-self-hosted-runners-using-groups). Set the group in your config:

```yaml
runners:
  group: Production
```

The default group name is `Default`, which is available to all repositories in the organization.

## Repository Runners

Repository-level runners are scoped to a single repository.

```yaml
scope: repo
repo:
  owner: lamtuanvu
  name: gh-runner-ctl
```

**When to use:**
- You need runners dedicated to a specific repository
- You don't have organization admin access
- You want strict isolation between projects

**Token scope required:** `repo`

## Comparison

| Aspect | Organization | Repository |
|--------|-------------|------------|
| Scope | All repos in org (via groups) | Single repo |
| Token scope | `admin:org` | `repo` |
| Runner groups | Supported | Not applicable |
| Config fields | `scope: org`, `org` | `scope: repo`, `repo.owner`, `repo.name` |
| Best for | Teams, shared CI | Individual projects, isolation |

## Switching Scope

To switch from org to repo runners (or vice versa):

1. Remove existing runners: `ghr down --all`
2. Update `config.yaml` with the new scope and credentials
3. Create new runners: `ghr up`

{{< callout type="warning" >}}
Always remove runners before switching scope. Runners registered under the old scope won't be deregistered automatically.
{{< /callout >}}
