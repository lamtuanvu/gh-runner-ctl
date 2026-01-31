package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for ghr.

To load completions:

  bash:
    source <(ghr completion bash)

  zsh:
    source <(ghr completion zsh)

  fish:
    ghr completion fish | source

  powershell:
    ghr completion powershell | Invoke-Expression

Or run "ghr completion install" to install them permanently.`,
	}

	cmd.AddCommand(
		newCompletionBashCmd(),
		newCompletionZshCmd(),
		newCompletionFishCmd(),
		newCompletionPowershellCmd(),
		newCompletionInstallCmd(),
	)

	return cmd
}

func newCompletionBashCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bash",
		Short: "Generate bash completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenBashCompletionV2(os.Stdout, true)
		},
	}
}

func newCompletionZshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "Generate zsh completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenZshCompletion(os.Stdout)
		},
	}
}

func newCompletionFishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fish",
		Short: "Generate fish completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		},
	}
}

func newCompletionPowershellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "powershell",
		Short: "Generate powershell completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		},
	}
}

func newCompletionInstallCmd() *cobra.Command {
	var shell string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install completion script for your shell",
		Long:  `Detect your current shell and install the completion script to the appropriate location.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if shell == "" {
				shell = detectShell()
				if shell == "" {
					return fmt.Errorf("could not detect shell from $SHELL; use --shell to specify one (bash, zsh, fish)")
				}
			}
			return installCompletion(cmd.Root(), shell)
		},
	}

	cmd.Flags().StringVar(&shell, "shell", "", "shell type (bash, zsh, fish)")

	return cmd
}

func detectShell() string {
	s := os.Getenv("SHELL")
	if s == "" {
		return ""
	}
	base := filepath.Base(s)
	switch base {
	case "bash", "zsh", "fish":
		return base
	default:
		return ""
	}
}

func installCompletion(root *cobra.Command, shell string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	var dest string
	var generate func() error

	switch strings.ToLower(shell) {
	case "bash":
		dir := filepath.Join(home, ".bash_completion.d")
		dest = filepath.Join(dir, "ghr")
		generate = func() error {
			return writeCompletion(dest, func(f *os.File) error {
				return root.GenBashCompletionV2(f, true)
			})
		}
	case "zsh":
		dir := filepath.Join(home, ".zsh", "completions")
		dest = filepath.Join(dir, "_ghr")
		generate = func() error {
			return writeCompletion(dest, func(f *os.File) error {
				return root.GenZshCompletion(f)
			})
		}
	case "fish":
		dir := filepath.Join(home, ".config", "fish", "completions")
		dest = filepath.Join(dir, "ghr.fish")
		generate = func() error {
			return writeCompletion(dest, func(f *os.File) error {
				return root.GenFishCompletion(f, true)
			})
		}
	default:
		return fmt.Errorf("unsupported shell %q; supported: bash, zsh, fish", shell)
	}

	if err := generate(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Completion script installed to %s\n", dest)

	switch strings.ToLower(shell) {
	case "bash":
		fmt.Fprintln(os.Stderr, "Ensure your .bashrc sources files from ~/.bash_completion.d/")
		fmt.Fprintln(os.Stderr, `  e.g.: for f in ~/.bash_completion.d/*; do [ -f "$f" ] && source "$f"; done`)
	case "zsh":
		fmt.Fprintln(os.Stderr, "Ensure your .zshrc includes ~/.zsh/completions in fpath:")
		fmt.Fprintln(os.Stderr, "  e.g.: fpath=(~/.zsh/completions $fpath); autoload -Uz compinit && compinit")
	case "fish":
		fmt.Fprintln(os.Stderr, "Fish completions are loaded automatically from ~/.config/fish/completions/")
	}

	return nil
}

func writeCompletion(dest string, gen func(*os.File) error) error {
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("could not create directory %s: %w", dir, err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", dest, err)
	}
	defer f.Close()

	return gen(f)
}
