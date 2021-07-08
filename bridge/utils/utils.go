package utils

import (
	"fmt"

	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/decimal"
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
		return d, fmt.Errorf("Error parsing string: %s", err)
	}

	if isNumeric(s) == false {
		return d, fmt.Errorf("Error string is not valid number")
	}

	if !d.IsPositive() {
		return d, fmt.Errorf("string number is not postive interger")
	}

	return d, nil

}
