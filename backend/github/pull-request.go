package github

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/data-drift/data-drift/urlgen"
	"github.com/google/go-github/v56/github"
)

func handlePullRequestOpened(event *github.PullRequestEvent) error {
	if event.GetAction() != "opened" {
		return nil
	}
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	log.Printf("Pull request opened: owner=%s, repo=%s, installation_id=%d, number=%d, title=%s, url=%s", event.Repo.Owner.GetLogin(), event.Repo.GetName(), event.Installation.GetID(), event.PullRequest.GetNumber(), event.PullRequest.GetTitle(), event.PullRequest.GetHTMLURL())

	commitDiffUrl := urlgen.BuildReportDiffUrl(urlgen.BuildReportDiffBaseUrl(owner, repo), *event.PullRequest.Head.SHA, url.Values{})
	log.Printf("commitDiffUrl: %s", commitDiffUrl)
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
