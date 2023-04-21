package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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
	KPIColumnName := os.Getenv("KPI_COLUMN")

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
	lineCountAndKPIByDateByVersion := make(map[string]map[string]struct {
		Lines int
		KPI   float64
	})
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
		var kpiColumn int

		for i, columnName := range records[0] {
			if columnName == dateColumnName {
				dateColumn = i
			}
			if columnName == KPIColumnName {
				kpiColumn = i
			}
		}
		for _, record := range records[1:] { // Skip the header row.
			dateStr := record[dateColumn]

			if lineCountAndKPIByDateByVersion[dateStr] == nil {
				lineCountAndKPIByDateByVersion[dateStr] = make(map[string]struct {
					Lines int
					KPI   float64
				})
			}

			kpiStr := record[kpiColumn]
			kpi, _ := strconv.ParseFloat(kpiStr, 64)

			newLineCount := lineCountAndKPIByDateByVersion[dateStr][*commit.SHA].Lines + 1
			newKPI := lineCountAndKPIByDateByVersion[dateStr][*commit.SHA].KPI + kpi

			lineCountAndKPIByDateByVersion[dateStr][*commit.SHA] = struct {
				Lines int
				KPI   float64
			}{Lines: newLineCount + 1, KPI: newKPI}
		}

	}

	// Print the line count for each reporting date.
	for dateStr, lineCounts := range lineCountAndKPIByDateByVersion {

		var countsStr string
		var kpiStr string
		for _, count := range lineCounts {
			countsStr += fmt.Sprintf("%d ", count.Lines)
			kpiStr += fmt.Sprintf("%.2f ", count.KPI)
		}
		fmt.Printf("Line Count %s: %s\n", dateStr, countsStr)
		fmt.Printf("       KPI %s: %s\n", dateStr, kpiStr)
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
