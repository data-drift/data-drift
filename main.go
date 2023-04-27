// main.go
package main

import (
	"fmt"

	"github.com/data-drift/kpi-git-history/charts"
	"github.com/data-drift/kpi-git-history/history"
	"github.com/data-drift/kpi-git-history/reports"
)

func main() {
	// Call your custom function here
	performTask()
}

func performTask() {
	filepath := history.ProcessHistory()
	// Call functions from charts.go and reports.go
	chartResults := charts.ProcessCharts(filepath)

	for _, chartResult := range chartResults {
		reports.CreateReport(chartResult)
	}
	// ...
	fmt.Println("Custom function completed. Chart result:", filepath)
	fmt.Println("Custom function completed. Chart result:", chartResults)
}
