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
	}{
		{common.PeriodKey("2023-06-01"), common.PeriodKey("2023-06-02")},
		{common.PeriodKey("2023-06-30"), common.PeriodKey("2023-07-01")},
		{common.PeriodKey("2023-12-31"), common.PeriodKey("2024-01-01")},
		{common.PeriodKey("2023-W01"), common.PeriodKey("2023-W02")},
		{common.PeriodKey("2023-W09"), common.PeriodKey("2023-W10")},
		{common.PeriodKey("2023-W52"), common.PeriodKey("2024-W01")},
		// {common.PeriodKey("2024-W52"), common.PeriodKey("2024-W53")}, //@SamoxTODO To be fixed or understood
		// {common.PeriodKey("2024-W53"), common.PeriodKey("2025-W01")},
		{common.PeriodKey("2023-12"), common.PeriodKey("2024-01")},
		{common.PeriodKey("2023-06"), common.PeriodKey("2023-07")},
		{common.PeriodKey("2025-01"), common.PeriodKey("2025-02")},
		{common.PeriodKey("2025-Q1"), common.PeriodKey("2025-Q2")},
		{common.PeriodKey("2025-Q3"), common.PeriodKey("2025-Q4")},
		{common.PeriodKey("2025-Q4"), common.PeriodKey("2026-Q1")},
		{common.PeriodKey("2025"), common.PeriodKey("2026")},
	}

	for _, tc := range testCases {
		actual, err := GetNextPeriod(tc.periodKey)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if actual != tc.expected {
			t.Errorf("GetNextPeriod(%s) = %s, expected %s", tc.periodKey, actual, tc.expected)
		}
	}
}
