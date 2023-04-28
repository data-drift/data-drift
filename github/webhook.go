package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/xeipuuv/gojsonschema"
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

type GithubWebhookPayload struct {
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"repository"`
	Repositories []struct {
		Name string `json:"name"`
	} `json:"repositories"`

	Installation struct {
		ID      int `json:"id"`
		Account struct {
			Login string `json:"login"`
		} `json:"account"`
	} `json:"installation"`
}

func HandleWebhook(c *gin.Context) {
	var payload GithubWebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("ref", payload.Installation.ID)

	InstallationId := payload.Installation.ID
	client, err := CreateClientFromGithubApp(int64(InstallationId))
	if err != nil {
		fmt.Println("wahou1")

		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx := context.Background()

	fmt.Println(payload)

	if payload.Repository.Owner.Name != "" {

		confidIsValid, err := verifyConfigFile(client, payload.Repository.Owner.Name, payload.Repository.Name, ctx)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Webhook processed", "configIsValie": confidIsValid, "installationId": InstallationId})
	} else if payload.Installation.Account.Login != "" {
		confidIsValid, err := verifyConfigFile(client, payload.Installation.Account.Login, payload.Repositories[0].Name, ctx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Webhook processed", "configIsValie": confidIsValid, "installationId": InstallationId})
	} else {
		fmt.Println("No repository or account found")

		c.JSON(http.StatusBadRequest, gin.H{"error": "No repository or account found"})
		return
	}
}

func verifyConfigFile(client *github.Client, RepoOwner string, RepoName string, ctx context.Context) (bool, error) {

	commit, _, _ := client.Repositories.GetCommit(ctx, RepoOwner, RepoName, "main")

	configFilePath := "datadrift-config.json"

	file, _, _, err := client.Repositories.GetContents(ctx, RepoOwner, RepoName, configFilePath, &github.RepositoryContentGetOptions{
		Ref: *commit.SHA,
	})

	if err != nil {
		fmt.Println(err)
		return false, err
	}
	content, _ := file.GetContent()
	schemaLoader := gojsonschema.NewReferenceLoader("file://./json-schema.json")
	documentLoader := gojsonschema.NewStringLoader(content)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if result.Errors() != nil {
		fmt.Println("result.Errors()", result.Errors())
		return false, fmt.Errorf("invalid config file")
	}
	fmt.Println(result.Valid())
	return true, nil
}

func ValidateConfigHandler(c *gin.Context) {
	// Parse the JSON configuration from the request body
	var config Config
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
