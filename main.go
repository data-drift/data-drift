package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("GITHUB_TOKEN")

	// Create a new GitHub client with authentication using the token.
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Create a new GitHub client without authentication.

	// Set the repository owner and name.
	owner := "data-drift"
	name := "fake-dataset"

	// Set the path to the CSV file in the repository.
	path := "mart/organisation-mrr.csv"

	// Set the start and end dates to display the history for.
	startDateStr := "2022-03-20"
	endDate := time.Now()
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Fatalf("Error parsing start date: %v", err)
	}

	// Get the commit history for the repository.
	// Get the commit history for the file.
	commits, _, err := client.Repositories.ListCommits(context.Background(), owner, name, &github.CommitsListOptions{
		Path:        path,
		SHA:         "",
		Since:       startDate,
		Until:       endDate,
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		log.Fatalf("Error getting commit history: %v", err)
	}

	// Print the number of commits.
	fmt.Printf("Number of commits: %d\n", len(commits))

	// Iterate through the commits and display the number of lines in the file for each commit between the start and end dates.
	var numLines int
	for _, commit := range commits {
		if commit.Commit.Author.Date.After(endDate) || commit.Commit.Author.Date.Before(startDate) {
			continue
		}

		// Get the file contents for the commit.
		fileContents, err := getFileContentsForCommit(client, owner, name, path, *commit.SHA)
		if err != nil {
			log.Printf("Error getting file contents for commit %s: %v", *commit.SHA, err)
			continue
		}

		// Parse the CSV file and count the number of lines.
		r := csv.NewReader(bytes.NewReader(fileContents))
		lines, err := r.ReadAll()
		if err != nil {
			log.Printf("Error parsing CSV file for commit %s: %v", *commit.SHA, err)
			continue
		}
		numLines = len(lines)

		// Display the number of lines and the commit message for the commit.
		fmt.Printf("%d lines - %s\n", numLines, *commit.Commit.Message)
	}
}

func getFileContentsForCommit(client *github.Client, owner, name, path, sha string) ([]byte, error) {
	opts := &github.RepositoryContentGetOptions{Ref: sha}
	fileContents, _, resp, err := client.Repositories.GetContents(context.Background(), owner, name, path, opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := fileContents.GetContent()
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}
