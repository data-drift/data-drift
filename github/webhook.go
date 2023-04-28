package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GithubWebhookPayload struct {
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"repository"`
	Installation struct {
		ID int `json:"id"`
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

	if payload.Repository.Owner.Name != "" {

		// Get the last commit of the repository
		commit, _, err := client.Repositories.GetCommit(ctx, payload.Repository.Owner.Name, payload.Repository.Name, "main")

		if err != nil {
			fmt.Println("wahou 2")
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Webhook processed", "commit": commit.GetSHA(), "installationId": InstallationId})
	}
}
