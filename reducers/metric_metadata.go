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
	RelativeHistory map[time.Duration]RelativeHistoricalEvent
}

func ProcessMetricMetadata(metricConfig common.MetricConfig, metrics common.Metrics) map[common.TimeGrain]map[common.PeriodKey]MetricMetadata {

	metricMetadatas := map[common.PeriodKey]MetricMetadata{}
	for _, metric := range metrics {
		if metric.TimeGrain != common.Month {
			continue
		}
		if metric.Dimension != "none" {
			continue
		}
		metricMetadata := getMetadataOfMetric(metric)
		metricMetadatas[metric.Period] = metricMetadata
	}
	return map[common.TimeGrain]map[common.PeriodKey]MetricMetadata{
		common.Month: metricMetadatas,
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

	relativeHistory := make(map[time.Duration]RelativeHistoricalEvent)
	initialValue := sortedAndFilteredArray[0].KPI
	for _, commitData := range sortedAndFilteredArray {
		durationFromFirstComputation := getDuration(commitData.CommitTimestamp, firstDateOfPeriod)
		relativeHistoricalEvent := RelativeHistoricalEvent{
			RelativeValue:         commitData.KPI.Div(initialValue).Mul(decimal.NewFromInt(100)),
			DaysFromHistorization: decimal.NewFromFloat(durationFromFirstComputation.Hours() / 24),
		}
		relativeHistory[durationFromFirstComputation] = relativeHistoricalEvent
	}

	var metricMetadata MetricMetadata = MetricMetadata{
		TimeGrain:       metric.TimeGrain,
		PeriodKey:       metric.Period,
		InitialValue:    initialValue,
		FirstDate:       firstDateOfPeriod,
		RelativeHistory: relativeHistory,
	}
	return metricMetadata
}

func getDuration(commitTimestamp int64, firstDateOfPeriod time.Time) time.Duration {
	commitTime := time.Unix(commitTimestamp, 0)
	return commitTime.Sub(firstDateOfPeriod)
}
