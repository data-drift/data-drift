package helpers

import "github.com/shopspring/decimal"

func GetFloat(dec decimal.Decimal) float64 {
	f, _ := dec.Float64()
	return f
}
