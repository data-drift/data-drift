package local_store

import (
	"log"
	"net/http"
	"strings"

	"github.com/data-drift/data-drift/common"

	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func MetricHandler(c *gin.Context) {
	store := c.Param("store")
	table := c.Param("table")
	metricName := "metric_value"
	// Est-ce que je fait pour des grosse table metric ? est-ce que je fais pour des tables aggregée ? je veux rapidement avoir un resultat montrable
	metricHistory, err := getMetricHistory(store, table, metricName)
	print(metricHistory)
	print(err)
	c.JSON(http.StatusOK, gin.H{
		"store":         store,
		"table":         table,
		"metricHistory": metricHistory,
	})
}

func getMetricHistory(store string, table string, metricName string) ([]common.MetricHistoryEvent, error) {
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

	var history []common.MetricHistoryEvent
	// Filter the commit history to only include commits that modified the file

	err = commitIter.ForEach(func(commit *object.Commit) error {
		file, _ := commit.File(filePath)
		content, err := file.Contents()
		if err != nil {
			print(err)
			return err
		}
		records := strings.Split(content, "\n")
		print(records[0])

		log.Println("it worked")
		history = append(history, common.MetricHistoryEvent{})

		return nil
	})
	if err != nil {
		print(err)
	}
	return history, nil
}
