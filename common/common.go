package common

type KPIInfo struct {
	KPIName         string `json:"kpiName"`
	GraphQLURL      string `json:"graphqlUrl"`
	FirstRoundedKPI int    `json:"firstRoundedKPI"`
	LastRoundedKPI  int    `json:"lastRoundedKPI"`
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
