package github

import (
	"context"

	"github.com/google/go-github/v42/github"
)

func GetCommitComments(client *github.Client, ctx context.Context, RepoOwner string, RepoName string, commitSha string) []*github.IssueComment {
	pullRequests, _, err := client.PullRequests.ListPullRequestsWithCommit(ctx, RepoOwner, RepoName, commitSha, nil)
	if err != nil {
		println("[DATADRIFT ERROR] getting pull request:", err.Error())
		return []*github.IssueComment{}
	}
	if len(pullRequests) == 0 {
		println("No pull request found")
		return []*github.IssueComment{}
	}
	comments, _, err := client.Issues.ListComments(ctx, RepoOwner, RepoName, pullRequests[0].GetNumber(), nil)
	if err != nil {
		println("[DATADRIFT ERROR] getting comments:", err.Error())
		return []*github.IssueComment{}
	}

	return comments
}
