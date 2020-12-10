package chaincode_test

import (
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

	chaincodeStub.GetStateReturns(err, nil)
	fileContract := chaincode.SmartContract{}
	check, err := fileContract.UpdateFileProgress(transactionContext, "", "", 0)
	require.NoError(t, err, check)
}

func TestFileExists(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	fileExists := &chaincode.SmartContract{}
	_, err := fileExists.FileExists(transactionContext, "someId1")
	require.NoError(t, err)
}

func TestCancelFileProgress(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	cancelFileProgress := &chaincode.SmartContract{}
	_, err := cancelFileProgress.CancelFileProgress(transactionContext, "someId1,", "", 0)
	require.NoError(t, err)
}
