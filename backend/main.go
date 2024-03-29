// main.go
package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/debug"
	"github.com/data-drift/data-drift/github"
	"github.com/data-drift/data-drift/local_store"
	"github.com/data-drift/data-drift/metrics"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&github.GithubConnection{})

	GithubService := github.NewGithubService(db)

	redisClient, err := common.GetRedisClient()
	if err != nil {
		panic("failed to connect redis")
	}

	KpiRepository := common.NewKpiRepository(redisClient)

	metricsService := metrics.NewMetricService(KpiRepository)

	port := defaultIfEmpty(os.Getenv("PORT"), "8080")

	go github.ProcessWebhooks(redisClient)

	router := gin.New()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	// config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowHeaders = append(config.AllowHeaders, "Installation-Id")
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	router.Use(cors.New(config))

	router.Use(gin.Logger())

	router.GET("/ghhealth", github.HealthCheck)
	router.GET("/ghhealth/:installation-id", github.HealthCheckInstallation)

	router.POST("webhooks/github", GithubService.HandleWebhook)
	router.GET("gh/:owner/:repo/commit/:commit-sha", GithubService.GithubClientGuard, github.GetCommitDiff)
	router.GET("gh/:owner/:repo/compare/:base-commit-sha/:head-commit-sha", GithubService.GithubClientGuard, github.CompareCommit)
	router.GET("gh/:owner/:repo/compare-between-date", GithubService.GithubClientGuard, github.CompareCommitBetweenDates)
	router.GET("gh/:owner/:repo/commits", GithubService.GithubClientGuard, github.GetCommitList)
	router.GET("gh/:owner/:repo/metrics/:metric-name/cohorts/:timegrain", GithubService.GithubClientGuard, metricsService.GetMetricCohort)
	router.GET("gh/:owner/:repo/metrics/:metric-name/reports", GithubService.GithubClientGuard, metricsService.GetMetricReport)
	router.GET("config/:owner/:repo", GithubService.GithubClientGuard, github.GetConfigHandler)

	router.GET("metrics/:metric-name/cohorts/:timegrain", metricsService.GetMetricCohort)
	router.GET("metrics/:metric-name/reports", metricsService.GetMetricReport)
	router.GET("stores/:store/tables", local_store.TablesHandler)
	router.GET("stores/:store/tables/:table", local_store.TableHandler)
	router.POST("stores/:store/tables/:table", local_store.StoreTableHandler)
	router.POST("stores/:store/tables/:table/metrics", local_store.MetricHandler)
	router.GET("stores/:store/tables/:table/measurements", local_store.MeasurementsHandler)
	router.GET("stores/:store/tables/:table/measurements/:measurementId", local_store.MeasurementHandler)

	router.POST("validate-config", github.ValidateConfigHandler)

	staticFilesPath := "./dist-app"
	router.Static("/assets", filepath.Join(staticFilesPath, "assets"))
	router.StaticFile("/logo.png", filepath.Join(staticFilesPath, "logo.png"))
	// If the route does not match any API or static file, serve index.html
	// This is useful for handling HTML5 history API used in single-page applications.
	router.NoRoute(func(c *gin.Context) {
		print("No route")
		c.File(filepath.Join(staticFilesPath, "index.html"))
	})

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
