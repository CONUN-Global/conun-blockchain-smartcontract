package utils

import (
	"fmt"

	"math/big"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

func GetMsgForSign(_address, swapId string, _amount *big.Int) (string, error) {

	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	addressTy, _ := abi.NewType("address", "address", nil)
	swapID, _ := abi.NewType("bytes32", "bytes32", nil)

	arguments := abi.Arguments{
		{
			Type: swapID,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
	}

	bytes, err := arguments.Pack(
		common.HexToHash("0x"+swapId),
		_amount,
		common.HexToAddress(_address),
	)

	if err != nil {
		return "", err
	}

	var buf []byte
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes)
	buf = hash.Sum(buf)

	haa2 := common.HexToHash(hexutil.Encode(buf))

	prefixedHash := crypto.Keccak256Hash(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%v", len(common.HexToHash(hexutil.Encode(buf))))),
		haa2.Bytes(),
	)
	return prefixedHash.Hex(), nil

}
