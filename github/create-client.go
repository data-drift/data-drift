package github

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func CreateClientFromGithubApp() {}

func CreateClientFromGithubToken(token string) *github.Client {
	client := github.NewClient(nil)
	if strings.HasPrefix(token, "github_pat") {
		// Create a new GitHub client with authentication using the token.
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}
	return client
}
