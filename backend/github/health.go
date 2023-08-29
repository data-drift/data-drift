package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CheckGithubAppConnectionForInstallation(installationId int64) (string, error) {

	client, err := CreateClientFromGithubApp(installationId)
	if err != nil {
		return "", err
	}
	ctx := context.Background()

	// Get the last commit of the repository
	result, _, err := client.Apps.ListRepos(ctx, nil)
	if err != nil {
		return "", err
	}

	for _, repo := range result.Repositories {
		fmt.Println("repo:", repo.GetName())
	}
	return "ok", nil
}

func CheckGithubAppConnection() error {
	privateKeyPath := os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH")
	privateKey := os.Getenv("GITHUB_APP_PRIVATE_KEY")

	if privateKeyPath == "" && privateKey == "" {
		return fmt.Errorf("missing GitHub App private key information, please provide GITHUB_APP_PRIVATE_KEY_PATH or GITHUB_APP_PRIVATE_KEY")
	}

	appIDStr := os.Getenv("GITHUB_APP_ID")
	if appIDStr == "" {
		return fmt.Errorf("missing GitHub App ID, please provide GITHUB_APP_ID")
	}
	githubAppId, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		return err
	}
	fmt.Println("Github App ID: ", githubAppId)
	return nil
}

func HealthCheck(c *gin.Context) {
	err := CheckGithubAppConnection()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})

}

func HealthCheckInstallation(c *gin.Context) {
	installationIdStr := c.Param("installation-id")
	installationId, err := strconv.ParseInt(installationIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"status": "ERROR", "error": err.Error()})
	}
	sha, err := CheckGithubAppConnectionForInstallation(installationId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "commit": sha})

}
