package github

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GithubWebhookPayload struct {
	Ref    string `json:"ref"`
	After  string `json:"after"`
	Before string `json:"before"`
	Forced bool   `json:"forced"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Sender struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		HTMLURL   string `json:"html_url"`
		AvatarURL string `json:"avatar_url"`
	} `json:"sender"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
	} `json:"commits"`
	Repository struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"repository"`
}

func HandleWebhook(c *gin.Context) {
	var payload GithubWebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Access the parsed payload fields
	ref := payload.Ref
	fmt.Println(ref)
	// Process the webhook payload...

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}
