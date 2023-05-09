package charts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/data-drift/kpi-git-history/common"
)

type ChartResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
}

func ProcessCharts(historyFilepath string, metric common.Metric) []common.KPIReport {

	data, err := getKeysFromJSON(historyFilepath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var kpiInfos []common.KPIReport

	for key := range data {
		fmt.Println("Key:", key)
		// Access the value associated with the key: data[key]
		// Additional logic for processing the value
		// ...
		kpi := OrderDataAndCreateChart(metric.MetricName+" "+key, data[key])
		kpiInfos = append(kpiInfos, kpi)
	}

	return kpiInfos
}

func OrderDataAndCreateChart(KPIName string, unsortedResults map[string]struct {
	Lines           int
	KPI             float64
	CommitTimestamp int64
	CommitUrl       string
}) common.KPIReport {
	// Extract the values from the map into a slice of struct objects
	var dataSortableArray []struct {
		Lines           int
		KPI             float64
		CommitTimestamp int64
		CommitUrl       string
	}
	for _, stats := range unsortedResults {
		dataSortableArray = append(dataSortableArray, struct {
			Lines           int
			KPI             float64
			CommitTimestamp int64
			CommitUrl       string
		}{
			Lines:           stats.Lines,
			KPI:             stats.KPI,
			CommitTimestamp: stats.CommitTimestamp,
			CommitUrl:       stats.CommitUrl,
		})
	}

	// Sort the slice by CommitTimestamp
	sort.Slice(dataSortableArray, func(i, j int) bool {
		return dataSortableArray[i].CommitTimestamp < dataSortableArray[j].CommitTimestamp
	})

	var diff []interface{}
	var labels []interface{}
	var colors []interface{}
	initialcolor := "rgb(151 154 155)"
	upcolor := "rgb(82 156 202)"
	downcolor := "rgb(255 163 68)"
	var prevKPI int
	var firstRoundedKPI int
	var lastRoundedKPI int
	var events []common.EventObject
	minOfChart := 0

	for _, v := range dataSortableArray {
		roundedKPI := int(math.Round(v.KPI))
		roundedMin := int(math.Round(v.KPI * 0.98))
		timestamp := int64(v.CommitTimestamp) // Unix timestamp for May 26, 2022 12:00:00 AM UTC
		timeObj := time.Unix(timestamp, 0)    // Convert the Unix timestamp to a time.Time object
		dateStr := timeObj.Format("2006-01-02")
		if prevKPI == 0 {
			firstRoundedKPI = roundedKPI
			prevKPI = roundedKPI
			minOfChart = roundedMin
			labels = append(labels, dateStr)
			diff = append(diff, roundedKPI)
			colors = append(colors, initialcolor)
			event := common.EventObject{
				CommitTimestamp: timestamp,
				Diff:            0,
				EventType:       common.EventTypeCreate,
				CommitUrl:       v.CommitUrl,
			}
			events = append(events, event)
		} else {
			d := roundedKPI - prevKPI
			if d == 0 {

			} else {
				diff = append(diff, []int{prevKPI, roundedKPI})
				labels = append(labels, dateStr)
				if prevKPI < roundedKPI {
					colors = append(colors, upcolor)
				} else {
					colors = append(colors, downcolor)
					minOfChart = roundedMin
				}
				event := common.EventObject{
					CommitTimestamp: timestamp,
					Diff:            d,
					EventType:       common.EventTypeUpdate,
					CommitUrl:       v.CommitUrl,
				}
				events = append(events, event)
			}
			prevKPI = roundedKPI
			lastRoundedKPI = roundedKPI
		}
	}
	fmt.Println(diff)

	chartUrl := createChart(diff, labels, colors, KPIName, minOfChart)
	kpi1 := common.KPIReport{
		KPIName:         KPIName,
		GraphQLURL:      chartUrl,
		FirstRoundedKPI: firstRoundedKPI,
		LastRoundedKPI:  lastRoundedKPI,
		Events:          events,
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
		fmt.Println("Error parsing JSON:", err)
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

func getKeysFromJSON(path string) (map[string]map[string]struct {
	Lines           int
	KPI             float64
	CommitTimestamp int64
	CommitUrl       string
}, error) {
	// Read the file at the given path
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the desired type
	var data map[string]map[string]struct {
		Lines           int
		KPI             float64
		CommitTimestamp int64
		CommitUrl       string
	}
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
