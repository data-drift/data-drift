package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	url := "https://quickchart.io/chart/create"
	jsonBody := map[string]interface{}{
		"backgroundColor":  "#fff",
		"width":            500,
		"height":           300,
		"devicePixelRatio": 1.0,
		"chart": map[string]interface{}{
			"type": "bar",
			"data": map[string]interface{}{
				"labels": []int{2012, 2013, 2014, 2015, 2016},
				"datasets": []map[string]interface{}{
					{
						"label": "Users",
						"data":  []interface{}{750, []int{500, 750}, []int{360, 500}, []int{200, 360}, 200},
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
