package github

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/data-drift/data-drift/helpers"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stringContentUrl := content.GetDownloadURL()

	resp, err := http.Get(stringContentUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(records) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no records in CSV file"})
		return
	}

	firstRecord := records[0]

	patch := csvFile.GetPatch()
	patchToLarge := false

	if patch == "" {
		patchToLarge = true
		patch, err = getPatchIfEmpty(client, c, owner, repo, commit, csvFile, records)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting patch when patch is empty"})
			return
		}
	}

	jsonData, err := json.Marshal(gin.H{"patch": patch, "headers": firstRecord, "filename": csvFile.GetFilename(), "date": commit.GetCommit().GetCommitter().GetDate(), "commitLink": commit.GetHTMLURL(), "patchToLarge": patchToLarge})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error marshaling JSON"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
}

func getPatchIfEmpty(client *github.Client, ctx *gin.Context, owner string, repo string, commit *github.RepositoryCommit, file *github.CommitFile, currentRecord [][]string) (string, error) {
	parentCommitSha := commit.Parents[0].GetSHA()
	previousFileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, *file.Filename, &github.RepositoryContentGetOptions{Ref: parentCommitSha})
	if err != nil {
		fmt.Println("Error getting github file content:", err)
		return "", err
	}
	stringContentUrl := previousFileContent.GetDownloadURL()
	resp, err := http.Get(stringContentUrl)
	if err != nil {
		fmt.Println("Error getting file from url:", err)
		return "", err
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	previousRecords, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading csv:", err)
		return "", err
	}
	patch, err := helpers.GenerateCsvPatch(currentRecord, previousRecords)
	lines := strings.Split(patch, "\n")
	if len(lines) > 10000 {
		lines = lines[:10000]
	}
	patch = strings.Join(lines, "\n")
	return patch, err
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	date := c.Query("date")
	if date != "" {
		start, _ := time.Parse(time.RFC3339, date+"T00:00:00Z")
		end, _ := time.Parse(time.RFC3339, date+"T23:59:59Z")
		opt.Since = start
		opt.Until = end
	}

	commits, _, err := client.Repositories.ListCommits(c, owner, repo, opt)
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
