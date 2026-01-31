package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/lamtuanvu/gh-runner-ctl/internal/runner"
)

// PrintRunnerTable prints a formatted table of runner info.
func PrintRunnerTable(w io.Writer, runners []runner.RunnerInfo, showGitHub bool) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if showGitHub {
		fmt.Fprintln(tw, "NUM\tNAME\tCONTAINER\tDOCKER STATUS\tGITHUB\tBUSY")
		fmt.Fprintln(tw, "---\t----\t---------\t-------------\t------\t----")
		for _, r := range runners {
			busy := ""
			if r.GitHubStatus != "" {
				if r.Busy {
					busy = "yes"
				} else {
					busy = "no"
				}
			}
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n",
				r.Num, r.Name, r.ContainerID, statusWithState(r), r.GitHubStatus, busy)
		}
	} else {
		fmt.Fprintln(tw, "NUM\tNAME\tCONTAINER\tSTATUS")
		fmt.Fprintln(tw, "---\t----\t---------\t------")
		for _, r := range runners {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n",
				r.Num, r.Name, r.ContainerID, statusWithState(r))
		}
	}
	tw.Flush()
}

func statusWithState(r runner.RunnerInfo) string {
	if r.DockerStatus != "" {
		return r.DockerStatus
	}
	return r.DockerState
}

// PrintStatusSummary prints the ghr status overview.
func PrintStatusSummary(w io.Writer, configPath string, scope, org, repoOwner, repoName string,
	image string, labels []string, total, running, stopped int) {

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Config:\t%s\n", configPath)
	fmt.Fprintf(tw, "Scope:\t%s\n", scope)
	if scope == "org" {
		fmt.Fprintf(tw, "Org:\t%s\n", org)
	} else {
		fmt.Fprintf(tw, "Repo:\t%s/%s\n", repoOwner, repoName)
	}
	fmt.Fprintf(tw, "Image:\t%s\n", image)
	fmt.Fprintf(tw, "Labels:\t%s\n", strings.Join(labels, ", "))
	fmt.Fprintf(tw, "Runners:\t%d total (%d running, %d stopped)\n", total, running, stopped)
	tw.Flush()
}
