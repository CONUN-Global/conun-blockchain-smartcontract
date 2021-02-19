package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type SmartContract struct {
}

const (
	OK    = 200
	ERROR = 500
)

/*
	Smart cpm
*/

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// func (s *SmartContract) Invoke(ApIstub shim.ChaincodeStubInterface) pb.Response {

// 	function, args := APIstub.GetFunctionAndParametrs()

// 	if function == "update" {
// 		return "gooo"
// 	}
// }

func (s *SmartContract) update(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("incorrect number of arguments")
	}
	name := args[0]
	op := args[2]
	_, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("you got error")
	}

	txid := APIstub.GetTxID()
	compositeIndexName := "varName~op~value~TxID"

	compositeKey, compositeErr := APIstub.CreateCompositeKey(compositeIndexName, []string{name, op, args[1], txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("could not create a composite key for"))
	}

	compositePutErr := APIstub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("coould not put operation"))
	}
	return shim.Success([]byte(fmt.Sprintf("succesfully added ")))
}
