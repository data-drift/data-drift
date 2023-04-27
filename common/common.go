package common

type KPIInfo struct {
	KPIName    string `json:"kpiName"`
	GraphQLURL string `json:"graphqlUrl"`
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
