package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/lamtuanvu/gh-runner-ctl/internal/config"
	"github.com/lamtuanvu/gh-runner-ctl/internal/docker"
	"github.com/lamtuanvu/gh-runner-ctl/internal/runner"
)

var (
	Version    = "dev"
	cfgFile    string
	cfg        *config.Config
	cfgPath    string
	dockerCli  *docker.Client
	mgr        *runner.Manager
)

// skipConfigLoad returns true if the command (or any of its parents) is one
// that should run without loading config or connecting to Docker.
func skipConfigLoad(cmd *cobra.Command) bool {
	skip := map[string]bool{
		"init":               true,
		"version":            true,
		"completion":         true,
		"help":               true,
		"__complete":         true,
		"__completeNoDesc":   true,
	}
	for c := cmd; c != nil; c = c.Parent() {
		if skip[c.Name()] {
			return true
		}
	}
	return false
}

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "ghr",
		Short: "GitHub self-hosted runner manager",
		Long:  "Manage GitHub Actions self-hosted runners via Docker containers.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load .env file if present (does not overwrite existing env vars)
			// Check ~/.ghr/.env first, then ./.env
			config.LoadDotenv(config.DotenvPath())
			config.LoadDotenv(".env")

			// Commands that don't need config
			if skipConfigLoad(cmd) {
				return nil
			}

			var err error
			cfg, cfgPath, err = config.Load(cfgFile)
			if err != nil {
				return err
			}
			if err := config.Validate(cfg); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}

			dockerCli, err = docker.NewClient(cfg.Docker.Socket)
			if err != nil {
				return err
			}
			mgr = runner.NewManager(cfg, dockerCli)
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if dockerCli != nil {
				dockerCli.Close()
			}
		},
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.ghr/config.yaml)")

	root.AddCommand(
		newCompletionCmd(),
		newInitCmd(),
		newUpCmd(),
		newDownCmd(),
		newScaleCmd(),
		newListCmd(),
		newLogsCmd(),
		newStopCmd(),
		newStartCmd(),
		newRmCmd(),
		newStatusCmd(),
		newVersionCmd(),
	)

	return root
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
