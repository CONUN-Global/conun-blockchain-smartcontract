package util

import (
	"fmt"

	"github.com/drive/base"
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

// [ParsePositive string its postive?]
func ParsePositive(s string) (decimal.Decimal, error) {
	var d decimal.Decimal
	var err error

	if d, err = decimal.NewFromString(s); err != nil {
		return d, fmt.Errorf("%s %s", base.NumberError, err)
	}
	if isNumeric(s) == false {
		return d, fmt.Errorf("%s %s", base.NumberError, err)
	}
	if !d.IsPositive() {
		return d, fmt.Errorf("%s %s", base.NumberError, err)
	}
	return d, nil
}
