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

func FilterAndSortByCommitTimestamp[T ObjectWithDate](dataSortableArray []T, startDate time.Time) []T {
	filteredArray := make([]T, 0, len(dataSortableArray))
	for i := range dataSortableArray {
		timestamp := time.Unix(dataSortableArray[i].Timestamp(), 0)
		if timestamp.After(startDate) {
			filteredArray = append(filteredArray, dataSortableArray[i])
		}
	}

	sort.Slice(filteredArray, func(i, j int) bool {
		return filteredArray[i].Timestamp() < filteredArray[j].Timestamp()
	})

	return filteredArray
}

func LegacyGetFirstComputationDateOfPeriod(periodKeyParam common.PeriodKey) (time.Time, error) {
	timegrain, timeGrainError := reports.GetTimeGrain(periodKeyParam)
	periodKey := string(periodKeyParam)
	var firstDate time.Time
	if timeGrainError != nil {
		fmt.Println("Error:", timeGrainError.Error())
		return firstDate, timeGrainError
	}
	switch timegrain {
	case common.Day:
		firstDate, _ = time.Parse("2006-01-02", periodKey)
	case common.Week:
		periodTime, _ := reports.ParseYearWeek(periodKey)
		firstDate = periodTime.AddDate(0, 0, 6).Add(time.Duration(23)*time.Hour + time.Duration(59)*time.Minute + time.Duration(59)*time.Second)
	case common.Month:
		periodTime, _ := time.Parse("2006-01", periodKey)

		firstDate = periodTime.AddDate(0, 1, -1).Add(time.Duration(23)*time.Hour + time.Duration(59)*time.Minute + time.Duration(59)*time.Second)
	case common.Quarter:
		periodTime, _ := reports.ParseQuarterDate(periodKey)

		firstDate = periodTime
	case common.Year:
		periodTime, _ := time.Parse("2006", periodKey)
		firstDate = time.Date(periodTime.Year(), 12, 31, 23, 59, 59, 0, time.UTC)
	default:
		fmt.Printf("Invalid time grain: %s", timegrain)
		return firstDate, fmt.Errorf("invalid time grain: %s", timegrain)
	}
	return firstDate, nil

}

func GetStartDateEndDateAndNextPeriod(periodKey common.PeriodKey) (time.Time, time.Time, common.PeriodKey, error) {
	timegrain, timeGrainError := reports.GetTimeGrain(periodKey)
	if timeGrainError != nil {
		fmt.Println("Error:", timeGrainError.Error())
		return time.Now(), time.Now(), "", timeGrainError
	}
	periodKeyString := string(periodKey)
	switch timegrain {
	case common.Day:
		startDate, err := time.Parse("2006-01-02", periodKeyString)
		nextStartDate := startDate.AddDate(0, 0, 1)
		nextPeriodKey := common.PeriodKey(nextStartDate.Format("2006-01-02"))
		return startDate, nextStartDate, nextPeriodKey, err

	case common.Week:
		startDate, err := reports.GetFirstDateOfYearISOWeek(periodKeyString)
		nextStartDate := startDate.AddDate(0, 0, 7)
		year, week := nextStartDate.ISOWeek()

		nextPeriodKey := common.PeriodKey(fmt.Sprintf("%d-W%02d", year, week))
		return startDate, nextStartDate, nextPeriodKey, err
	case common.Month:
		startDate, err := time.Parse("2006-01", periodKeyString)
		nextStartDate := startDate.AddDate(0, 1, 0)
		nextPeriodKey := common.PeriodKey(nextStartDate.Format("2006-01"))
		return startDate, nextStartDate, nextPeriodKey, err
	case common.Quarter:
		startDate, err := reports.GetFirstDayOfQuarter(periodKeyString)
		nextStartDate := startDate.AddDate(0, 3, 0)
		nextPeriodKey := common.PeriodKey(fmt.Sprintf("%d-Q%d", nextStartDate.Year(), (nextStartDate.Month()-1)/3+1))
		return startDate, nextStartDate, nextPeriodKey, err
	case common.Year:
		startDate, err := time.Parse("2006", periodKeyString)
		nextStartDate := startDate.AddDate(1, 0, 0)
		nextPeriodKey := common.PeriodKey(nextStartDate.Format("2006"))
		return startDate, nextStartDate, nextPeriodKey, err

	default:
		fmt.Printf("Invalid time grain: %s", timegrain)
		return time.Now(), time.Now(), periodKey, fmt.Errorf("invalid time grain: %s", timegrain)
	}
}
