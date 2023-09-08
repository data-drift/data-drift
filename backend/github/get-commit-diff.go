package github

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v42/github"
)

func GetCommitDiff(c *gin.Context) {
	InstallationId, err := strconv.ParseInt(c.Request.Header.Get("Installation-Id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing Installation-Id header"})
		return
	}
	owner := c.Param("owner")
	repo := c.Param("repo")
	commitSha := c.Param("commit-sha")

	client, err := CreateClientFromGithubApp(int64(InstallationId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	commit, _, ghErr := client.Repositories.GetCommit(c, owner, repo, commitSha, nil)
	if ghErr != nil {
		fmt.Println(ghErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ghErr})
		return
	}
	if len(commit.Files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files in commit"})
		return
	}
	var csvFile *github.CommitFile
	for _, file := range commit.Files {
		if strings.HasSuffix(file.GetFilename(), ".csv") {
			csvFile = file
			break
		}
	}

	if csvFile == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no CSV files in commit"})
		return
	}

	content, _, _, err := client.Repositories.GetContents(c, owner, repo, csvFile.GetFilename(), &github.RepositoryContentGetOptions{Ref: commitSha})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	stringContent, _ := content.GetContent()
	contentReader := strings.NewReader(stringContent)
	csvReader := csv.NewReader(contentReader)
	records, err := csvReader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if len(records) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no records in CSV file"})
		return
	}

	firstRecord := records[0]

	patch := csvFile.GetPatch()

	jsonData, err := json.Marshal(gin.H{"patch": patch, "headers": firstRecord, "filename": csvFile.GetFilename(), "date": commit.GetCommit().GetCommitter().GetDate(), "commitLink": commit.GetHTMLURL()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error marshaling JSON"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
}

func GetCommitList(c *gin.Context) {
	InstallationId, err := strconv.ParseInt(c.Request.Header.Get("Installation-Id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing Installation-Id header"})
		return
	}
	owner := c.Param("owner")
	repo := c.Param("repo")

	client, err := CreateClientFromGithubApp(int64(InstallationId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	commits, _, err := client.Repositories.ListCommits(c, owner, repo, nil)
	if err != nil {
		return
	}

	jsonData, err := json.Marshal(commits)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error marshaling JSON"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
}
