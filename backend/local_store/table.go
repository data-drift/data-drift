package local_store

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitInfo struct {
	Message string
	Date    time.Time
	Sha     string
}

func StoreTableHandler(c *gin.Context) {

	store := c.Param("store")
	table := c.Param("table")
	fileName := table + ".csv"

	// Get the commit message and date from the request
	commitMessage := c.PostForm("commitMessage")
	commitDateString := c.PostForm("commitDateRFC3339")
	commitDate, err := time.Parse(time.RFC3339, commitDateString)
	if err != nil {
		commitDate = time.Now()
	}
	log.Println("commitDate:", commitDate)

	file, err := c.FormFile("csvfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repoDir, err := getStoreDir(store)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.SaveUploadedFile(file, repoDir+"/"+fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			repo, err = git.PlainInit(repoDir, false)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			wt, err := repo.Worktree()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			commit, err := wt.Commit("Init DB", &git.CommitOptions{
				Author: &object.Signature{
					Name: "Driftdb",
					When: time.Now(),
				},
				AllowEmptyCommits: true,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fmt.Printf("Init repo with commit %s\n", commit)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	wt, err := repo.Worktree()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = wt.Add(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status, err := wt.Status()
	if err != nil {
		log.Fatalf("Failed to get status of working tree: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if status.IsClean() {
		log.Println("No changes to commit")
		c.JSON(http.StatusAlreadyReported, gin.H{"message": "No changes to commit"})
		return
	}
	commit, err := wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name: "Driftdb",
			When: commitDate,
		},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	obj, err := repo.CommitObject(commit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"commit": obj.Hash.String()})
}

func TableHandler(c *gin.Context) {
	store := c.Param("store")
	table := c.Param("table")
	tableColumns, err := getListOfColumnsFromTable(store, table)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	commits, err := getCommitsForFile(store, table+".csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"store":        store,
		"table":        table,
		"tableColumns": tableColumns,
		"commits":      commits,
	})

}

func getListOfColumnsFromTable(store string, table string) ([]string, error) {
	repoDir, err := getStoreDir(store)
	if err != nil {
		fmt.Println("Error getting store directory:", err)
		return nil, err
	}

	filename := table + ".csv"
	file, err := os.Open(filepath.Join(repoDir, filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	headers := records[0]

	var columns []string

	for _, header := range headers {
		if header != "date" && header != "unique_key" {
			columns = append(columns, header)
		}
	}
	return columns, nil
}

func getCommitsForFile(store string, filePath string) ([]CommitInfo, error) {
	repoDir, err := getStoreDir(store)
	if err != nil {
		print("Error getting store directory")
		return nil, err
	}
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		print("Error opening repo")
		return nil, err
	}

	// Get the HEAD reference
	if err != nil {
		print("Error getting HEAD reference")
		return nil, err
	}

	// Get the commit history for the file
	commitIter, err := repo.Log(&git.LogOptions{FileName: &filePath})
	if err != nil {
		print("Error getting commit history")
		return nil, err
	}

	// Filter the commit history to only include commits that modified the file

	var commits []CommitInfo

	err = commitIter.ForEach(func(commit *object.Commit) error {

		commits = append(commits, CommitInfo{
			Message: commit.Message,
			Date:    commit.Author.When,
			Sha:     commit.Hash.String(),
		})

		return nil
	})
	if err != nil && err != io.EOF {
		print("Error iterating commits", err.Error())
		return nil, err
	}

	return commits, nil
}
