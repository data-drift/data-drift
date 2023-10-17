package local_store

import (
	"bufio"
	"encoding/csv"
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

func getMetricHistory(store string, table string, metricName string, periodKey common.PeriodKey) ([]common.MetricMeasurement, error) {
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

	var history []common.MetricMeasurement

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

		commitComments := []common.CommitComments{
			{
				CommentAuthor: commit.Author.Name,
				CommentBody:   commit.Message,
			},
		}
		metricEvent := computeMetricHistoryEvent(records, metricName, periodKey, time.Unix(commit.Author.When.Unix(), 0), commitComments)

		if len(history) == 0 {
			history = append(history, metricEvent)
		} else if !history[len(history)-1].Metric.Equal(metricEvent.Metric) {
			history = append(history, metricEvent)
		} else if history[len(history)-1].Metric.Equal(metricEvent.Metric) {
			history[len(history)-1] = metricEvent
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

func computeMetricHistoryEvent(records [][]string, metricName string, periodKey common.PeriodKey, measureDate time.Time, commitComments []common.CommitComments) common.MetricMeasurement {
	headers := records[0]
	metricIndex := findMetricIndex(headers, metricName)
	dateIndex := findMetricIndex(headers, "date")

	firstDateOfPeriod, firstDateOfNextPeriod, _, _ := reducers.GetStartDateEndDateAndNextPeriod(periodKey)

	isAfterPeriod := measureDate.After(firstDateOfNextPeriod) || measureDate.Equal(firstDateOfNextPeriod)

	measurementMetaData := common.MeasurementMetaData{
		IsMeasureAfterPeriod: isAfterPeriod,
		MeasurementTimestamp: measureDate.Unix(),
		MeasurementDate:      measureDate.Format("2006-01-02"),
		MeasurementDateTime:  measureDate.Format("2006-01-02 15:04:05"),
		MeasurementComments:  commitComments,
	}

	var historyEvent = common.MetricMeasurement{
		MeasurementMetaData: measurementMetaData,
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

		newLineCount := historyEvent.LineCount + 1
		newKPI := kpi.Add(historyEvent.Metric)
		historyEvent.Metric = newKPI
		historyEvent.LineCount = newLineCount

	}

	return historyEvent
}
