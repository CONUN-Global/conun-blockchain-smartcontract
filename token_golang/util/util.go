package util

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func isNumeric(s string) bool {
	for _, v := range s {
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

// [ParsePositive string its positive?]
func ParsePositive(s string) (decimal.Decimal, error) {
	var d decimal.Decimal
	var err error

	if d, err = decimal.NewFromString(s); err != nil {
		return d, fmt.Errorf("is not integer %s", err)
	}
	if isNumeric(s) == false {
		return d, fmt.Errorf("is not integer string %s", err)
	}
	if !d.IsPositive() {
		return d, fmt.Errorf(" is either 0 or negative  %s", err)
	}
	return d, nil
}

func ParseNotNegative(s string) (decimal.Decimal, error) {
	var d decimal.Decimal
	var err error

	if d, err = decimal.NewFromString(s); err != nil {
		return d, fmt.Errorf("%s is not integer string", s)

	}
	if isNumeric(s) == false {
		return d, fmt.Errorf("%s is not integer", s)
	}
	if d.IsNegative() {
		return d, fmt.Errorf("%s is negative", s)
	}
	return d, nil
}
