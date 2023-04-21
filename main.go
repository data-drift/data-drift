package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
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

	// Create a new GitHub client without authentication.

	// Set the repository owner and name.
	owner := os.Getenv("GITHUB_REPO_OWNER")
	name := os.Getenv("GITHUB_REPO_NAME")

	// Set the path to the CSV file in the repository.
	path := os.Getenv("GITHUB_REPO_FILE_PATH")

	// Set the start and end dates to display the history for.
	startDateStr := os.Getenv("START_DATE")
	endDate := time.Now()
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Fatalf("Error parsing start date: %v", err)
	}

	dateColumnName := os.Getenv("DATE_COLUMN")
	if dateColumnName == "" {
		log.Fatalf("Error no date column name provided")
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

	// Group the lines of the CSV file by reporting date.
	lineCountByDateByVersion := make(map[string]map[string]int)

	for _, commit := range commits {
		fileContents, err := getFileContentsForCommit(client, owner, name, path, *commit.SHA)
		if err != nil {
			log.Printf("Error getting file contents for commit %s: %v", *commit.SHA, err)
			continue
		}
		r := csv.NewReader(bytes.NewReader(fileContents))
		records, err := r.ReadAll()
		if err != nil {
			log.Printf("Error parsing CSV file for commit %s: %v", *commit.SHA, err)
			continue
		}
		var dateColumn int

		for i, columnName := range records[0] {
			if columnName == dateColumnName {
				dateColumn = i
				break
			}
		}
		for _, record := range records[1:] { // Skip the header row.
			dateStr := record[dateColumn]

			if lineCountByDateByVersion[dateStr] == nil {
				lineCountByDateByVersion[dateStr] = make(map[string]int)
			}
			lineCountByDateByVersion[dateStr][*commit.SHA]++
		}

	}

	// Print the line count for each reporting date.
	for dateStr, lineCounts := range lineCountByDateByVersion {
		// for commitSha, count := range line {
		// 	fmt.Printf("%s %s: %d\n", dateStr, commitSha, count)
		// }
		var countsStr string
		for _, count := range lineCounts {
			countsStr += fmt.Sprintf("%d ", count)
		}
		fmt.Printf("%s: %s\n", dateStr, countsStr)
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
