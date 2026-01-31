---
title: Docker Labels
weight: 1
---

ghr uses Docker labels to track which containers it manages. This design eliminates the need for state files or databases.

## Label Schema

Every container created by ghr is tagged with the following labels:

| Label | Example Value | Description |
|-------|---------------|-------------|
| `dev.ghr.managed` | `true` | Identifies the container as ghr-managed. All ghr operations filter on this label. |
| `dev.ghr.runner-num` | `3` | The runner's sequential number, used for naming and ordering. |
| `dev.ghr.scope` | `org` | The scope at creation time (`org` or `repo`). |
| `dev.ghr.org` | `my-org` | The GitHub organization name (set when scope is `org`). |
| `dev.ghr.repo-owner` | `lamtuanvu` | The repository owner (set when scope is `repo`). |
| `dev.ghr.repo-name` | `gh-runner-ctl` | The repository name (set when scope is `repo`). |

## How It Works

When ghr lists, stops, or removes containers, it queries the Docker daemon with a label filter for `dev.ghr.managed=true`. This means:

- ghr only touches containers it created
- Other Docker containers are completely unaffected
- No state file can become stale or corrupted
- Multiple ghr instances can coexist (though they should manage different scopes)

## Stateless Design

Traditional runner managers maintain a state file that maps runner IDs to container IDs. If this file gets out of sync (e.g., a container is manually removed), the manager enters an inconsistent state.

ghr avoids this problem entirely. Every operation queries the Docker daemon for the current truth. If a container was manually removed via `docker rm`, ghr simply won't see it in the next listing -- no cleanup or reconciliation needed.

## Label Filtering

ghr uses the Docker API's label filter to efficiently query only its own containers:

```
label=dev.ghr.managed=true
```

This is the equivalent of:

```bash
docker ps --filter "label=dev.ghr.managed=true"
```

## Inspecting Labels

You can inspect the labels on any ghr-managed container using standard Docker commands:

```bash
docker inspect --format '{{json .Config.Labels}}' ghr-runner-1
```

```json
{
  "dev.ghr.managed": "true",
  "dev.ghr.runner-num": "1",
  "dev.ghr.scope": "org",
  "dev.ghr.org": "my-org"
}
```
