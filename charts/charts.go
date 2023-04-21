package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	url := "https://quickchart.io/chart/create"
	jsonStr := []byte(`{
		"backgroundColor": "#fff",
		"width": 500,
		"height": 300,
		"devicePixelRatio": 1.0,
		"chart": {
		  "type": "bar",
		  "data": {
			"labels": [2012, 2013, 2014, 2015, 2016],
			"datasets": [
			  {
				"label": "Users",
				"data": [750, [500, 750], [360, 500], [200, 360], 200]
			  }
			]
		  }
		}
	  }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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
