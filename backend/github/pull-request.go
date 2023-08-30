package github

import (
	"context"
	"fmt"

	"github.com/data-drift/data-drift/urlgen"
	"github.com/google/go-github/v42/github"
)

func handlePullRequestOpened(event *github.PullRequestEvent) error {
	if event.GetAction() != "opened" {
		return nil
	}
	print("Pull request opened", event)
	// Get the owner and repository name from the event.
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	installationId := event.GetInstallation().GetID()

	// Get the pull request number from the event.
	number := event.GetNumber()

	commitDiffUrl := urlgen.BuildReportDiffUrl(urlgen.BuildReportDiffBaseUrl(fmt.Sprint(installationId), owner, repo), "036f9d6b685ee02a14faa70ed05e0bd60650c477")

	// Create a comment on the pull request.
	comment := &github.IssueComment{
		Body: github.String(fmt.Sprintf("The diff is available [here](%s).", commitDiffUrl)),
	}
	client, _ := CreateClientFromGithubApp(*event.Installation.ID)

	_, _, err := client.Issues.CreateComment(context.Background(), owner, repo, number, comment)

	if err != nil {
		print(err.Error())
	}
	return err
}
