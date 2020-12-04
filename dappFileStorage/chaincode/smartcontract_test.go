package chaincode_test

import (
	"encoding/json"
	"testing"

	"chaincode/mocks"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/require"
)

// go generate counterfeither -o mocks/transaction.go -fake-name TransactionContext . transactionContext

type transactionContext interface {
	contractapi.TransactionContextInterface
}

// go generate counterfeither -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

// go generate counterfeither -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

// test ============

func TestCreatFile(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	createFile := dappFileStorage.SmartContract{}
	err := createFile.CreateFile(transactionContext, "someId1", "", "", 0, 0, 0, 0, "")
	require.NoError(t, err)

	chaincodeStub.GetStubReturns([]byte{}, nil)
	err = createFile.CreateFile(transactionContext, "someId1", "", "", 0, 0, 0, 0, "")
	require.EqualError(t, err, "the file with with someId1 is already exists")

}

func TestUpdateFileProgress(t *Testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedFileProgress := &chaincode.FileData{ID: "someId1"}
	bytes, err := json.Marshal(expectedFileProgress)
	require.noError(t, err)

	chaincodeStub.GetStubReturns(bytes, nil)
	fileContract := chaincode.SmartContract{}
	err = fileContract.UpdateFileProgress(transactionContext, "", "", 0)
	require.NoError(t, err)
}
