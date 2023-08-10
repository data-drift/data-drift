// main.go
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/debug"
	"github.com/data-drift/data-drift/github"
	"github.com/data-drift/data-drift/history"
	"github.com/data-drift/data-drift/metrics"
	"github.com/data-drift/data-drift/reducers"
	"github.com/data-drift/data-drift/reports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var debugEnabled bool

func init() {
	// Parse command line flags
	flag.BoolVar(&debugEnabled, "debug", false, "Enable debug mode")
	flag.Parse()
}

func main() {
	godotenv.Load()

	if debugEnabled {
		debug.DebugFunction()
		return
	}

	port := defaultIfEmpty(os.Getenv("PORT"), "8080")

	router := gin.New()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	// config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowHeaders = append(config.AllowHeaders, "Installation-Id")
	router.Use(cors.New(config))

	router.Use(gin.Logger())

	router.GET("/", HealthCheck)

	router.POST("/", ManualSync)

	router.GET("/ghhealth", github.HealthCheck)
	router.GET("/ghhealth/:installation-id", github.HealthCheckInstallation)

	router.POST("webhooks/github", github.HandleWebhook)
	router.GET("gh/:owner/:repo/commit/:commit-sha", github.GetCommitDiff)

	router.GET("metrics/:metric-name/cohorts/:timegrain", metrics.GetMetricCohort)

	router.POST("validate-config", github.ValidateConfigHandler)

	router.Run(":" + port)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func ManualSync(c *gin.Context) {
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
}

func performTask(syncConfig common.SyncConfig) error {
	filepath := os.Getenv("DEFAULT_FILE_PATH")
	fmt.Println(filepath)

	if filepath == "" {
		client := github.CreateClientFromGithubToken(syncConfig.GithubToken)
		newFilepath, err := history.ProcessHistory(client,
			syncConfig.GithubRepoOwner,
			syncConfig.GithubRepoName,
			common.MetricConfig{
				MetricName:     "Default metric name",
				KPIColumnName:  syncConfig.KpiColumn,
				DateColumnName: syncConfig.DateColumn,
				Filepath:       syncConfig.GithubRepoFilePath,
				TimeGrains:     []common.TimeGrain{common.Day},
				Dimensions:     []string{},
			},
			int(0),
		)

		if err != nil {
			return err
		}
		filepath = newFilepath
	}
	// Call functions from charts.go and reports.go
	chartResults := reducers.ProcessMetricHistory(filepath, common.MetricConfig{MetricName: "Default metric name"})

	for _, chartResult := range chartResults {
		err := reports.CreateReport(syncConfig, chartResult)
		if err != nil {
			fmt.Println("[DATADRIFT_ERROR]", err.Error())
		}
	}
	fmt.Println("Custom function completed. Chart result:", filepath)
	fmt.Println("Custom function completed. Chart result:", chartResults)
	return nil
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
