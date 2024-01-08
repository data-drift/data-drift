package github

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/data-drift/data-drift/helpers"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
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
		patch, err = getPatchIfEmpty(client, owner, repo, commit.Parents[0].GetSHA(), csvFile, records)
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

func CompareCommit(c *gin.Context) {
	InstallationId, err := strconv.ParseInt(c.Request.Header.Get("Installation-Id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing Installation-Id header"})
		return
	}
	owner := c.Param("owner")
	repo := c.Param("repo")
	baseCommitSha := c.Param("base-commit-sha")
	headCommitSha := c.Param("head-commit-sha")
	table := c.Query("table")
	jsonData, err := compareCommit(InstallationId, owner, repo, baseCommitSha, headCommitSha, table)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonData)
}

func CompareCommitBetweenDates(c *gin.Context) {
	InstallationId, err := strconv.ParseInt(c.Request.Header.Get("Installation-Id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing Installation-Id header"})
		return
	}
	client, err := CreateClientFromGithubApp(int64(InstallationId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	owner := c.Param("owner")
	repo := c.Param("repo")
	startDateStr := c.Query("start-date")
	table := c.Query("table")
	if !strings.HasSuffix(table, ".csv") {
		table += ".csv"
	}
	if table == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table query param is required"})
		return
	}
	endDateStr := c.Query("end-date")
	beginDate, err := time.Parse("2006-01-02", startDateStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inclusiveBeginDate := beginDate.AddDate(0, 0, 1)

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inclusiveEndDate := endDate.AddDate(0, 0, 1)
	var firstCommit, latestCommit string

	optBefore := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		Until: inclusiveBeginDate,
		Path:  table,
	}
	commitsBefore, _, err := client.Repositories.ListCommits(c, owner, repo, optBefore)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(commitsBefore) > 0 {
		firstCommit = commitsBefore[0].GetSHA()

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No commits before date for table " + table})
		return
	}

	optAfter := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		Since: beginDate,
		Until: inclusiveEndDate,
		Path:  table,
	}
	commitsAfter, _, err := client.Repositories.ListCommits(c, owner, repo, optAfter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(commitsAfter) > 0 {
		latestCommit = commitsAfter[0].GetSHA()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No commits between dates for table " + table})
		return
	}

	jsonData, err := compareCommit(InstallationId, owner, repo, firstCommit, latestCommit, table)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonData)
}

func compareCommit(InstallationId int64, owner string, repo string, baseCommitSha string, headCommitSha string, table string) ([]byte, error) {

	c := context.Background()
	client, err := CreateClientFromGithubApp(int64(InstallationId))
	if err != nil {
		return nil, err
	}

	baseCommit, _, _ := client.Repositories.GetCommit(c, owner, repo, baseCommitSha, nil)
	headCommit, _, _ := client.Repositories.GetCommit(c, owner, repo, headCommitSha, nil)

	opts := &github.ListOptions{}

	comparison, _, ghErr := client.Repositories.CompareCommits(c, owner, repo, baseCommitSha, headCommitSha, opts)
	if ghErr != nil {
		fmt.Println(ghErr.Error())
		return nil, ghErr
	}
	var csvFile *github.CommitFile

	for _, file := range comparison.Files {
		if file.GetFilename() == table {
			csvFile = file
			break
		}
	}
	if csvFile == nil {
		return nil, fmt.Errorf("table %s not updated between those dates", table)
	}
	content, _, _, err := client.Repositories.GetContents(c, owner, repo, csvFile.GetFilename(), &github.RepositoryContentGetOptions{Ref: headCommitSha})
	if err != nil {
		return nil, err
	}
	stringContentUrl := content.GetDownloadURL()

	resp, err := http.Get(stringContentUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no records in CSV file")
	}

	firstRecord := records[0]

	patchToLarge := true
	patch, err := getPatchIfEmpty(client, owner, repo, baseCommitSha, csvFile, records)
	if err != nil {
		return nil, fmt.Errorf("error getting patch when patch is empty: %v", err)
	}
	jsonData, err := json.Marshal(gin.H{"patch": patch, "headers": firstRecord, "filename": csvFile.GetFilename(), "patchToLarge": patchToLarge, "fromCommitDate": baseCommit.Commit.Committer.Date, "toCommitDate": headCommit.Commit.Committer.Date})
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func getPatchIfEmpty(client *github.Client, owner string, repo string, parentCommitSha string, file *github.CommitFile, currentRecord [][]string) (string, error) {
	ctx := context.Background()
	previousRecords, err := getPreviousRecords(parentCommitSha, client, ctx, owner, repo, file)

	if err != nil {
		fmt.Println("Error getting PreviousRecords:", err)
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

func getPreviousRecords(parentCommitSha string, client *github.Client, ctx context.Context, owner string, repo string, file *github.CommitFile) ([][]string, error) {
	previousFileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, *file.Filename, &github.RepositoryContentGetOptions{Ref: parentCommitSha})
	if err != nil {
		fmt.Println("Error getting github file content:", err)
		if errResp, ok := err.(*github.ErrorResponse); ok && errResp.Response.StatusCode == http.StatusNotFound {
			fmt.Println("File not found")
			return [][]string{{"No file"}}, nil
		}
		return nil, err
	}
	stringContentUrl := previousFileContent.GetDownloadURL()
	resp, err := http.Get(stringContentUrl)
	if err != nil {
		fmt.Println("Error getting file from url:", err)
		return nil, err
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	previousRecords, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading csv:", err)
		return nil, err
	}
	return previousRecords, nil
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
