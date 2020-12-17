package chaincode

import (
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// action struct
type Action struct {
	ActionId   int    `json:"actionId"`
	ActionCode int    `json:"actionCode"`
	CcId       string `json:"ccId"`
	Timestamp  string `json:"timestamp"`
	TxId       string `json:"TxId"`
}

//user struct
type User struct {
	User    string   `json:"user"`
	Actions []Action `json:"actions"`
}

// SmartContract provides functions for transferring tokens between accounts
type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) ActionWrite(ctx contractapi.TransactionContextInterface, user string, actionId, actionCode int, ccid, txId string) error {

	action := &Action{
		ActionId:   actionId,
		ActionCode: actionCode,
		CcId:       ccid,
		TxId:       txId,
	}
	content, err := json.Marshal(action)
	if err != nil {
		return err
	}
	if err = ctx.GetStub().PutState(strconv.Itoa(actionId), content); err != nil {
		return err
	}
	userByte, err := ctx.GetStub().GetState(user)
	if err != nil {
		return err
	}

	if userByte == nil {
		err = ctx.GetStub().PutState(user, []byte(""))
	}

	return nil
}
