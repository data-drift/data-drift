package common

type KPIInfo struct {
	KPIName         string        `json:"kpiName"`
	GraphQLURL      string        `json:"graphqlUrl"`
	FirstRoundedKPI int           `json:"firstRoundedKPI"`
	LastRoundedKPI  int           `json:"lastRoundedKPI"`
	Events          []EventObject `json:"events"`
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
	Diff            int       `json:"diff"`
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

type Metric struct {
	Filepath       string `json:"filepath"`
	DateColumnName string `json:"dateColumnName"`
	KPIColumnName  string `json:"KPIColumnName"`
	MetricName     string `json:"metricName"`
}
