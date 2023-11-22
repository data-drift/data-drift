package github

import (
	"context"
	"fmt"
	"log"

	"github.com/data-drift/data-drift/urlgen"
	"github.com/google/go-github/v56/github"
)

func handleIssueOpened(event *github.IssuesEvent) error {
	if event.GetAction() != "opened" {
		return nil
	}
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	installationId := event.GetInstallation().GetID()
	log.Printf("Issue opened: owner=%s, repo=%s, installation_id=%d, number=%d, title=%s, url=%s", event.Repo.Owner.GetLogin(), event.Repo.GetName(), event.Installation.GetID(), *event.Issue.Number, event.Issue.GetTitle(), event.Issue.GetHTMLURL())

	commitDiffUrl := urlgen.BuildOverviewUrl(fmt.Sprint(installationId), owner, repo)
	log.Printf("commitDiffUrl: %s", commitDiffUrl)
	comment := &github.IssueComment{
		Body: github.String(fmt.Sprintf("The diff is available [here](%s).", commitDiffUrl)),
	}
	client, _ := CreateClientFromGithubApp(*event.Installation.ID)
	number := event.Issue.Number

	_, _, err := client.Issues.CreateComment(context.Background(), owner, repo, *number, comment)

	if err != nil {
		print(err.Error())
	}
	return err
}
