package history

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/data-drift/kpi-git-history/common"
	"github.com/google/go-github/github"
)

func ProcessHistory(client *github.Client, syncConfig common.SyncConfig) (error, string) {

	// Create a new GitHub client without authentication.

	// Set the repository owner and name.
	owner := syncConfig.GithubRepoOwner
	name := syncConfig.GithubRepoName

	// Set the path to the CSV file in the repository.
	path := syncConfig.GithubRepoFilePath

	// Set the start and end dates to display the history for.
	startDateStr := syncConfig.StartDate
	endDate := time.Now()
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return fmt.Errorf("error parsing start date: %v", err), ""
	}

	dateColumnName := syncConfig.DateColumn
	if dateColumnName == "" {
		return fmt.Errorf("error no date column name provided"), ""
	}
	KPIColumnName := syncConfig.KpiColumn

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
		return fmt.Errorf("error getting commit history: %v", err), ""
	}

	// Print the number of commits.
	fmt.Printf("Number of commits: %d\n", len(commits))

	// Group the lines of the CSV file by reporting date.
	lineCountAndKPIByDateByVersion := make(map[string]map[string]struct {
		Lines           int
		KPI             float64
		CommitTimestamp int64
	})
	for index, commit := range commits {
		fmt.Printf("\r Commit %d/%d", index, len(commits))

		commitDate := commit.Commit.Author.Date
		commitTimestamp := commitDate.Unix()
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
					Lines           int
					KPI             float64
					CommitTimestamp int64
				})
			}

			kpiStr := record[kpiColumn]
			kpi, _ := strconv.ParseFloat(kpiStr, 64)

			newLineCount := lineCountAndKPIByDateByVersion[dateStr][*commit.SHA].Lines + 1
			newKPI := lineCountAndKPIByDateByVersion[dateStr][*commit.SHA].KPI + kpi

			lineCountAndKPIByDateByVersion[dateStr][*commit.SHA] = struct {
				Lines           int
				KPI             float64
				CommitTimestamp int64
			}{Lines: newLineCount, KPI: newKPI, CommitTimestamp: commitTimestamp}
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

	if _, err := os.Stat("dist"); os.IsNotExist(err) {
		if err := os.Mkdir("dist", 0755); err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}

	// Create a slice to hold the date strings.
	dates := make([]string, 0, len(lineCountAndKPIByDateByVersion))

	// Add the date strings to the slice.
	for dateStr := range lineCountAndKPIByDateByVersion {
		dates = append(dates, dateStr)
	}

	// Sort the date strings in ascending order.
	sort.Strings(dates)

	// Create a map to hold the line counts by date by version, ordered by date.
	orderedLineCountsByDateByVersion := make(map[string]map[string]struct {
		Lines           int
		KPI             float64
		CommitTimestamp int64
	})

	// Copy the line counts by date by version to the new map, ordered by date.
	for _, dateStr := range dates {
		orderedLineCountsByDateByVersion[dateStr] = lineCountAndKPIByDateByVersion[dateStr]
	}
	// Generate a timestamp to include in the JSON file name.
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	// Open a file to write the line counts by date by version in JSON format.
	filepath := fmt.Sprintf("dist/lineCountAndKPIByDateByVersion_%s.json", timestamp)
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// Write the line counts and KPI values to the JSON file.
	enc := json.NewEncoder(file)
	if err := enc.Encode(orderedLineCountsByDateByVersion); err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}
	fmt.Println("Results written to lineCountsAndKPIs.json")
	return nil, filepath
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
