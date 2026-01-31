package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "start [NAME_OR_NUMBER]",
		Short: "Start stopped runners",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !all && len(args) == 0 {
				return fmt.Errorf("specify a runner or use --all")
			}
			nameOrNum := ""
			if len(args) > 0 {
				nameOrNum = args[0]
			}
			return mgr.Start(cmd.Context(), nameOrNum, all)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "start all managed runners")
	return cmd
}
