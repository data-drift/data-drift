package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/data-drift/kpi-git-history/charts"
	"github.com/data-drift/kpi-git-history/common"
	"github.com/data-drift/kpi-git-history/database/notion_database"
	"github.com/data-drift/kpi-git-history/github"
	"github.com/data-drift/kpi-git-history/history"
	"github.com/data-drift/kpi-git-history/reports"
)

func DebugFunction() {
	// Perform debugging operations
	fmt.Println("Manual Sync ...")
	// githubToken := os.Getenv("GITHUB_TOKEN")
	githubRepoOwner := os.Getenv("GITHUB_REPO_OWNER")
	githubRepoName := os.Getenv("GITHUB_REPO_NAME")
	githubRepoFilePath := os.Getenv("GITHUB_REPO_FILE_PATH")
	githubApplicationIdStr := os.Getenv("GITHUB_APPLICATION_INSTALLATION_ID")
	dateColumn := os.Getenv("DATE_COLUMN")
	kpiColumn := os.Getenv("KPI_COLUMN")
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	notionDatabaseID := os.Getenv("NOTION_DATABASE_ID")

	filepath := os.Getenv("DEFAULT_FILE_PATH")
	githubApplicationId, _ := strconv.ParseInt(githubApplicationIdStr, 10, 64)

	_ = notion_database.AssertDatabaseHasDatadriftProperties(notionDatabaseID, notionAPIKey)

	client, _ := github.CreateClientFromGithubApp(int64(githubApplicationId))

	if filepath == "" {
		newFilepath, err := history.ProcessHistory(client, githubRepoOwner, githubRepoName, common.Metric{
			MetricName:     "Default metric name",
			KPIColumnName:  kpiColumn,
			DateColumnName: dateColumn,
			Filepath:       githubRepoFilePath,
			TimeGrains:     []common.TimeGrain{common.Quarter, common.Year},
			Dimensions:     []string{},
		})

		if err != nil {
			println(err)
		}
		filepath = newFilepath
	}

	chartResults := charts.ProcessCharts(filepath, common.Metric{MetricName: "Default metric name"})

	// if (len(chartResults)) != 0 {
	// 	println("Stop exectution here")
	// 	return
	// }

	for _, chartResult := range chartResults {
		err := reports.CreateReport(common.SyncConfig{NotionAPIKey: notionAPIKey, NotionDatabaseID: notionDatabaseID}, chartResult)
		if err != nil {
			println(err)
		}
	}
	println(filepath)
}
