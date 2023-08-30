package github

import (
	"context"
	"fmt"
	"log"

	"github.com/data-drift/data-drift/urlgen"
	"github.com/google/go-github/v42/github"
)

func handlePullRequestOpened(event *github.PullRequestEvent) error {
	if event.GetAction() != "opened" {
		return nil
	}
	// Get the owner and repository name from the event.
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	installationId := event.GetInstallation().GetID()
	log.Printf("Pull request opened: owner=%s, repo=%s, installation_id=%d, number=%d, title=%s, url=%s", event.Repo.Owner.GetLogin(), event.Repo.GetName(), event.Installation.GetID(), event.PullRequest.GetNumber(), event.PullRequest.GetTitle(), event.PullRequest.GetHTMLURL())
	// Get the pull request number from the event.

	commitDiffUrl := urlgen.BuildReportDiffUrl(urlgen.BuildReportDiffBaseUrl(fmt.Sprint(installationId), owner, repo), *event.PullRequest.Head.SHA)

	// Create a comment on the pull request.
	comment := &github.IssueComment{
		Body: github.String(fmt.Sprintf("The diff is available [here](%s).", commitDiffUrl)),
	}
	client, _ := CreateClientFromGithubApp(*event.Installation.ID)
	number := event.GetNumber()

	_, _, err := client.Issues.CreateComment(context.Background(), owner, repo, number, comment)

	if err != nil {
		print(err.Error())
	}
	return err
}
