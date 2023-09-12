package helpers

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestFormatWithSeparator(t *testing.T) {
	testCases := []struct {
		input decimal.Decimal
		want  string
	}{
		{decimal.NewFromFloat(1234567.89), "1,234,567.89"},
		{decimal.NewFromFloat(0), "0"},
		{decimal.NewFromFloat(-1234567.89), "-1,234,567.89"},
		{decimal.NewFromFloat(-1234567.8), "-1,234,567.80"},
		{decimal.NewFromFloat(0.123456789), "0.123456789"},
		{decimal.NewFromFloat(0.000000001), "0.000000001"},
	}

	for _, tc := range testCases {
		got := FormatWithSeparator(tc.input)
		if got != tc.want {
			t.Errorf("FormatWithSeparator(%v) = %v; want %v", tc.input, got, tc.want)
		}
	}
}
