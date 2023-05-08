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

type CommitSha string
type PeriodId string

type PeriodCommitData struct {
	Lines           int
	KPI             float64
	CommitTimestamp int64
	CommitUrl       string
}

type PeriodData map[PeriodId]map[CommitSha]PeriodCommitData

func ProcessHistory(client *github.Client, repoOwner string, repoName string, metric common.Metric) (string, error) {

	filePath := metric.Filepath
	dateColumnName := metric.DateColumnName
	KPIColumnName := metric.KPIColumnName
	metricName := metric.MetricName

	fmt.Println(metric)
	// Set the start and end dates to display the history for.
	endDate := time.Now()

	if dateColumnName == "" {
		return "", fmt.Errorf("error no date column name provided")
	}

	// Get the commit history for the repository.
	// Get the commit history for the file.
	commits, _, err := client.Repositories.ListCommits(context.Background(), repoOwner, repoName, &github.CommitsListOptions{
		Path:        filePath,
		SHA:         "",
		Until:       endDate,
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return "", fmt.Errorf("error getting commit history: %v", err)
	}

	// Print the number of commits.
	fmt.Printf("Number of commits: %d\n", len(commits))

	// Group the lines of the CSV file by reporting date.
	lineCountAndKPIByDateByVersion := make(PeriodData)
	for index, commit := range commits {
		fmt.Printf("\r Commit %d/%d", index, len(commits))

		commitSha := CommitSha(*commit.SHA)

		commitDate := commit.Commit.Author.Date
		commitTimestamp := commitDate.Unix()
		fileContents, err := getFileContentsForCommit(client, repoOwner, repoName, filePath, *commit.SHA)
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
			for _, timegrain := range metric.TimeGrains {
				var periodKey PeriodId
				periodTime, _ := time.Parse("2006-01-02", record[dateColumn])

				switch timegrain {
				case common.Day:
					periodKey = PeriodId(periodTime.Format("2006-01-02"))
				case common.Week:
					_, week := periodTime.ISOWeek()
					periodKey = PeriodId(fmt.Sprintf("%d-%d", periodTime.Year(), week))
				case common.Month:
					periodKey = PeriodId(periodTime.Format("2006-01"))
				case common.Year:
					periodKey = PeriodId(periodTime.Format("2006"))
				default:
					log.Fatalf("Invalid time grain: %s", timegrain)
				}

				if lineCountAndKPIByDateByVersion[periodKey] == nil {
					lineCountAndKPIByDateByVersion[periodKey] = make(map[CommitSha]PeriodCommitData)
				}

				kpiStr := record[kpiColumn]
				kpi, _ := strconv.ParseFloat(kpiStr, 64)

				newLineCount := lineCountAndKPIByDateByVersion[periodKey][commitSha].Lines + 1
				newKPI := lineCountAndKPIByDateByVersion[periodKey][commitSha].KPI + kpi

				lineCountAndKPIByDateByVersion[periodKey][commitSha] = struct {
					Lines           int
					KPI             float64
					CommitTimestamp int64
					CommitUrl       string
				}{Lines: newLineCount, KPI: newKPI, CommitTimestamp: commitTimestamp, CommitUrl: *commit.HTMLURL}
			}
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
		dates = append(dates, string(dateStr))
	}

	// Sort the date strings in ascending order.
	sort.Strings(dates)

	// Create a map to hold the line counts by date by version, ordered by date.
	orderedLineCountsByDateByVersion := make(map[string]map[CommitSha]PeriodCommitData)

	// Copy the line counts by date by version to the new map, ordered by date.
	for _, dateStr := range dates {
		orderedLineCountsByDateByVersion[dateStr] = lineCountAndKPIByDateByVersion[PeriodId(dateStr)]
	}
	// Generate a timestamp to include in the JSON file name.
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	// Open a file to write the line counts by date by version in JSON format.
	filepath := fmt.Sprintf("dist/"+metricName+"lineCountAndKPIByDateByVersion_%s.json", timestamp)
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
	return filepath, nil
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
