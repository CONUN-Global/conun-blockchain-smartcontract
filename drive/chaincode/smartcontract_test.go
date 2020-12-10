package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/drive/chaincode"
	"github.com/drive/chaincode/mocks"

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

	createFile := chaincode.SmartContract{}
	err := createFile.CreateFile(transactionContext, "someId1", "aziz", "", 1, 1, 1, 1, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = createFile.CreateFile(transactionContext, "someId1", "", "", 0, 0, 0, 0, "")
	require.EqualError(t, err, "the file with with someId1 is already exists")

}
func TestOrderFileFromAuthor(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	orderFileFromAuthor := chaincode.SmartContract{}
	err := orderFileFromAuthor.OrderFileFromAuthor(transactionContext, "someId", "Sara")
	//require.NoError(t, err, resp)
	fmt.Println(err)
}

func TestUpdateFileProgress(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedFileProgress := &chaincode.FileData{ID: "someId1"}
	bytes, err := json.Marshal(expectedFileProgress)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	fileContract := chaincode.SmartContract{}
	check, err := fileContract.UpdateFileProgress(transactionContext, "", "", 0)
	require.NoError(t, err, check)
}
