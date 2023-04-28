package github

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckGithubAppConnection() (string, error) {

	client, err := CreateClientFromGithubApp(36944435)
	if err != nil {
		return "", err
	}
	ctx := context.Background()

	// Get the last commit of the repository
	commit, _, err := client.Repositories.GetCommit(ctx, "Samox", "copy-libeo-data-history", "main")

	if err != nil {
		return "", err
	}
	return commit.GetSHA(), nil
}

func HealthCheck(c *gin.Context) {
	sha, err := CheckGithubAppConnection()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "commit": sha})

}
