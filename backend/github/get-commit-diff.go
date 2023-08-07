package github

import "github.com/gin-gonic/gin"

func GetCommitDiff(c *gin.Context) {
	installationId := c.Request.Header.Get("Installation-Id")
	owner := c.Param("owner")
	repo := c.Param("repo")
	commitSha := c.Param("commit-sha")

	print(installationId, owner, repo, commitSha)
}
