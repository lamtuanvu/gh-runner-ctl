---
title: Token Setup
weight: 2
---

ghr needs a GitHub Personal Access Token to register runners with the GitHub API.

## Token Formats

The `token` field in `config.yaml` supports two formats:

| Format | Example | Description |
|--------|---------|-------------|
| `env:VAR` | `env:GH_TOKEN` | Reads the token from an environment variable at runtime. **Recommended.** |
| Literal | `ghp_abc123...` | Token stored directly in the config file. Avoid for shared or version-controlled configs. |

### Environment Variable Reference (recommended)

```yaml
token: env:GH_TOKEN
```

At runtime, ghr reads the value of `GH_TOKEN` from the process environment. This keeps the token out of the config file and works well with `.env` files, CI secrets, and secret managers.

### Literal Token

```yaml
token: ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

The token is stored directly in the config file. This is convenient for local-only setups but should be avoided if the config might be shared or committed to version control.

## Required Scopes

The token must be a **Personal Access Token (classic)**. Fine-grained tokens are not currently supported by the runner registration API.

| Scope | When Required |
|-------|---------------|
| `admin:org` | Registering organization-level runners (`scope: org`) |
| `repo` | Registering repository-level runners (`scope: repo`) |

Create a token at [github.com/settings/tokens](https://github.com/settings/tokens).

## Using a .env File

ghr automatically loads environment variables from `~/.ghr/.env` before resolving the token. The `.env` file uses standard `KEY=VALUE` format:

```bash
# ~/.ghr/.env
GH_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

- Lines starting with `#` are treated as comments
- Existing environment variables are **not** overwritten (the process environment takes precedence)
- Inline comments after values are stripped

During `ghr init`, if a `.env` file is found in `~/.ghr/` or the current directory, ghr offers to import settings from it.

## Security Recommendations

1. **Use `env:VAR` syntax** -- keeps the token out of the config file
2. **Use a `.env` file** -- store the token in `~/.ghr/.env` with restrictive permissions:
   ```bash
   chmod 600 ~/.ghr/.env
   ```
3. **Never commit tokens** -- add `.env` and `.ghr.yaml` to `.gitignore`
4. **Rotate tokens regularly** -- generate a new token periodically and update your `.env` file
5. **Use minimal scopes** -- only grant `admin:org` or `repo` as needed, not both
