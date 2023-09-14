package reducers

import (
	"testing"
	"time"

	"github.com/data-drift/data-drift/common"
)

func TestGetFirstDateOfPeriod_Day(t *testing.T) {
	expected := time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)
	result, _ := LegacyGetFirstComputationDateOfPeriod("2023-05-01")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFirstDateOfPeriod_Week(t *testing.T) {
	expected := time.Date(2023, time.April, 30, 23, 59, 59, 0, time.UTC)
	result, _ := LegacyGetFirstComputationDateOfPeriod("2023-W17")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFirstDateOfPeriod_Month(t *testing.T) {
	expected := time.Date(2023, time.May, 31, 23, 59, 59, 0, time.UTC)
	result, _ := LegacyGetFirstComputationDateOfPeriod("2023-05")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFirstDateOfPeriod_Quarter(t *testing.T) {
	expected := time.Date(2023, time.March, 31, 0, 0, 0, 0, time.UTC)
	result, _ := LegacyGetFirstComputationDateOfPeriod("2023-Q1")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFirstDateOfPeriod_Year(t *testing.T) {
	expected := time.Date(2023, time.December, 31, 23, 59, 59, 0, time.UTC)
	result, _ := LegacyGetFirstComputationDateOfPeriod("2023")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetNextPeriod(t *testing.T) {
	testCases := []struct {
		periodKey common.PeriodKey
		expected  common.PeriodKey
		start     time.Time
		end       time.Time
	}{
		{common.PeriodKey("2023-06-01"), common.PeriodKey("2023-06-02"), time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 6, 2, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2023-06-30"), common.PeriodKey("2023-07-01"), time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC), time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2023-12-31"), common.PeriodKey("2024-01-01"), time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2023-W01"), common.PeriodKey("2023-W02"), time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2023-W09"), common.PeriodKey("2023-W10"), time.Date(2023, 2, 27, 0, 0, 0, 0, time.UTC), time.Date(2023, 3, 6, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2023-W52"), common.PeriodKey("2024-W01"), time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2024-W53"), common.PeriodKey("2025-W02"), time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC), time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2023-12"), common.PeriodKey("2024-01"), time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2023-06"), common.PeriodKey("2023-07"), time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2025-01"), common.PeriodKey("2025-02"), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2025-Q1"), common.PeriodKey("2025-Q2"), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2025-Q3"), common.PeriodKey("2025-Q4"), time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2025-Q4"), common.PeriodKey("2026-Q1"), time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{common.PeriodKey("2025"), common.PeriodKey("2026"), time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		start, end, next, err := GetStartDateEndDateAndNextPeriod(tc.periodKey)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if next != tc.expected {
			t.Errorf("GetNextPeriod(%s) = %s, expected %s", tc.periodKey, next, tc.expected)
		}
		if !start.Equal(tc.start) {
			t.Errorf("GetNextPeriod(%s) = %s, expected %s", tc.periodKey, start, tc.start)
		}
		if !end.Equal(tc.end) {
			t.Errorf("GetNextPeriod(%s) = %s, expected %s", tc.periodKey, end, tc.end)
		}
	}
}
