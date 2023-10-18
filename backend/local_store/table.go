package local_store

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type CommitInfo struct {
	Message string
	Date    time.Time
	Sha     string
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

		log.Println(commit.Hash.String())
		commits = append(commits, CommitInfo{
			Message: commit.Message,
			Date:    commit.Author.When,
			Sha:     commit.Hash.String(),
		})
		log.Println("it worked")

		return nil
	})
	if err != nil && err.Error() != "EOF" {
		print("Error iterating commits", err.Error())
		return nil, err
	}

	return commits, nil
}
