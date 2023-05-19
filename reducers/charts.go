package reducers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/data-drift/kpi-git-history/common"
	"github.com/shopspring/decimal"
)

type ChartResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
}

func ProcessCharts(historyFilepath string, metric common.MetricConfig) []common.KPIReport {

	data, err := getKeysFromJSON(historyFilepath)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	var kpiInfos []common.KPIReport

	for periodIdAndDimensionKey := range data {
		key := string(periodIdAndDimensionKey)
		fmt.Println("Key:", key)
		// Access the value associated with the key: data[key]
		// Additional logic for processing the value
		// ...
		kpi := OrderDataAndCreateChart(metric.MetricName+" "+key, data[periodIdAndDimensionKey].Period, data[periodIdAndDimensionKey].History, data[periodIdAndDimensionKey].DimensionValue)
		kpiInfos = append(kpiInfos, kpi)
	}

	return kpiInfos
}

func OrderDataAndCreateChart(KPIName string, periodId common.PeriodKey, unsortedResults common.MetricHistory, dimensionValue common.DimensionValue) common.KPIReport {
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

	sortedAndFilteredArray := FilterAndSortByCommitTimestamp(dataSortableArray, getFirstDateOfPeriod(periodId))

	if len(sortedAndFilteredArray) == 0 {
		return common.KPIReport{}
	}

	var diff []interface{}
	var labels []interface{}
	var colors []interface{}
	initialcolor := "rgb(151 154 155)"
	upcolor := "rgb(82 156 202)"
	downcolor := "rgb(255 163 68)"
	var prevKPI decimal.Decimal
	initialValue := sortedAndFilteredArray[0].KPI
	latestValue := sortedAndFilteredArray[len(sortedAndFilteredArray)-1].KPI
	var events []common.EventObject
	minOfChart := 0

	for _, v := range sortedAndFilteredArray {
		roundedKPI := v.KPI
		// TODO
		roundedMin := 32000
		timestamp := int64(v.CommitTimestamp) // Unix timestamp for May 26, 2022 12:00:00 AM UTC
		timeObj := time.Unix(timestamp, 0)    // Convert the Unix timestamp to a time.Time object
		dateStr := timeObj.Format("2006-01-02")
		if prevKPI.IsZero() {
			prevKPI = v.KPI
			minOfChart = roundedMin
			labels = append(labels, dateStr)
			diff = append(diff, roundedKPI)
			colors = append(colors, initialcolor)
			event := common.EventObject{
				CommitTimestamp: timestamp,
				Diff:            0,
				EventType:       common.EventTypeCreate,
				CommitUrl:       v.CommitUrl,
				CommitComments:  v.CommitComments,
			}
			events = append(events, event)
		} else {
			d := roundedKPI.Sub(prevKPI)
			if d.IsZero() {

			} else {
				// Maybe diff does not work with float
				diff = append(diff, []decimal.Decimal{prevKPI, roundedKPI})
				labels = append(labels, dateStr)
				if prevKPI.LessThan(roundedKPI) {
					colors = append(colors, upcolor)
				} else {
					colors = append(colors, downcolor)
					minOfChart = roundedMin
				}
				diff, _ := d.Float64()
				event := common.EventObject{
					CommitTimestamp: timestamp,
					Diff:            diff,
					EventType:       common.EventTypeUpdate,
					CommitUrl:       v.CommitUrl,
					CommitComments:  v.CommitComments,
				}
				events = append(events, event)
			}
			prevKPI = roundedKPI
		}
	}
	fmt.Println(diff)

	chartUrl := createChart(diff, labels, colors, KPIName, minOfChart)
	kpi1 := common.KPIReport{
		KPIName:        KPIName,
		PeriodId:       periodId,
		DimensionValue: dimensionValue,
		GraphQLURL:     chartUrl,
		InitialValue:   initialValue,
		LatestValue:    latestValue,
		Events:         events,
	}
	return kpi1
}

func createChart(diff []interface{}, labels []interface{}, colors []interface{}, KPIDate string, minOfChart int) string {
	url := "https://quickchart.io/chart/create"
	jsonBody := map[string]interface{}{
		"version":          "4",
		"backgroundColor":  "transparent",
		"width":            250,
		"height":           150,
		"devicePixelRatio": 2.0,
		"format":           "svg",
		"chart": map[string]interface{}{
			"type": "bar",
			"data": map[string]interface{}{
				"labels": labels,

				"datasets": []map[string]interface{}{
					{
						"backgroundColor": colors,
						"label":           KPIDate,
						"data":            diff,
					},
				},
			},
			"options": map[string]interface{}{
				"scales": map[string]interface{}{
					"y": map[string]interface{}{
						"min": minOfChart,
						"ticks": map[string]interface{}{
							"font": map[string]interface{}{
								"size":   8,
								"family": "Sans-Serif Workhorse",
							},
						},
					},
					"x": map[string]interface{}{
						"ticks": map[string]interface{}{
							"font": map[string]interface{}{
								"size":   8,
								"family": "Sans-Serif Workhorse",
							},
						},
					},
				},
				"legend": map[string]interface{}{
					"display": false,
				},
				"plugins": map[string]interface{}{
					"legend": map[string]interface{}{
						"display": false,
					},
				},
			},
		},
	}

	newData, _ := json.Marshal(jsonBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

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

func convertToChartMakerURL(url string) string {
	chartMakerURL := strings.Replace(url, "chart/render", "chart-maker/view", 1)
	return chartMakerURL
}

func getKeysFromJSON(path string) (common.Metrics, error) {
	// Read the file at the given path
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the desired type
	var data common.Metrics
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
