package metrics

import (
	"net/http"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/reducers"
	"github.com/gin-gonic/gin"
)

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
	cohortDates := []string{}
	reportData := make(map[int64]map[string]interface{})
	cohortsMetricsMetadata := make(map[string]reducers.MetricMetadata)

	for cohortName, cohort := range metrics {
		if cohort.TimeGrain == timeGrain && cohort.Dimension == "none" {
			cohortDates = append(cohortDates, string(cohortName))
			metricMetadata, err := reducers.GetMetadataOfMetric(cohort)
			if err != nil {
				print(err.Error())
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
