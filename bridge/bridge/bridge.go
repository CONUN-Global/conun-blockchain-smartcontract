package bridge

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func Bridge(ctx contractapi.TransactionContextInterface, fcName, toWallet, amount, msg, signature string) (bool, error) {

	params := []string{fcName, toWallet, amount, msg, signature}
	queryArgs := make([][]byte, len(params))
	for i, args := range params {
		queryArgs[i] = []byte(args)
	}

	res := ctx.GetStub().InvokeChaincode("conx", queryArgs, "mychannel")
	if res.Payload == nil {
		return false, fmt.Errorf("error occured while invoking chaincode %s", res.Payload)
	}
	if res.Payload != nil {
		return true, nil
	}

	return false, fmt.Errorf("error while invoking chaincode")
}
