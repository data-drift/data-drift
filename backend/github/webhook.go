package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/database/notion_database"
	"github.com/data-drift/data-drift/history"
	"github.com/data-drift/data-drift/reducers"
	"github.com/data-drift/data-drift/reports"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
	"github.com/xeipuuv/gojsonschema"
	"gorm.io/gorm"
)

const configFilePath = "datadrift-config.json"

type GithubConnection struct {
	gorm.Model
	Owner          string
	Repository     string
	InstallationID int64 `gorm:"uniqueIndex"`
}

type GithubService struct {
	DB *gorm.DB
}

func NewGithubService(db *gorm.DB) *GithubService {
	return &GithubService{DB: db}
}

func (h *GithubService) HandleWebhook(c *gin.Context) {
	payload, err := github.ValidatePayload(c.Request, []byte(""))

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	switch event := event.(type) {
	case *github.PushEvent:
		fmt.Println("Installation ID: ", event.Installation.ID)

		InstallationId := *event.Installation.ID
		client, err := CreateClientFromGithubApp(InstallationId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx := context.Background()

		ownerName := *event.Repo.Owner.Name
		repoName := *event.Repo.Name
		h.DB.Create(&GithubConnection{Owner: ownerName, Repository: repoName, InstallationID: InstallationId})

		config, err := VerifyConfigFile(client, ownerName, repoName, ctx)

		fmt.Println("config", config)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("config", config)
		c.JSON(http.StatusOK, gin.H{"message": "Webhook processed", "configIsValie": config, "installationId": InstallationId})

		webhookChannel <- WebhookToProcess{config: config, InstallationId: int(InstallationId), client: client, ownerName: ownerName, repoName: repoName}

	case *github.InstallationEvent:
		fmt.Println("Installation ID: ", event.Installation.ID)

		InstallationId := *event.Installation.ID
		client, err := CreateClientFromGithubApp(InstallationId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx := context.Background()

		if len(event.Repositories) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Webhook ignored"})
			return
		}

		ownerName := *event.Installation.Account.Login
		repoName := *event.Repositories[0].Name

		h.DB.Create(&GithubConnection{Owner: ownerName, Repository: repoName, InstallationID: InstallationId})

		config, err := VerifyConfigFile(client, ownerName, repoName, ctx)

		fmt.Println("config", config)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("config", config)
		c.JSON(http.StatusOK, gin.H{"message": "Webhook processed", "configIsValie": config, "installationId": InstallationId})

		webhookChannel <- WebhookToProcess{config: config, InstallationId: int(InstallationId), client: client, ownerName: ownerName, repoName: repoName}
		return

	case *github.PullRequestEvent:
		err := handlePullRequestOpened(event)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
		}
		return
	case *github.IssuesEvent:
		err := handleIssueOpened(event)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
		}
		return
	default:
		c.JSON(http.StatusOK, gin.H{"message": "Webhook ignored"})
		return
	}

}

// Define a type for the webhook data.
type WebhookToProcess struct {
	config         common.Config
	InstallationId int
	client         *github.Client
	ownerName      string
	repoName       string
}

var webhookChannel = make(chan WebhookToProcess, 100)

func ProcessWebhooks() {
	log.Println("Starting to consume the channel")
	for {
		webhookData := <-webhookChannel

		log.Println("Consuming the channel", webhookData.InstallationId, webhookData.client, webhookData.ownerName, webhookData.repoName)
		processWebhookInTheBackground(webhookData.config, webhookData.InstallationId, webhookData.client, webhookData.ownerName, webhookData.repoName)
		time.Sleep(10 * time.Second)
	}
}

func processWebhookInTheBackground(config common.Config, InstallationId int, client *github.Client, ownerName string, repoName string) bool {

	fmt.Println("starting sync")

	err := notion_database.AssertDatabaseHasDatadriftProperties(config.NotionDatabaseID, config.NotionAPIToken)

	if err != nil {
		fmt.Println("[DATADRIFT_ERROR] db notion", err.Error())
	}

	for _, metric := range config.Metrics {

		filepath, err := history.ProcessHistory(client, ownerName, repoName, metric, InstallationId)
		if err != nil {
			fmt.Println("[DATADRIFT_ERROR] process history", err.Error())

		}

		chartResults := reducers.ProcessMetricHistory(filepath, metric, common.GithubInstallationId(fmt.Sprint(InstallationId)))

		for _, chartResult := range chartResults {
			err = reports.CreateReport(common.SyncConfig{NotionAPIKey: config.NotionAPIToken, NotionDatabaseID: config.NotionDatabaseID}, chartResult)
			if err != nil {
				fmt.Println("[DATADRIFT_ERROR] create report", err.Error())
			}
		}

		metadataChartResults, metadataChartError := reducers.ProcessMetricMetadataCharts(filepath, metric)
		if metadataChartError != nil {
			fmt.Println("[DATADRIFT_ERROR] create summary report", metadataChartError.Error())
		} else {
			reports.CreateSummaryReport(common.SyncConfig{NotionAPIKey: config.NotionAPIToken, NotionDatabaseID: config.NotionDatabaseID}, metric, metadataChartResults, fmt.Sprint(InstallationId))
		}
	}
	return false
}

func VerifyConfigFile(client *github.Client, RepoOwner string, RepoName string, ctx context.Context) (common.Config, error) {

	repository, _, getRepoError := client.Repositories.Get(ctx, RepoOwner, RepoName)
	if getRepoError != nil {
		fmt.Println("[DATADRIFT_ERROR]", getRepoError.Error())

		return common.Config{}, getRepoError
	}
	commit, _, getCommitError := client.Repositories.GetCommit(ctx, RepoOwner, RepoName, *repository.DefaultBranch, nil)

	if getCommitError != nil {
		fmt.Println("[DATADRIFT_ERROR]", getCommitError.Error())
		return common.Config{}, getCommitError
	}

	file, _, _, err := client.Repositories.GetContents(ctx, RepoOwner, RepoName, configFilePath, &github.RepositoryContentGetOptions{
		Ref: *commit.SHA,
	})

	if err != nil {
		fmt.Println("[DATADRIFT_ERROR]", err.Error())
		return common.Config{}, err
	}
	content, _ := file.GetContent()
	schemaLoader := gojsonschema.NewReferenceLoader("file://./json-schema.json")
	documentLoader := gojsonschema.NewStringLoader(content)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println("[DATADRIFT_ERROR]", err.Error())
		return common.Config{}, err
	}
	if result.Errors() != nil {
		fmt.Println("result.Errors()", result.Errors())
		return common.Config{}, fmt.Errorf("invalid config file")
	}
	fmt.Println(result.Valid())
	var config common.Config
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		fmt.Println("[DATADRIFT_ERROR]", err.Error())
		return common.Config{}, err
	}
	return config, nil
}

func GetConfigHandler(c *gin.Context) {
	InstallationId, err := strconv.ParseInt(c.Request.Header.Get("Installation-Id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	owner := c.Param("owner")
	repo := c.Param("repo")
	client, err := CreateClientFromGithubApp(InstallationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx := context.Background()

	config, err := VerifyConfigFile(client, owner, repo, ctx)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get config"})
		return
	}
	config.NotionAPIToken = ""
	config.NotionDatabaseID = ""
	c.JSON(http.StatusOK, gin.H{"config": config})
}

func ValidateConfigHandler(c *gin.Context) {
	// Parse the JSON configuration from the request body
	var config common.Config
	err := c.ShouldBindJSON(&config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON configuration"})
		return
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://./json-schema.json")
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Failed to load JSON Schema"})

		return
	}

	// Validate the configuration against the JSON Schema
	configLoader := gojsonschema.NewGoLoader(config)
	result, err := schema.Validate(configLoader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Check validation results
	if !result.Valid() {
		validationErrors := make([]string, len(result.Errors()))
		for i, desc := range result.Errors() {
			validationErrors[i] = desc.String()
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Configuration is valid
	c.JSON(http.StatusOK, gin.H{"message": "Configuration is valid"})
}
