// main.go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/data-drift/kpi-git-history/charts"
	"github.com/data-drift/kpi-git-history/history"
	"github.com/data-drift/kpi-git-history/reports"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	router.POST("/", func(c *gin.Context) {
		performTask()
	})

	router.Run(":" + port)
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
