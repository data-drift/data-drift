package urlgen

import (
	"fmt"
	"net/url"

	"github.com/data-drift/data-drift/common"
)

func MetricCohortUrl(owner string, repo string, metricName string, timegrain common.TimeGrain) string {
	return fmt.Sprintf("https://app.data-drift.io/report/%s/%s/metrics/%s/cohorts/%s", owner, repo, metricName, timegrain)
}

func MetricReportUrl(owner string, repo string, metricName string, period common.PeriodKey, dimensionValue string) string {
	url := fmt.Sprintf("https://app.data-drift.io/report/%s/%s/metrics/%s/report/%s", owner, repo, metricName, string(period))
	if dimensionValue != "" {
		url += fmt.Sprintf("?dimensionValue=%s", dimensionValue)
	}
	return url
}

func BuildReportDiffBaseUrl(repoOwner, repoName string) string {
	reportBaseUrl := fmt.Sprintf("https://app.data-drift.io/report/%s/%s/commit", repoOwner, repoName)
	return reportBaseUrl
}

func BuildOverviewUrl(repoOwner, repoName, snapshotDate, tableName string) string {
	queryString := url.Values{
		"snapshotDate": {snapshotDate},
		"tableName":    {tableName},
	}
	url := fmt.Sprintf("https://app.data-drift.io/%s/%s/overview?%s", repoOwner, repoName, queryString.Encode())
	return url
}

func BuildReportDiffUrl(reportBaseUrl, commitSha string, queryString url.Values) string {
	url := fmt.Sprintf("%s/%s?%s", reportBaseUrl, commitSha, queryString.Encode())
	return url
}
