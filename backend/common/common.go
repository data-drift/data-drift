package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

type KPIReport struct {
	KPIName        string          `json:"kpiName"`
	PeriodId       PeriodKey       `json:"periodId"`
	DimensionValue DimensionValue  `json:"dimensionValue"`
	GraphQLURL     string          `json:"graphqlUrl"`
	InitialValue   decimal.Decimal `json:"firstRoundedKPI"`
	LatestValue    decimal.Decimal `json:"lastRoundedKPI"`
	Events         []EventObject   `json:"events"`
}

type SyncConfig struct {
	GithubToken        string
	GithubRepoOwner    string
	GithubRepoName     string
	GithubRepoFilePath string
	DateColumn         string
	StartDate          string
	KpiColumn          string
	NotionAPIKey       string
	NotionDatabaseID   string
}

type EventObject struct {
	CommitTimestamp int64            `json:"commitTimestamp"`
	CommitUrl       string           `json:"commitUrl"`
	Diff            float64          `json:"diff"`
	EventType       EventType        `json:"eventType"`
	CommitComments  []CommitComments `json:"commitComments"`
}

type EventType string

type CommitComments struct {
	CommentAuthor string
	CommentBody   string
}

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
)

type CommitData struct {
	Lines           int
	KPI             decimal.Decimal
	CommitTimestamp int64
	CommitUrl       string
	CommitComments  []CommitComments
}

func (c CommitData) Timestamp() int64 {
	return c.CommitTimestamp
}

type CommitSha string
type PeriodKey string
type PeriodAndDimensionKey string
type Dimension string
type DimensionValue string
type MetricHistory map[CommitSha]CommitData
type Metric struct {
	TimeGrain      TimeGrain
	Period         PeriodKey
	Dimension      Dimension
	DimensionValue DimensionValue
	History        MetricHistory
}
type Metrics map[PeriodAndDimensionKey]Metric

type Config struct {
	NotionAPIToken   string         `json:"notionAPIToken"`
	NotionDatabaseID string         `json:"notionDatabaseId"`
	Metrics          []MetricConfig `json:"metrics"`
}

type TimeGrain string

const (
	Day     TimeGrain = "day"
	Week    TimeGrain = "week"
	Month   TimeGrain = "month"
	Quarter TimeGrain = "quarter"
	Year    TimeGrain = "year"
)

type MetricConfig struct {
	Filepath       string      `json:"filepath"`
	DateColumnName string      `json:"dateColumnName"`
	KPIColumnName  string      `json:"KPIColumnName"`
	MetricName     string      `json:"metricName"`
	TimeGrains     []TimeGrain `json:"timeGrains"`
	Dimensions     []string    `json:"dimensions"`
}

type FilePathString string

var ctx = context.Background()

func GetKeysFromJSON(path FilePathString) (Metrics, error) {
	// Read the file at the given path
	jsonFile, err := os.ReadFile(string(path))
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the desired type
	var data Metrics
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func StoreMetricMetadataAndAggregatedData(installationId int, metricName string, lineCountAndKPIByDateByVersion Metrics) FilePathString {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	metricStoredFilePath := GetMetricFilepath(fmt.Sprint(installationId), metricName, timestamp)

	file, err := os.Create(string(metricStoredFilePath))
	if err != nil {
		log.Fatalf("Error creating file: %v", err.Error())
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(lineCountAndKPIByDateByVersion); err != nil {
		log.Fatalf("Error writing JSON to file: %v", err.Error())
	}
	fmt.Println("Results written to lineCountsAndKPIs.json")
	fmt.Println("Storing results in Redis")
	// Connect to the Redis database.
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       5,  // use default DB
	})
	jsonData, err := json.Marshal(lineCountAndKPIByDateByVersion)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Err: %s", err)
	}
	err = rdb.Set(ctx, string(metricStoredFilePath), jsonData, 0).Err()
	if err != nil {
		log.Fatalf("Could not set key. Err: %s", err)
	}
	return metricStoredFilePath
}

func GetMetricFilepath(installationId string, metricName string, timestamp string) FilePathString {
	metricNameEncoded := url.PathEscape(metricName)
	filepath := fmt.Sprintf("dist/%s_%s_lineCountAndKPIByDateByVersion_%s.json", installationId, metricNameEncoded, timestamp)
	return FilePathString(filepath)
}

func GetLatestMetricFile(installationId string, metricName string) (FilePathString, error) {
	filepathPattern := fmt.Sprintf("dist/%s_%s_lineCountAndKPIByDateByVersion_*.json", installationId, metricName)
	files, err := filepath.Glob(filepathPattern)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found matching pattern %q", filepathPattern)
	}

	// Check the most recent file

	return FilePathString(files[0]), nil
}
