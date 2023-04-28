package github

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
)

func CheckGithubAppConnection() (string, error) {
	privateKeyPath := "data-drift.2023-04-28.private-key.pem"

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 3252700, 36944435, privateKeyPath)

	if err != nil {
		return "", err
	}

	// Use installation transport with client.
	client := github.NewClient(&http.Client{Transport: itr})
	ctx := context.Background()

	// Get the last commit of the repository
	commit, _, err := client.Repositories.GetCommit(ctx, "Samox", "copy-libeo-data-history", "main")

	if err != nil {
		return "", err
	}
	return commit.GetSHA(), nil
}
