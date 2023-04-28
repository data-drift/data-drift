package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
)

func CheckGithubAppConnection() (string, error) {
	privateKeyPath := os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH")
	privateKey := os.Getenv("GITHUB_APP_PRIVATE_KEY")

	appIDStr := os.Getenv("GITHUB_APP_ID")
	githubAppId, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		return "", err
	}

	itr, err := CreateGithubTransport(privateKeyPath, privateKey, githubAppId, 36944435)
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

func CreateGithubTransport(privateKeyPath string, privateKey string, githubAppId int64, githubInstallationId int64) (*ghinstallation.Transport, error) {
	if privateKeyPath != "" {
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, githubAppId, githubInstallationId, privateKeyPath)
		if err != nil {
			return nil, err
		}
		return itr, nil
	} else if privateKey != "" {
		itr, err := ghinstallation.New(http.DefaultTransport, githubAppId, githubInstallationId, []byte(privateKey))
		if err != nil {
			return nil, err
		}
		return itr, nil
	} else {
		return nil, fmt.Errorf("missing GitHub App private key information")
	}
}
