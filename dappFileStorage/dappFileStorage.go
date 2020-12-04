package main

import (
	"dappFileStorage/chaincode"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	dappFileStorage, err := contractapi.NewChaincode(&chaincode.Smartcontract{})
	if err != nil {
		fmt.Println(fmt.Sprintf("Error init SmartContract %s", err))
	}
	if err := dappFileStorage.Start(); err != nil {
		fmt.Println(fmt.Sprintf("Error starting SmartContract FileStorage %s", err))
	}
}
