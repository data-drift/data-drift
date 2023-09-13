package urlgen

import (
	"fmt"

	"github.com/data-drift/data-drift/common"
)

func MetricCohortUrl(installationId string, metricName string, timegrain common.TimeGrain) string {
	return fmt.Sprintf("https://app.data-drift.io/report/%s/metrics/%s/cohorts/%s", installationId, metricName, timegrain)
}

func MetricReportUrl(installationId common.GithubInstallationId, metricName string, period common.PeriodKey) string {
	return fmt.Sprintf("https://app.data-drift.io/report/%s/metrics/%s/report/%s", string(installationId), metricName, string(period))
}

func BuildReportDiffBaseUrl(installationId, repoOwner, repoName string) string {
	reportBaseUrl := fmt.Sprintf("https://app.data-drift.io/report/%s/%s/%s/commit", installationId, repoOwner, repoName)
	return reportBaseUrl
}

func BuildReportDiffUrl(reportBaseUrl, commitSha, queryString string) string {
	url := fmt.Sprintf("%s/%s?%s", reportBaseUrl, commitSha, queryString)
	return url
}
