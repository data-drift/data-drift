package reducers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/data-drift/kpi-git-history/common"
	"github.com/data-drift/kpi-git-history/helpers"
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

func mapChartDataToDatasets(chartData map[time.Duration]RelativeHistoricalEvent) []map[string]interface{} {

	var data []map[string]interface{}
	for _, point := range chartData {

		data = append(data, map[string]interface{}{
			"x": helpers.GetFloat(point.DaysFromHistorization),
			"y": helpers.GetFloat(point.RelativeValue),
		})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i]["x"].(float64) < data[j]["x"].(float64)
	})

	return removeDuplicatesByY(data)
}

func CreateMetadataChart(metricMetadatas map[common.PeriodKey]MetricMetadata) string {
	datasets := []map[string]interface{}{}
	for _, metricMetadata := range metricMetadatas {

		datasets = append(datasets, map[string]interface{}{
			"label":           metricMetadata.PeriodKey,
			"showLine":        true,
			"lineTension":     0,
			"borderColor":     helpers.GetColorFromString(string(metricMetadata.PeriodKey)),
			"backgroundColor": "transparent",
			"data":            mapChartDataToDatasets(metricMetadata.RelativeHistory),
		})
	}
	jsonBody := map[string]interface{}{
		"version":          "4",
		"backgroundColor":  "transparent",
		"width":            500,
		"height":           300,
		"devicePixelRatio": 2.0,
		"format":           "svg",
		"chart": map[string]interface{}{
			"type": "scatter",
			"data": map[string]interface{}{

				"datasets": datasets,
			},
		},
	}

	helpers.WriteMetadataToFile(jsonBody, "dist/payload.json")

	newData, _ := json.Marshal(jsonBody)
	url := "https://quickchart.io/chart/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newData))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())

	var chartResponse ChartResponse
	jsonUnmarshalError := json.Unmarshal(buf.Bytes(), &chartResponse)
	if jsonUnmarshalError != nil {
		fmt.Println("Error parsing JSON:", err.Error())
		return "" // Return an empty string or handle the error as needed
	}

	interactiveUrl := convertToChartMakerURL(chartResponse.URL)
	fmt.Println("Interactive URL:", interactiveUrl)

	// Return only the URL
	return interactiveUrl
}

func removeDuplicatesByY(data []map[string]interface{}) []map[string]interface{} {

	var result []map[string]interface{}
	var lastKPI float64
	for _, point := range data {
		if lastKPI == 0 || point["y"] != lastKPI {
			result = append(result, point)
			lastKPI = point["y"].(float64)
		}
	}

	return result
}