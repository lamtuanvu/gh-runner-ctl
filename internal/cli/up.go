package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up [COUNT]",
		Short: "Create and start runners",
		Long:  "Create and start COUNT new runners (additive). Defaults to runners.count from config.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			count := cfg.Runners.Count
			if len(args) == 1 {
				var err error
				count, err = strconv.Atoi(args[0])
				if err != nil || count < 1 {
					return fmt.Errorf("invalid count: %s", args[0])
				}
			}

			created, err := mgr.Up(cmd.Context(), count)
			if err != nil {
				return err
			}
			fmt.Printf("\n%d runner(s) created.\n", len(created))
			return nil
		},
	}
}
