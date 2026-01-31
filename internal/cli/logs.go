package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func newLogsCmd() *cobra.Command {
	var follow bool

	cmd := &cobra.Command{
		Use:   "logs NAME_OR_NUMBER",
		Short: "Show runner container logs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Resolve name: could be a number, name, or container ID prefix
			nameOrNum := args[0]

			// Find the container to get its full name
			runners, err := mgr.List(cmd.Context())
			if err != nil {
				return err
			}

			var containerName string
			for _, r := range runners {
				if r.Name == nameOrNum ||
					fmt.Sprintf("%d", r.Num) == nameOrNum ||
					r.ContainerID == nameOrNum {
					containerName = r.Name
					break
				}
			}
			if containerName == "" {
				// Try using it directly as a container name/ID
				containerName = nameOrNum
			}

			reader, err := dockerCli.ContainerLogs(cmd.Context(), containerName, follow)
			if err != nil {
				return fmt.Errorf("getting logs for %s: %w", containerName, err)
			}
			defer reader.Close()

			_, err = io.Copy(os.Stdout, reader)
			return err
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow log output")
	return cmd
}
