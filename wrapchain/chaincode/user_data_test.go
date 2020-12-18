package chaincode_test

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// go generate counterfeither -o mocks/trnasction.go -fake -name Transaction
type transactionContext interface {
	contractapi.SettableTransactionContextInterface
}

type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}
