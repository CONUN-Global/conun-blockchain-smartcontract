package utils

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

func ParsePositive(s string) (decimal.Decimal, error) {
	var d decimal.Decimal
	var err error

	if d, err = decimal.NewFromString(s); err != nil {
		return d, fmt.Errorf("error parsing string: %s", err)
	}

	if !isNumeric(s) {
		return d, fmt.Errorf("error string is not valid number")
	}

	if !d.IsPositive() {
		return d, fmt.Errorf("string number is not postive interger")
	}

	return d, nil

}
