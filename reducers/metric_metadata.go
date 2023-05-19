package reducers

import (
	"time"

	"github.com/data-drift/kpi-git-history/common"
	"github.com/shopspring/decimal"
)

type RelativeHistoricalEvent struct {
	RelativeValue         decimal.Decimal
	DaysFromHistorization decimal.Decimal
}

type MetricMetadata struct {
	TimeGrain       common.TimeGrain
	PeriodKey       common.PeriodKey
	InitialValue    decimal.Decimal
	FirstDate       time.Time
	RelativeHistory []RelativeHistoricalEvent
}

func ProcessMetricMetadata(metricConfig common.MetricConfig, metrics common.Metrics) map[common.TimeGrain]MetricMetadata {

	metric := metrics[common.PeriodAndDimensionKey("2023-02 FR")]
	metricMetadata := getMetadataOfMetric(metric)
	return map[common.TimeGrain]MetricMetadata{
		common.Month: metricMetadata,
	}
}

func getMetadataOfMetric(metric common.Metric) MetricMetadata {
	firstDateOfPeriod := getFirstDateOfPeriod(metric.Period)
	var dataSortableArray []common.CommitData

	for _, stats := range metric.History {
		KPI := stats.KPI
		dataSortableArray = append(dataSortableArray, common.CommitData{
			Lines:           stats.Lines,
			KPI:             KPI,
			CommitTimestamp: stats.CommitTimestamp,
			CommitUrl:       stats.CommitUrl,
			CommitComments:  stats.CommitComments,
		})
	}

	sortedAndFilteredArray := FilterAndSortByCommitTimestamp(dataSortableArray, firstDateOfPeriod)

	relativeHistory := []RelativeHistoricalEvent{}

	var metricMetadata MetricMetadata = MetricMetadata{
		TimeGrain:       metric.TimeGrain,
		PeriodKey:       metric.Period,
		InitialValue:    sortedAndFilteredArray[0].KPI,
		FirstDate:       firstDateOfPeriod,
		RelativeHistory: relativeHistory,
	}
	return metricMetadata
}
