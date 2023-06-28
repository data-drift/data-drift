package reports

import (
	"testing"
	"time"
)

func TestParseYearWeek(t *testing.T) {
	yearWeek := "2023-W17" // Monday 24 April to Sunday 30 April 2023
	expectedFirstDay := time.Date(2023, time.April, 24, 0, 0, 0, 0, time.UTC)

	firstDay, err := ParseYearWeek(yearWeek)
	if err != nil {
		t.Errorf("Error parsing year week: %v", err.Error())
	}

	if !firstDay.Equal(expectedFirstDay) {
		t.Errorf("Expected first day of week for %s to be %v, but got %v", yearWeek, expectedFirstDay, firstDay)
	}
}
