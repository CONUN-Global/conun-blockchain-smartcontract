package util

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
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

// [ParsePositive string its positive?]
func ParsePositive(s string) (decimal.Decimal, error) {
	var d decimal.Decimal
	var err error

	if d, err = decimal.NewFromString(s); err != nil {
		return d, fmt.Errorf("is not integer %s", err)
	}
	if !isNumeric(s) {
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
	if !isNumeric(s) {
		return d, fmt.Errorf("%s is not integer", s)
	}
	if d.IsNegative() {
		return d, fmt.Errorf("%s is negative", s)
	}
	return d, nil
}

func VerifyMsgAddr(from, sign, msg string) (bool, error) {

	msgBytes, err := mustDecodeUtil(msg)
	if err != nil {
		return false, err
	}

	sig, err := mustDecodeUtil(sign)

	if err != nil {
		return false, err
	}
	if sig[64] != 27 && sig[64] != 28 {
		return false, fmt.Errorf("error signature is not valid type")
	}
	sig[64] -= 27
	sigPubKey, err := crypto.Ecrecover(msgBytes, sig)
	if err != nil {
		return false, fmt.Errorf("error verifying msg %s", err)
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write(sigPubKey[1:]) // 0x
	if strings.Compare(strings.ToLower(string(hexutil.Encode(hash.Sum(nil)[12:]))), strings.ToLower(from)) == 0 {
		return true, nil
	}
	return false, fmt.Errorf("error address doesnt mastch")
}

// MustDecode decodes a hex string with 0x prefix. It panics for invalid input.
func mustDecodeUtil(input string) ([]byte, error) {
	dec, err := hexutil.Decode(input)
	return dec, err
}
