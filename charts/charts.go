package charts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/data-drift/kpi-git-history/common"
)

type ChartResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
}

func ProcessCharts(historyFilepath string) []common.KPIInfo {

	data, err := getKeysFromJSON(historyFilepath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var kpiInfos []common.KPIInfo

	for key := range data {
		fmt.Println("Key:", key)
		// Access the value associated with the key: data[key]
		// Additional logic for processing the value
		// ...
		kpi := OrderDataAndCreateChart("MRR "+key, data[key])
		kpiInfos = append(kpiInfos, kpi)
	}

	return kpiInfos
}

func OrderDataAndCreateChart(KPIName string, unsortedResults map[string]struct {
	Lines           int
	KPI             float64
	CommitTimestamp int64
}) common.KPIInfo {
	// Extract the values from the map into a slice of struct objects
	var dataSortableArray []struct {
		Lines           int
		KPI             float64
		CommitTimestamp int64
	}
	for _, stats := range unsortedResults {
		dataSortableArray = append(dataSortableArray, struct {
			Lines           int
			KPI             float64
			CommitTimestamp int64
		}{
			Lines:           stats.Lines,
			KPI:             stats.KPI,
			CommitTimestamp: stats.CommitTimestamp,
		})
	}

	// Sort the slice by CommitTimestamp
	sort.Slice(dataSortableArray, func(i, j int) bool {
		return dataSortableArray[i].CommitTimestamp < dataSortableArray[j].CommitTimestamp
	})

	var diff []interface{}
	var labels []interface{}
	var colors []interface{}
	upcolor := "rgb(82 156 202)"
	downcolor := "rgb(255 163 68)"
	var prevKPI int
	var firstRoundedKPI int
	var lastRoundedKPI int
	var events []common.EventObject

	for _, v := range dataSortableArray {
		roundedKPI := int(math.Round(v.KPI))
		timestamp := int64(v.CommitTimestamp) // Unix timestamp for May 26, 2022 12:00:00 AM UTC
		timeObj := time.Unix(timestamp, 0)    // Convert the Unix timestamp to a time.Time object
		dateStr := timeObj.Format("2006-01-02")
		if prevKPI == 0 {
			firstRoundedKPI = roundedKPI
			prevKPI = roundedKPI
			labels = append(labels, dateStr)
			diff = append(diff, roundedKPI)
			colors = append(colors, upcolor)
			event := common.EventObject{
				CommitTimestamp: timestamp,
				Diff:            0,
				EventType:       common.EventTypeCreate,
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
				}
				event := common.EventObject{
					CommitTimestamp: timestamp,
					Diff:            d,
					EventType:       common.EventTypeUpdate,
				}
				events = append(events, event)
			}
			prevKPI = roundedKPI
			lastRoundedKPI = roundedKPI
		}
	}
	fmt.Println(diff)

	chartUrl := createChart(diff, labels, colors, KPIName)
	kpi1 := common.KPIInfo{
		KPIName:         KPIName,
		GraphQLURL:      chartUrl,
		FirstRoundedKPI: firstRoundedKPI,
		LastRoundedKPI:  lastRoundedKPI,
		Events:          events,
	}
	return kpi1
}

func createChart(diff []interface{}, labels []interface{}, colors []interface{}, KPIDate string) string {
	url := "https://quickchart.io/chart/create"
	jsonBody := map[string]interface{}{
		"backgroundColor":  "#fff",
		"width":            500,
		"height":           300,
		"devicePixelRatio": 1.0,
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
					"yAxes": []map[string]interface{}{
						{"suggestedMin": 35000},
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

	// Return only the URL
	return chartResponse.URL
}

func getKeysFromJSON(path string) (map[string]map[string]struct {
	Lines           int
	KPI             float64
	CommitTimestamp int64
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
	}
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
