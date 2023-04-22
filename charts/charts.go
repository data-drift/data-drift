package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
)

func main() {
	res, _ := getKeyFromJSON("../dist/lineCountAndKPIByDateByVersion_2023-04-21_17-09-03.json", "2022-02-01")
	var diff []interface{}
	var labels []interface{}
	var colors []interface{}
	upcolor := "rgb(100, 181, 246)"
	downcolor := "rgb(255, 107, 107)"
	var prevKPI int
	for _, v := range res {
		roundedKPI := int(math.Round(v.KPI))
		if prevKPI == 0 {
			prevKPI = roundedKPI
			labels = append(labels, "i")
			diff = append(diff, roundedKPI)
			colors = append(colors, upcolor)
		} else {
			d := roundedKPI - prevKPI
			if d == 0 {

			} else {
				diff = append(diff, []int{prevKPI, roundedKPI})
				labels = append(labels, "i")
				if prevKPI < roundedKPI {
					colors = append(colors, upcolor)
				} else {
					colors = append(colors, downcolor)
				}
			}
			prevKPI = roundedKPI
		}
	}
	fmt.Println(diff)
	createChart(diff, labels, colors)

}

func createChart(diff []interface{}, labels []interface{}, colors []interface{}) {
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
						"label":           "KPI",
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
}

func getKeyFromJSON(path string, key string) (map[string]struct {
	Lines int
	KPI   float64
}, error) {
	// Read the file at the given path
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a map[string]interface{}
	var data map[string]map[string]struct {
		Lines int
		KPI   float64
	}
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}

	// Extract the value associated with the given key
	value, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}
