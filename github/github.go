package github

import (
	"context"
)

func CheckGithubAppConnection() (string, error) {

	client, err := CreateClientFromGithubApp()
	if err != nil {
		return "", err
	}
	ctx := context.Background()

	// Get the last commit of the repository
	commit, _, err := client.Repositories.GetCommit(ctx, "Samox", "copy-libeo-data-history", "main")

	if err != nil {
		return "", err
	}
	return commit.GetSHA(), nil
}
