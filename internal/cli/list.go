package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/lamtuanvu/gh-runner-ctl/internal/config"
	ghclient "github.com/lamtuanvu/gh-runner-ctl/internal/github"
	"github.com/lamtuanvu/gh-runner-ctl/internal/output"
)

func newListCmd() *cobra.Command {
	var showGitHub bool

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List managed runners",
		RunE: func(cmd *cobra.Command, args []string) error {
			runners, err := mgr.List(cmd.Context())
			if err != nil {
				return err
			}
			if len(runners) == 0 {
				fmt.Println("No managed runners found.")
				return nil
			}

			if showGitHub {
				token, err := config.ResolveToken(cfg.Token)
				if err != nil {
					return fmt.Errorf("resolving token for GitHub API: %w", err)
				}
				ghc := ghclient.NewClient(cmd.Context(), token)

				var ghRunners []ghclient.RunnerStatus
				if cfg.Scope == "org" {
					ghRunners, err = ghc.ListOrgRunners(cmd.Context(), cfg.Org)
				} else {
					ghRunners, err = ghc.ListRepoRunners(cmd.Context(), cfg.Repo.Owner, cfg.Repo.Name)
				}
				if err != nil {
					return fmt.Errorf("fetching GitHub runner status: %w", err)
				}

				// Build lookup by name
				ghByName := make(map[string]ghclient.RunnerStatus)
				for _, gr := range ghRunners {
					ghByName[gr.Name] = gr
				}

				// Merge
				for i := range runners {
					if gr, ok := ghByName[runners[i].Name]; ok {
						runners[i].GitHubID = gr.ID
						runners[i].GitHubStatus = gr.Status
						runners[i].Busy = gr.Busy
					}
				}
			}

			output.PrintRunnerTable(os.Stdout, runners, showGitHub)
			return nil
		},
	}

	cmd.Flags().BoolVar(&showGitHub, "github", false, "show GitHub API runner status")
	return cmd
}
