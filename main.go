// main.go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/data-drift/kpi-git-history/charts"
	"github.com/data-drift/kpi-git-history/common"
	"github.com/data-drift/kpi-git-history/github"
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
		var syncConfig common.SyncConfig
		err := c.ShouldBindJSON(&syncConfig)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid syncConfig JSON"})
			return
		}

		err = performTask(syncConfig)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	router.GET("/ghhealth", func(c *gin.Context) {
		sha, err := github.CheckGithubAppConnection()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "commit": sha})

	})

	router.Run(":" + port)
}

func performTask(syncConfig common.SyncConfig) error {
	err, filepath := history.ProcessHistory(syncConfig)
	if err != nil {
		return err
	}
	// Call functions from charts.go and reports.go
	chartResults := charts.ProcessCharts(filepath)

	for _, chartResult := range chartResults {
		err = reports.CreateReport(syncConfig, chartResult)
		if err != nil {
			return err
		}
	}
	// ...
	fmt.Println("Custom function completed. Chart result:", filepath)
	fmt.Println("Custom function completed. Chart result:", chartResults)
	return nil
}
