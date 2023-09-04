package reducers

import (
	"fmt"
	"sort"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/reports"
)

type ObjectWithDate interface {
	Timestamp() int64
}

func FilterAndSortByCommitTimestamp[T ObjectWithDate](dataSortableArray []T, driftDay time.Time) []T {
	filteredArray := make([]T, 0, len(dataSortableArray))
	for i := range dataSortableArray {
		timestamp := time.Unix(dataSortableArray[i].Timestamp(), 0)
		if timestamp.After(driftDay) {
			filteredArray = append(filteredArray, dataSortableArray[i])
		}
	}

	sort.Slice(filteredArray, func(i, j int) bool {
		return filteredArray[i].Timestamp() < filteredArray[j].Timestamp()
	})

	return filteredArray
}

func GetFirstDateOfPeriod(periodKeyParam common.PeriodKey) (time.Time, error) {
	timegrain, timeGrainError := reports.GetTimeGrain(periodKeyParam)
	periodKey := string(periodKeyParam)
	var lastDay time.Time
	if timeGrainError != nil {
		fmt.Println("Error:", timeGrainError.Error())
		return lastDay, timeGrainError
	}
	switch timegrain {
	case common.Day:
		lastDay, _ = time.Parse("2006-01-02", periodKey)
	case common.Week:
		periodTime, _ := reports.ParseYearWeek(periodKey)
		lastDay = periodTime.AddDate(0, 0, 6).Add(time.Duration(23)*time.Hour + time.Duration(59)*time.Minute + time.Duration(59)*time.Second)
	case common.Month:
		periodTime, _ := time.Parse("2006-01", periodKey)

		lastDay = periodTime.AddDate(0, 1, -1).Add(time.Duration(23)*time.Hour + time.Duration(59)*time.Minute + time.Duration(59)*time.Second)
	case common.Quarter:
		periodTime, _ := reports.ParseQuarterDate(periodKey)

		lastDay = periodTime
	case common.Year:
		periodTime, _ := time.Parse("2006", periodKey)
		lastDay = time.Date(periodTime.Year(), 12, 31, 23, 59, 59, 0, time.UTC)
	default:
		fmt.Printf("Invalid time grain: %s", timegrain)
		return lastDay, fmt.Errorf("invalid time grain: %s", timegrain)
	}
	return lastDay, nil

}
