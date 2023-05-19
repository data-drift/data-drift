package reducers

import (
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
	RelativeHistory []RelativeHistoricalEvent
}

func ProcessMetricMetadata(metricConfig common.MetricConfig, metrics common.Metrics) map[common.TimeGrain]MetricMetadata {
	var metricMetadata MetricMetadata
	return map[common.TimeGrain]MetricMetadata{
		common.Month: metricMetadata,
	}
}
