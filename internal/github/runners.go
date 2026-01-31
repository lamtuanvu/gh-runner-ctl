package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v68/github"
)

// RunnerStatus holds simplified GitHub runner info.
type RunnerStatus struct {
	ID     int64
	Name   string
	Status string // online, offline
	Busy   bool
	Labels []string
}

// ListOrgRunners lists all self-hosted runners for an organization.
func (c *Client) ListOrgRunners(ctx context.Context, org string) ([]RunnerStatus, error) {
	var all []RunnerStatus
	opts := &gh.ListOptions{PerPage: 100}

	for {
		runners, resp, err := c.gh.Actions.ListOrganizationRunners(ctx, org, &gh.ListRunnersOptions{ListOptions: *opts})
		if err != nil {
			return nil, fmt.Errorf("listing org runners: %w", err)
		}
		for _, r := range runners.Runners {
			var labels []string
			for _, l := range r.Labels {
				labels = append(labels, l.GetName())
			}
			all = append(all, RunnerStatus{
				ID:     r.GetID(),
				Name:   r.GetName(),
				Status: r.GetStatus(),
				Busy:   r.GetBusy(),
				Labels: labels,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return all, nil
}

// ListRepoRunners lists all self-hosted runners for a repository.
func (c *Client) ListRepoRunners(ctx context.Context, owner, repo string) ([]RunnerStatus, error) {
	var all []RunnerStatus
	opts := &gh.ListOptions{PerPage: 100}

	for {
		runners, resp, err := c.gh.Actions.ListRunners(ctx, owner, repo, &gh.ListRunnersOptions{ListOptions: *opts})
		if err != nil {
			return nil, fmt.Errorf("listing repo runners: %w", err)
		}
		for _, r := range runners.Runners {
			var labels []string
			for _, l := range r.Labels {
				labels = append(labels, l.GetName())
			}
			all = append(all, RunnerStatus{
				ID:     r.GetID(),
				Name:   r.GetName(),
				Status: r.GetStatus(),
				Busy:   r.GetBusy(),
				Labels: labels,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return all, nil
}
