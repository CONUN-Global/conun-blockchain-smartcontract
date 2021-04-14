package conos

import (
	"fmt"

	"github.com/drive/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// initialize contract
type SmartContract struct {
	contractapi.Contract
}

/**
Send Token
@param string ccid
@param string approveTo
@returns string approveTo
*/
// [SendToken invoke contract conos]
func SendToken(ctx contractapi.TransactionContextInterface, fromWallet, toWallet, amount string) (bool, error) {
	params := []string{"Transfer", fromWallet, toWallet, amount}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}
	res := ctx.GetStub().InvokeChaincode("conos", queryArgs, "mychannel")
	if res.Payload == nil {
		return false, fmt.Errorf("%s %d", base.InvokeChaincodeError, res.Status)
	}
	if res.Payload != nil {
		return true, nil
	}
	return false, fmt.Errorf(base.InvokeChaincodeError)
}
