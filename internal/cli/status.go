package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/lamtuanvu/gh-runner-ctl/internal/output"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show config summary and runner counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			runners, err := mgr.List(cmd.Context())
			if err != nil {
				return err
			}

			running := 0
			stopped := 0
			for _, r := range runners {
				if r.DockerState == "running" {
					running++
				} else {
					stopped++
				}
			}

			output.PrintStatusSummary(os.Stdout, cfgPath,
				cfg.Scope, cfg.Org, cfg.Repo.Owner, cfg.Repo.Name,
				cfg.Runners.Image, cfg.Runners.Labels,
				len(runners), running, stopped,
			)
			return nil
		},
	}
}
