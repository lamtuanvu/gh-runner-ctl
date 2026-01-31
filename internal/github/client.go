package github

import (
	"context"

	gh "github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client.
type Client struct {
	gh *gh.Client
}

// NewClient creates a GitHub API client with the given token.
func NewClient(ctx context.Context, token string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &Client{gh: gh.NewClient(tc)}
}
