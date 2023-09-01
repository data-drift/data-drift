// main.go
package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/data-drift/data-drift/debug"
	"github.com/data-drift/data-drift/github"
	"github.com/data-drift/data-drift/metrics"
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

	router.GET("/ghhealth", github.HealthCheck)
	router.GET("/ghhealth/:installation-id", github.HealthCheckInstallation)

	router.POST("webhooks/github", github.HandleWebhook)
	router.GET("gh/:owner/:repo/commit/:commit-sha", github.GetCommitDiff)

	router.GET("metrics/:metric-name/cohorts/:timegrain", metrics.GetMetricCohort)
	router.GET("metrics/:metric-name/reports", metrics.GetMetricReport)

	router.POST("validate-config", github.ValidateConfigHandler)

	router.Run(":" + port)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
