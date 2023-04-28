package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func CreateClientFromGithubApp(installationId int64) (*github.Client, error) {
	privateKeyPath := os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH")
	privateKey := os.Getenv("GITHUB_APP_PRIVATE_KEY")

	appIDStr := os.Getenv("GITHUB_APP_ID")
	githubAppId, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	itr, err := CreateGithubTransport(privateKeyPath, privateKey, githubAppId, installationId)
	if err != nil {
		return nil, err
	}

	// Use installation transport with client.
	client := github.NewClient(&http.Client{Transport: itr})
	return client, nil
}

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
