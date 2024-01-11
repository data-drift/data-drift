package reducers

import (
	"fmt"
	"strings"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/urlgen"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

type ChartResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
}

func ProcessMetricHistory(historyFilepath common.MetricStorageKey, redisClient *redis.Client, metric common.MetricConfig, ownerName string, repoName string) []common.KPIReport {
	kpiRepository := common.NewKpiRepository(redisClient)

	data, err := kpiRepository.ReadMetricKPI(historyFilepath)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	var kpiInfos []common.KPIReport

	for periodIdAndDimensionKey, datum := range data {
		key := string(periodIdAndDimensionKey)
		// Access the value associated with the key: data[key]
		// Additional logic for processing the value
		// ...
		kpi := OrderDataAndCreateChart(metric.MetricName+" "+key, datum.Period, datum.History, datum.DimensionValue, ownerName, repoName, metric.MetricName)
		kpiInfos = append(kpiInfos, kpi)
	}

	return kpiInfos
}

func OrderDataAndCreateChart(KPIName string, periodId common.PeriodKey, unsortedResults common.MetricHistory, dimensionValue common.DimensionValue, ownerName, repoName, metricName string) common.KPIReport {
	// Extract the values from the map into a slice of struct objects
	var dataSortableArray []common.CommitData

	for _, stats := range unsortedResults {
		KPI := stats.KPI
		dataSortableArray = append(dataSortableArray, common.CommitData{
			Lines:           stats.Lines,
			KPI:             KPI,
			CommitTimestamp: stats.CommitTimestamp,
			CommitUrl:       stats.CommitUrl,
			CommitComments:  stats.CommitComments,
		})
	}
	firstDateOfPeriod, firstDateOfPeriodErr := LegacyGetFirstComputationDateOfPeriod(periodId)
	if firstDateOfPeriodErr != nil {
		fmt.Println("Error:", firstDateOfPeriodErr.Error())
		return common.KPIReport{}
	}
	sortedAndFilteredArray := FilterAndSortByCommitTimestamp(dataSortableArray, firstDateOfPeriod)

	if len(sortedAndFilteredArray) == 0 {
		return common.KPIReport{}
	}

	var prevKPI decimal.Decimal
	initialValue := sortedAndFilteredArray[0].KPI
	latestValue := sortedAndFilteredArray[len(sortedAndFilteredArray)-1].KPI
	var events []common.EventObject

	for index, v := range sortedAndFilteredArray {
		roundedKPI := v.KPI
		timestamp := int64(v.CommitTimestamp) // Unix timestamp for May 26, 2022 12:00:00 AM UTC
		if index == 0 {
			prevKPI = v.KPI
			event := common.EventObject{
				CommitTimestamp: timestamp,
				Diff:            0,
				Current:         v.KPI,
				EventType:       common.EventTypeCreate,
				CommitUrl:       v.CommitUrl,
				CommitComments:  v.CommitComments,
			}
			events = append(events, event)
		} else {
			d := roundedKPI.Sub(prevKPI)
			if d.IsZero() {

			} else {

				diff, _ := d.Float64()
				event := common.EventObject{
					CommitTimestamp: timestamp,
					Diff:            diff,
					Current:         v.KPI,
					EventType:       common.EventTypeUpdate,
					CommitUrl:       v.CommitUrl,
					CommitComments:  v.CommitComments,
				}
				events = append(events, event)
			}
			prevKPI = roundedKPI
		}
	}
	dimensionValueForUrl := ""
	if dimensionValue != common.NoDimensionValue {
		dimensionValueForUrl = string(dimensionValue)
	}
	waterfallChartUrl := urlgen.MetricReportUrl(ownerName, repoName, metricName, periodId, dimensionValueForUrl)
	kpi1 := common.KPIReport{
		KPIName:           KPIName,
		PeriodId:          periodId,
		DimensionValue:    dimensionValue,
		WaterfallChartUrl: waterfallChartUrl,
		InitialValue:      initialValue,
		LatestValue:       latestValue,
		Events:            events,
	}
	return kpi1
}

func convertToChartMakerURL(url string) string {
	chartMakerURL := strings.Replace(url, "chart/render", "chart-maker/view", 1)
	return chartMakerURL
}
