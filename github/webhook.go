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

		// Get the last commit of the repository
		commit, _, err := client.Repositories.GetCommit(ctx, payload.Repository.Owner.Name, payload.Repository.Name, "main")

		if err != nil {
			fmt.Println("wahou 2")
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Webhook processed", "commit": commit.GetSHA(), "installationId": InstallationId})
	} else if payload.Installation.Account.Login != "" {
		// Get the last commit of the repository
		commit, _, err := client.Repositories.GetCommit(ctx, payload.Installation.Account.Login, payload.Repositories[0].Name, "main")

		if err != nil {
			fmt.Println("wahou 3")
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Webhook processed", "commit": commit.GetSHA(), "installationId": InstallationId})
	} else {
		fmt.Println("No repository or account found")

		c.JSON(http.StatusBadRequest, gin.H{"error": "No repository or account found"})
		return
	}
}
