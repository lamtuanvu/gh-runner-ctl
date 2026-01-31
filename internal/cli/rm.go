package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRmCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "rm [NAME_OR_NUMBER]",
		Short: "Remove stopped runners",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !all && len(args) == 0 {
				return fmt.Errorf("specify a runner or use --all")
			}
			nameOrNum := ""
			if len(args) > 0 {
				nameOrNum = args[0]
			}
			return mgr.Remove(cmd.Context(), nameOrNum, all)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "remove all managed runners")
	return cmd
}
