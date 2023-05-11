package charts

import (
	"testing"
	"time"
)

func TestGetFirstDateOfPeriod_Day(t *testing.T) {
	expected := time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)
	result := getFirstDateOfPeriod("2023-05-01")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFirstDateOfPeriod_Week(t *testing.T) {
	expected := time.Date(2023, time.April, 30, 23, 59, 59, 0, time.UTC)
	result := getFirstDateOfPeriod("2023-W17")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFirstDateOfPeriod_Month(t *testing.T) {
	expected := time.Date(2023, time.May, 31, 23, 59, 59, 0, time.UTC)
	result := getFirstDateOfPeriod("2023-05")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFirstDateOfPeriod_Quarter(t *testing.T) {
	expected := time.Date(2023, time.March, 31, 0, 0, 0, 0, time.UTC)
	result := getFirstDateOfPeriod("2023-Q1")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetFirstDateOfPeriod_Year(t *testing.T) {
	expected := time.Date(2023, time.December, 31, 23, 59, 59, 0, time.UTC)
	result := getFirstDateOfPeriod("2023")
	if !result.Equal(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
