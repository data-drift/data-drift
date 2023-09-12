package helpers

import (
	"strings"

	"github.com/shopspring/decimal"
)

func GetFloat(dec decimal.Decimal) float64 {
	f, _ := dec.Float64()
	return f
}

func FormatWithSeparator(d decimal.Decimal) string {
	isNegative := d.IsNegative()
	whole := d.Abs().Floor()
	frac := d.Abs().Sub(whole)

	var sb strings.Builder
	wholeLen := len(whole.String())

	if isNegative {
		sb.WriteRune('-')
	}

	for i, digit := range whole.String() {
		sb.WriteRune(digit)
		if (wholeLen-i-1)%3 == 0 && i < wholeLen-1 {
			sb.WriteRune(',')
		}
	}

	if !frac.IsZero() {
		sb.WriteRune('.')
		stringDecimalPart := frac.String()[2:]
		if len(stringDecimalPart) == 1 {
			stringDecimalPart = stringDecimalPart + "0"
		}
		sb.WriteString(stringDecimalPart)
	}

	return sb.String()
}
