package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStopCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "stop [NAME_OR_NUMBER]",
		Short: "Stop runners without removing",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !all && len(args) == 0 {
				return fmt.Errorf("specify a runner or use --all")
			}
			nameOrNum := ""
			if len(args) > 0 {
				nameOrNum = args[0]
			}
			return mgr.Stop(cmd.Context(), nameOrNum, all)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "stop all managed runners")
	return cmd
}
