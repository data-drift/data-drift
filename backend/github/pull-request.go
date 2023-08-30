package github

import (
	"context"

	"github.com/google/go-github/v42/github"
)

func handlePullRequestOpened(event *github.PullRequestEvent) error {
	print("Pull request opened", event)
	// Get the owner and repository name from the event.
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()

	// Get the pull request number from the event.
	number := event.GetNumber()

	// Create a comment on the pull request.
	comment := &github.IssueComment{
		Body: github.String("Thanks for opening this pull request!"),
	}
	client, _ := CreateClientFromGithubApp(*event.Installation.ID)

	_, _, err := client.Issues.CreateComment(context.Background(), owner, repo, number, comment)

	if err != nil {
		print(err.Error())
	}
	return err
}
