---
title: ghr completion
weight: 11
---

Generate shell completion scripts.

## Synopsis

```
ghr completion <shell>
ghr completion install
```

## Description

Generates autocompletion scripts for the specified shell. Supports bash, zsh, fish, and PowerShell. The `install` subcommand auto-detects your shell and writes the completion file to the appropriate location.

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `bash` | Generate bash completion script |
| `zsh` | Generate zsh completion script |
| `fish` | Generate fish completion script |
| `powershell` | Generate PowerShell completion script |
| `install` | Auto-detect shell and install completion |

## Examples

Auto-install completions for your current shell:

```bash
ghr completion install
```

Generate bash completion and source it manually:

```bash
source <(ghr completion bash)
```

Generate zsh completion:

```bash
ghr completion zsh > "${fpath[1]}/_ghr"
```

Generate fish completion:

```bash
ghr completion fish > ~/.config/fish/completions/ghr.fish
```

## Installation Paths

When using `ghr completion install`, the completion file is written to:

| Shell | Path |
|-------|------|
| Bash | `~/.bash_completion.d/ghr` |
| Zsh | `~/.zsh/completions/_ghr` |
| Fish | `~/.config/fish/completions/ghr.fish` |

## Related Commands

- [`ghr version`](../version) -- verify installation
