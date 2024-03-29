package metrics

import (
	"fmt"
	"net/http"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/github"
	"github.com/data-drift/data-drift/reducers"
	"github.com/gin-gonic/gin"
)

type MetricService struct {
	KpiRepository *common.KpiRepository
}

func NewMetricService(kpiRepository *common.KpiRepository) *MetricService {
	return &MetricService{KpiRepository: kpiRepository}
}

func (h *MetricService) GetMetricCohort(c *gin.Context) {

	InstallationId := c.Request.Header.Get("Installation-Id")
	metricName := c.Param("metric-name")
	timeGrain := c.Param("timegrain")
	var filepath common.MetricStorageKey
	if InstallationId == "" {
		githubConnectionValue, exists := c.Get("github_connection")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub client not found"})
			return
		}
		githubConnection, ok := githubConnectionValue.(github.GithubConnection)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid GitHub client"})
			return
		}
		filepath = common.NewGetMetricStorageKey(githubConnection.Owner, githubConnection.Repository, metricName)
	} else {
		filepath = common.LegacyGetMetricStorageKey(InstallationId, metricName)
	}

	metricHistory, err := h.KpiRepository.ReadMetricKPI(filepath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := GetReportData(metricHistory, common.TimeGrain(timeGrain))

	c.JSON(http.StatusOK, response)
}

func (h *MetricService) GetMetricReport(c *gin.Context) {

	InstallationId := c.Request.Header.Get("Installation-Id")
	metricName := c.Param("metric-name")
	var filepath common.MetricStorageKey

	if InstallationId == "" {
		githubConnectionValue, exists := c.Get("github_connection")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub client not found"})
			return
		}
		githubConnection, ok := githubConnectionValue.(github.GithubConnection)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid GitHub client"})
			return
		}
		filepath = common.NewGetMetricStorageKey(githubConnection.Owner, githubConnection.Repository, metricName)
	} else {
		filepath = common.LegacyGetMetricStorageKey(InstallationId, metricName)
	}

	metricHistory, err := h.KpiRepository.ReadMetricKPI(filepath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metricHistory)
}

func GetReportData(metrics common.Metrics, timeGrain common.TimeGrain) map[string]interface{} {
	cohortDates := []string{}
	reportData := make(map[int64]map[string]interface{})
	cohortsMetricsMetadata := make(map[string]reducers.MetricMetadata)

	for cohortName, cohort := range metrics {
		if cohort.TimeGrain == timeGrain && cohort.Dimension == "none" {
			cohortDates = append(cohortDates, string(cohortName))
			metricMetadata, err := reducers.GetMetadataOfMetric(cohort)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			cohortsMetricsMetadata[string(cohortName)] = metricMetadata
			for _, commit := range metricMetadata.RelativeHistory {
				timestampStr := commit.ComputationTimetamp
				if _, ok := reportData[timestampStr]; !ok {
					reportData[timestampStr] = make(map[string]interface{})
				}
				reportData[timestampStr][string(cohortName)] = commit.RelativeValue
			}

		}

	}

	// Return the response map.
	return map[string]interface{}{
		"timegrain":              timeGrain,
		"cohortDates":            cohortDates,
		"dataIndexedByTimestamp": reportData,
		"cohortsMetricsMetadata": cohortsMetricsMetadata,
	}
}
