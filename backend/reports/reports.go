package reports

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/database/notion_database"
	"github.com/dstotijn/go-notion"
	"github.com/gin-gonic/gin"
)

func CreateReport(syncConfig common.SyncConfig, KPIInfo common.KPIReport) error {
	timeGrain, _ := GetTimeGrain(KPIInfo.PeriodId)
	reportNotionPageId, shouldInitReport, findOrCreateError := notion_database.FindOrCreateReportPageId(syncConfig.NotionAPIKey, syncConfig.NotionDatabaseID, KPIInfo.KPIName, string(KPIInfo.PeriodId), timeGrain, KPIInfo.DimensionValue)
	if findOrCreateError != nil {
		return fmt.Errorf("failed to create reportNotionPageId: %v", findOrCreateError.Error())
	}

	if shouldInitReport {
		err := notion_database.InitChangeLogReport(syncConfig.NotionAPIKey, reportNotionPageId, KPIInfo)
		if err != nil {
			return fmt.Errorf("failed to create page: %v", err.Error())
		}
	} else {
		err := notion_database.UpdateChangeLogReport(syncConfig.NotionAPIKey, reportNotionPageId, KPIInfo)
		if err != nil {
			print("Updating report error", err.Error())
		}
	}

	return nil
}

func CreateSummaryReport(syncConfig common.SyncConfig, metricConfig common.MetricConfig, chartUrls map[common.TimeGrain]string) error {
	fmt.Println("Creating summary report")
	reportNotionPageId, findOrCreateError := notion_database.FindOrCreateSummaryReportPage(syncConfig.NotionAPIKey, syncConfig.NotionDatabaseID, "Summary of "+metricConfig.MetricName)
	fmt.Println(reportNotionPageId)
	fmt.Println(findOrCreateError)

	var children []notion.Block
	for _, timeGrain := range []common.TimeGrain{common.Day, common.Week, common.Month, common.Quarter, common.Year} {
		chartUrl := chartUrls[timeGrain]
		if chartUrl == "" {
			continue
		}
		children = append(children, notion.Heading1Block{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "Cohorts " + string(timeGrain),
					},
				},
			},
		},
			notion.EmbedBlock{
				URL: chartUrl,
			},
		)
		reports, err := notion_database.QueryDatabaseWithMetricAndTimegrain(syncConfig.NotionAPIKey, syncConfig.NotionDatabaseID, metricConfig.MetricName, timeGrain)
		for _, report := range reports {
			children = append(children, notion.LinkToPageBlock{
				Type:   "page_id",
				PageID: report.ID,
			})
		}
		if err != nil {
			print("Error getting links for report in summary page", err.Error())
		}
		// get report of metric and timegrain, ordered by name and push them in children
	}
	err := notion_database.UpdateMetadataReport(syncConfig.NotionAPIKey, reportNotionPageId, children, &notion.DatabasePageProperties{})
	if err != nil {
		return fmt.Errorf("failed to create page: %v", err.Error())
	}

	return nil
}

func GetTimeGrain(periodKeyParam common.PeriodKey) (common.TimeGrain, error) {
	periodKey := string(periodKeyParam)
	_, err := time.Parse("2006-01-02", periodKey)
	if err == nil {
		return common.Day, nil
	}
	_, err = ParseYearWeek(periodKey)
	if err == nil {
		return common.Week, nil
	}
	_, err = time.Parse("2006-01", periodKey)
	if err == nil {
		return common.Month, nil
	}
	_, err = ParseQuarterDate(periodKey)
	if err == nil {
		return common.Quarter, nil
	}
	_, err = time.Parse("2006", periodKey)
	if err == nil {
		return common.Year, nil
	}
	return "", fmt.Errorf("invalid period key: %s", periodKey)
}

func ParseYearWeek(yearWeek string) (time.Time, error) {
	if len(yearWeek) != 8 {
		return time.Time{}, fmt.Errorf("invalid year week format: %s", yearWeek)
	}
	year, err := strconv.Atoi(yearWeek[0:4])
	if err != nil {
		return time.Time{}, err
	}

	week, err := strconv.Atoi(yearWeek[6:])
	if err != nil {
		return time.Time{}, err
	}

	// Get the first day of the week (Monday)
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 7*(week-1)+1)

	return firstDay, nil
}

func ParseQuarterDate(s string) (time.Time, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid quarter date format: %s", s)
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year format in quarter date: %s", s)
	}
	quarter := parts[1]
	switch quarter {
	case "Q1":
		return time.Date(year, time.March, 31, 0, 0, 0, 0, time.UTC), nil
	case "Q2":
		return time.Date(year, time.June, 30, 0, 0, 0, 0, time.UTC), nil
	case "Q3":
		return time.Date(year, time.September, 30, 0, 0, 0, 0, time.UTC), nil
	case "Q4":
		return time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC), nil
	default:
		return time.Time{}, fmt.Errorf("invalid quarter format in quarter date: %s", s)
	}
}

func GetMetricCohort(c *gin.Context) {

	InstallationId := c.Request.Header.Get("Installation-Id")

	if InstallationId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No installation id provided"})
		return
	}

	metricName := c.Param("metric-name")
	timeGrain := c.Param("timegrain")

	filepath, err := common.GetLatestMetricFile(InstallationId, metricName)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metricHistory, err := common.GetKeysFromJSON(filepath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := GetReportData(metricHistory, common.TimeGrain(timeGrain))

	c.JSON(http.StatusOK, response)
}

func GetReportData(metrics common.Metrics, timeGrain common.TimeGrain) map[string]interface{} {
	// Calculate the cohort dates based on the time grain.
	cohortDates := []string{}

	reportData := make(map[int64]map[string]interface{})

	for cohortName, cohort := range metrics {
		if cohort.TimeGrain == timeGrain && cohort.Dimension == "none" {
			print(string(cohort.Dimension))
			cohortDates = append(cohortDates, string(cohortName))
			for _, commit := range cohort.History {
				timestampStr := commit.CommitTimestamp
				if _, ok := reportData[timestampStr]; !ok {
					reportData[timestampStr] = make(map[string]interface{})
				}
				reportData[timestampStr][string(cohortName)] = commit.KPI
			}

		}

	}

	// Calculate the data based on the metrics.
	// for _, cohort := range metrics {
	// 	versionData := map[string]interface{}{
	// 		"version": version,
	// 	}
	// 	for i := 0; i < metrics.CohortSize; i++ {
	// 		date := now.Add(-time.Duration(i) * timeGrain.Duration()).Format("2006-01-02")
	// 		count := rand.Intn(100)
	// 		versionData[date] = count
	// 	}
	// 	data = append(data, versionData)
	// }

	// Return the response map.
	return map[string]interface{}{
		"cohortDates": cohortDates,
		"data":        reportData,
	}
}
