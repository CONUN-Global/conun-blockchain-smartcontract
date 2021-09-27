package utils

import (
	"fmt"

	"math/big"

	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/decimal"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
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

func GetMsgForSign(_address string, _amount int64) (string, error) {

	uint256Ty := abi.NewType("uint256", "uint256", nil)
	addressTy := abi.NewType("address", "address", nil)

	arguments := abi.Arguments{
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
	}

	bytes, err := arguments.Pack(
		big.NewInt(_amount),
		common.HexToAddress(_address),
	)

	if err != nil {
		return "", err
	}

	var buf []byte
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes)
	buf = hash.Sum(buf)

	return hexutil.Encode(buf), nil

}
