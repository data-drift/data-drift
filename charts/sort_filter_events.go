package charts

import (
	"log"
	"sort"
	"time"

	"github.com/data-drift/kpi-git-history/common"
	"github.com/data-drift/kpi-git-history/reports"
)

func FilterAndSortByCommitTimestamp(dataSortableArray []CommitData, driftDay time.Time) []CommitData {
	filteredArray := make([]CommitData, 0, len(dataSortableArray))
	for i := range dataSortableArray {
		timestamp := time.Unix(dataSortableArray[i].CommitTimestamp, 0)
		if timestamp.After(driftDay) {
			filteredArray = append(filteredArray, dataSortableArray[i])
		}
	}

	sort.Slice(filteredArray, func(i, j int) bool {
		return filteredArray[i].CommitTimestamp < filteredArray[j].CommitTimestamp
	})

	return filteredArray
}

func getFirstDateOfPeriod(periodKey string) time.Time {
	timegrain, _ := reports.GetTimeGrain(periodKey)
	var lastDay time.Time
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
		log.Fatalf("Invalid time grain: %s", timegrain)
	}
	return lastDay

}
