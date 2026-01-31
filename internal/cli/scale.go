package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newScaleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scale COUNT",
		Short: "Scale to exactly COUNT runners",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := strconv.Atoi(args[0])
			if err != nil || target < 0 {
				return fmt.Errorf("invalid count: %s", args[0])
			}
			return mgr.Scale(cmd.Context(), target)
		},
	}
}
