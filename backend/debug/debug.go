package debug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/database/notion_database"
	"github.com/data-drift/data-drift/github"
	"github.com/data-drift/data-drift/history"
	"github.com/data-drift/data-drift/reducers"
	"github.com/data-drift/data-drift/reports"
)

func DebugFunction() {
	// Perform debugging operations
	fmt.Println("Manual Sync ...")
	// githubToken := os.Getenv("GITHUB_TOKEN")
	githubRepoOwner := os.Getenv("GITHUB_REPO_OWNER")
	githubRepoName := os.Getenv("GITHUB_REPO_NAME")
	githubRepoFilePath := os.Getenv("GITHUB_REPO_FILE_PATH")
	githubApplicationIdStr := os.Getenv("GITHUB_APPLICATION_INSTALLATION_ID")
	metricName := os.Getenv("METRIC_NAME")
	dateColumn := os.Getenv("DATE_COLUMN")
	kpiColumn := os.Getenv("KPI_COLUMN")
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	notionDatabaseID := os.Getenv("NOTION_DATABASE_ID")

	filepath := common.MetricStorageKey(os.Getenv("DEFAULT_FILE_PATH"))
	githubApplicationId, _ := strconv.ParseInt(githubApplicationIdStr, 10, 64)

	_ = notion_database.AssertDatabaseHasDatadriftProperties(notionDatabaseID, notionAPIKey)

	metricConfig := common.MetricConfig{
		MetricName:     metricName,
		KPIColumnName:  kpiColumn,
		DateColumnName: dateColumn,
		Filepath:       githubRepoFilePath,
		TimeGrains:     []common.TimeGrain{common.Month},
		Dimensions:     []string{},
	}
	if filepath == "" {
		client, _ := github.CreateClientFromGithubApp(int64(githubApplicationId))
		if client == nil {
			panic("Client not configured")
		}
		newFilepath, err := history.ProcessHistory(client, githubRepoOwner, githubRepoName, metricConfig, int(githubApplicationId))

		if err != nil {
			println(err)
		}
		filepath = newFilepath
	}
	notionSyncConfig := common.SyncConfig{NotionAPIKey: notionAPIKey, NotionDatabaseID: notionDatabaseID}

	// metadataChartResults, metadataChartError := reducers.ProcessMetricMetadataCharts(filepath, metricConfig)
	// if metadataChartError != nil {
	// 	println(metadataChartError)
	// 	panic("Error processing metadata charts")
	// }

	// reports.CreateSummaryReport(notionSyncConfig, metricConfig, metadataChartResults)
	// if metadataChartResults != nil {
	// 	panic("Stop execution here")
	// }

	chartResults := reducers.ProcessMetricHistory(filepath, common.MetricConfig{MetricName: "Default metric name"})

	for _, chartResult := range chartResults {
		err := reports.CreateReport(notionSyncConfig, chartResult)
		if err != nil {
			println(err)
		}
	}

	metadataChartResults, metadataChartError := reducers.ProcessMetricMetadataCharts(filepath, metricConfig)
	if metadataChartError != nil {
		fmt.Println("[DATADRIFT_ERROR] create summary report", metadataChartError.Error())
	} else {
		reports.CreateSummaryReport(notionSyncConfig, metricConfig, metadataChartResults, fmt.Sprint(githubApplicationId))
	}
	println(filepath)
}
