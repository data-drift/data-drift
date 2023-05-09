package common

type KPIReport struct {
	KPIName      string        `json:"kpiName"`
	GraphQLURL   string        `json:"graphqlUrl"`
	InitialValue float64       `json:"firstRoundedKPI"`
	LatestValue  float64       `json:"lastRoundedKPI"`
	Events       []EventObject `json:"events"`
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
	CommitTimestamp int64     `json:"commitTimestamp"`
	CommitUrl       string    `json:"commitUrl"`
	Diff            float64   `json:"diff"`
	EventType       EventType `json:"eventType"`
}

type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
)

type Config struct {
	NotionAPIToken   string   `json:"notionAPIToken"`
	NotionDatabaseID string   `json:"notionDatabaseId"`
	Metrics          []Metric `json:"metrics"`
}

type TimeGrain string

const (
	Day   TimeGrain = "day"
	Week  TimeGrain = "week"
	Month TimeGrain = "month"
	Year  TimeGrain = "year"
)

type Metric struct {
	Filepath       string      `json:"filepath"`
	DateColumnName string      `json:"dateColumnName"`
	KPIColumnName  string      `json:"KPIColumnName"`
	MetricName     string      `json:"metricName"`
	TimeGrains     []TimeGrain `json:"timeGrains"`
	Dimensions     []string    `json:"dimensions"`
}
