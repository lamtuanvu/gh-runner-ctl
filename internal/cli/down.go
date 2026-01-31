package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newDownCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "down [COUNT]",
		Short: "Stop and remove runners",
		Long:  "Stop and remove COUNT runners (highest-numbered first). Use --all to remove all.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !all && len(args) == 0 {
				return fmt.Errorf("specify COUNT or use --all")
			}
			count := 0
			if len(args) == 1 {
				var err error
				count, err = strconv.Atoi(args[0])
				if err != nil || count < 1 {
					return fmt.Errorf("invalid count: %s", args[0])
				}
			}
			return mgr.Down(cmd.Context(), count, all)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "remove all managed runners")
	return cmd
}
