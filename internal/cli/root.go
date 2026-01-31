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

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "ghr",
		Short: "GitHub self-hosted runner manager",
		Long:  "Manage GitHub Actions self-hosted runners via Docker containers.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Commands that don't need config
			if cmd.Name() == "init" || cmd.Name() == "version" {
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
	}

	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .ghr.yaml)")

	root.AddCommand(
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
