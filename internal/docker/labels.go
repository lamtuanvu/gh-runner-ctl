package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/filters"
)

const (
	LabelManaged   = "dev.ghr.managed"
	LabelRunnerNum = "dev.ghr.runner-num"
	LabelScope     = "dev.ghr.scope"
	LabelOrg       = "dev.ghr.org"
	LabelRepoOwner = "dev.ghr.repo-owner"
	LabelRepoName  = "dev.ghr.repo-name"
)

// ManagedLabels returns the base labels for a ghr-managed container.
func ManagedLabels(scope, org, repoOwner, repoName string, num int) map[string]string {
	labels := map[string]string{
		LabelManaged:   "true",
		LabelRunnerNum: fmt.Sprintf("%d", num),
		LabelScope:     scope,
	}
	if scope == "org" {
		labels[LabelOrg] = org
	} else {
		labels[LabelRepoOwner] = repoOwner
		labels[LabelRepoName] = repoName
	}
	return labels
}

// ManagedFilter returns a Docker filter that matches ghr-managed containers.
func ManagedFilter() filters.Args {
	return filters.NewArgs(
		filters.Arg("label", LabelManaged+"=true"),
	)
}
