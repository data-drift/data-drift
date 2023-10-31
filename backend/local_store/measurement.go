package local_store

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type MeasurementRequest struct {
	Metric    string           `json:"metric"`
	TimeGrain common.TimeGrain `json:"timegrain"`
}

func MeasurementsHandler(c *gin.Context) {
	store := c.Param("store")
	table := c.Param("table")
	queryDate := c.Query("date")
	date, err := time.Parse("2006-01-02", queryDate)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	measurements, err := getMeasurements(store, table, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Measurements": measurements})
}

func MeasurementHandler(c *gin.Context) {
	store := c.Param("store")
	table := c.Param("table")
	measurementId := c.Param("measurementId")

	commit, patch, headers, err := getMeasurement(store, table, measurementId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	measurementMetaData := common.MeasurementMetaData{

		MeasurementTimestamp: commit.Author.When.Unix(),
		MeasurementDate:      commit.Author.When.Format("2006-01-02"),
		MeasurementDateTime:  commit.Author.When.Format("2006-01-02 15:04:05"),
		MeasurementId:        measurementId,
		MeasurementComments: []common.CommitComments{
			{
				CommentAuthor: commit.Author.Name,
				CommentBody:   commit.Message,
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"MeasurementMetaData": measurementMetaData, "Patch": patch, "Headers": headers})
}

func getMeasurements(store string, table string, date time.Time) ([]CommitInfo, error) {
	repoDir, err := getStoreDir(store)
	filePath := table + ".csv"
	if err != nil {
		print("Error getting store directory")
		return nil, err
	}
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		print("Error opening repo")
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		print("Error fetching repo head")
		return nil, err
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), FileName: &filePath})
	if err != nil {
		print("Error getting commit log")
		return nil, err
	}

	commits := []CommitInfo{}

	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Author.When.Format("2006-01-02") == date.Format("2006-01-02") {
			commits = append(commits, CommitInfo{
				Message: c.Message,
				Date:    c.Author.When,
				Sha:     c.Hash.String(),
			})
		}
		return nil
	})
	if err != nil && err != io.EOF {
		print("Error iterating commits", err.Error())
		return nil, err
	}
	return commits, nil
}

func getMeasurement(store string, table string, commitSha string) (*object.Commit, string, []string, error) {
	repoDir, err := getStoreDir(store)
	log.Println("repoDir", repoDir)
	filePath := table + ".csv"
	if err != nil {
		print("Error getting store directory")
		return nil, "", nil, err
	}
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		print("Error opening repo")
		return nil, "", nil, err
	}

	if err != nil {
		print("Error getting HEAD reference")
		return nil, "", nil, err
	}
	hash := plumbing.NewHash(commitSha)

	commit, err := repo.CommitObject(hash)
	if err != nil {
		return nil, "", nil, err
	}
	file, err := commit.File(filePath)
	if err != nil {
		return nil, "", nil, fmt.Errorf("file not present in measurement")
	}

	// Retrieve the commit's parents
	parent, err := commit.Parent(0)

	log.Println("commit", commit.Hash, commit.Author.When, commit.Message)
	log.Println("parent", parent.Hash, parent.Author.When, commit.Message)
	if err != nil {
		return nil, "", nil, err
	}

	if err != nil {
		// Handle the error. For example:
		log.Fatalf("Failed to get file: %v", err)
	}

	currentContent, err := file.Contents()
	if err != nil {
		// Handle the error. For example:
		log.Fatalf("Failed to read file contents: %v", err)
	}

	// Convert the content to an io.Reader. Assuming content is a string:
	reader := csv.NewReader(strings.NewReader(currentContent))

	// read the headers from the CSV file
	currentRecord, err := reader.ReadAll()
	headers := currentRecord[0]
	if err != nil {
		// Handle the error. For example:
		log.Fatalf("Failed to read CSV headers: %v", err)
	}

	previousFile, _ := parent.File(filePath)
	previousContent, _ := previousFile.Contents()
	previousReader := csv.NewReader(strings.NewReader(previousContent))
	previousRecords, _ := previousReader.ReadAll()

	patch, err := helpers.GenerateCsvPatch(currentRecord, previousRecords)
	if err != nil {
		log.Fatalf("Failed to generate patch: %v", err)
	}
	lines := strings.Split(patch, "\n")
	if len(lines) > 10000 {
		lines = lines[:10000]
	}
	patch = strings.Join(lines, "\n")

	return commit, patch, headers, nil
}
