package token

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

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
