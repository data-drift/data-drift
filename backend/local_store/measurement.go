package local_store

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

	var req MeasurementRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commit, err := getMeasurement(store, table, measurementId)

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

	file, _ := commit.File(table + ".csv")
	content, _ := file.Contents()
	getMetricByTimeGrain(req, content)

	c.JSON(http.StatusOK, gin.H{"MeasurementMetaData": measurementMetaData})

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

func getMeasurement(store string, table string, commitSha string) (*object.Commit, error) {
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

	if err != nil {
		print("Error getting HEAD reference")
		return nil, err
	}
	hash := plumbing.NewHash(commitSha)

	commit, err := repo.CommitObject(hash)
	if err != nil {
		return nil, err
	}
	_, err = commit.File(filePath)
	if err != nil {
		return nil, fmt.Errorf("file not present in measurement")
	}
	return commit, err
}

func getMetricByTimeGrain(measurementRequest MeasurementRequest, fileContent string) error {
	reader := csv.NewReader(bufio.NewReader(strings.NewReader(fileContent)))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	fmt.Println(records[0])
	return nil
}
