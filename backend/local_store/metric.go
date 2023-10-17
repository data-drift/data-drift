package local_store

import (
	"bufio"
	"encoding/csv"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/reducers"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type MetricRequest struct {
	Period common.PeriodKey `json:"period"`
	Metric string           `json:"metric"`
}

func MetricHandler(c *gin.Context) {
	store := c.Param("store")
	table := c.Param("table")
	var req MetricRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metricName := req.Metric
	periodKey := req.Period
	metricHistory, err := getMetricHistory(store, table, metricName, periodKey)
	print(err) // TODO handle error
	c.JSON(http.StatusOK, gin.H{
		"store":         store,
		"table":         table,
		"metricHistory": metricHistory,
		"periodKey":     periodKey,
	})
}

func getMetricHistory(store string, table string, metricName string, periodKey common.PeriodKey) ([]common.MetricHistoryEvent, error) {
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

	commitIter, err := repo.Log(&git.LogOptions{FileName: &filePath})
	if err != nil {
		print("Error getting commit history")
		return nil, err
	}

	var history []common.MetricHistoryEvent

	err = commitIter.ForEach(func(commit *object.Commit) error {
		file, _ := commit.File(filePath)
		content, err := file.Contents()
		if err != nil {
			print(err)
			return err
		}
		reader := csv.NewReader(bufio.NewReader(strings.NewReader(content)))
		records, err := reader.ReadAll()
		if err != nil {
			return err
		}

		metricEvent := computeMetricHistoryEvent(records, metricName, periodKey, time.Unix(commit.Author.When.Unix(), 0))

		if len(history) == 0 {
			history = append(history, metricEvent)
		} else if !history[len(history)-1].KPI.Equal(metricEvent.KPI) {
			history = append(history, metricEvent)
		}

		return nil
	})
	if err != nil {
		print(err)
	}
	return history, nil
}

func findMetricIndex(headers []string, metricName string) int {
	for i, header := range headers {
		if header == metricName {
			return i
		}
	}
	return -1
}

func computeMetricHistoryEvent(records [][]string, metricName string, periodKey common.PeriodKey, measureDate time.Time) common.MetricHistoryEvent {
	headers := records[0]
	metricIndex := findMetricIndex(headers, metricName)
	dateIndex := findMetricIndex(headers, "date")

	firstDateOfPeriod, firstDateOfNextPeriod, _, _ := reducers.GetStartDateEndDateAndNextPeriod(periodKey)
	log.Println(firstDateOfPeriod)
	log.Println(firstDateOfNextPeriod)

	isAfterPeriod := measureDate.After(firstDateOfPeriod)
	var historyEvent = common.MetricHistoryEvent{
		IsAfterPeriod:   isAfterPeriod,
		CommitTimestamp: measureDate.Unix(),
		CommitDate:      measureDate.Format("2006-01-02"),
	}

	for i := 1; i < len(records); i++ {
		record := records[i]
		dateStr := record[dateIndex]
		recordDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			print(err)
		}
		if recordDate.After(firstDateOfNextPeriod) || recordDate.Equal(firstDateOfNextPeriod) {
			continue
		}
		if recordDate.Before(firstDateOfPeriod) {
			continue
		}
		kpiStr := record[metricIndex]
		kpi, _ := decimal.NewFromString(kpiStr)

		newLineCount := historyEvent.Lines + 1
		newKPI := kpi.Add(historyEvent.KPI)
		historyEvent.KPI = newKPI
		historyEvent.Lines = newLineCount

	}

	return historyEvent
}
