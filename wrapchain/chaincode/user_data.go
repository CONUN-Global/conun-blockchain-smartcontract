package chaincode

import (
	"encoding/json"
	"fmt"
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
type UserAr struct {
	Actions []*Action `json:"actions"`
}

// SmartContract provides functions for transferring tokens between accounts
type SmartContract struct {
	contractapi.Contract
}

/*
	ActionWrite saves every user action on conun blockchain

	@param {Context} ctx the transaction context
	@param {String} user the user address
	@param {Int} actionId the id of the action that user made
	@param {Int} actionCode the code of the action
	@param {String} ccid function name
	@param {String} the transaction id of that action


	returns Error

*/
func (s *SmartContract) ActionWrite(ctx contractapi.TransactionContextInterface, user string, actionId, actionCode int, ccid, txId string) error {

	if exist, err := s.ActionExists(ctx, strconv.Itoa(actionId)); err != nil {
		return err
	} else if exist {
		return fmt.Errorf("The action already exists: %d", actionId)
	}

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

	if userArjson, err := ctx.GetStub().GetState(user); err != nil {
		return err
	} else if userArjson == nil {
		var actionar []*Action
		actionar = append(actionar, action)
		userAr := &UserAr{
			Actions: actionar,
		}
		contentAr, err := json.Marshal(userAr)
		if err != nil {
			return err
		}
		if err = ctx.GetStub().PutState(user, []byte(contentAr)); err != nil {
			return err
		}
	} else {
		var userAr UserAr
		userArJson, err := ctx.GetStub().GetState(user)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(userArJson, &userAr); err != nil {
			return err
		}

		userAr.Actions = append(userAr.Actions, action)
		contentAr, err := json.Marshal(userAr)
		if err != nil {
			return err
		}
		_ = ctx.GetStub().PutState(user, contentAr)

	}

	return nil
}

/*
	ActionExists checks whether action exists in the blockchain

	@param {String} id the id of the action

	Retunrns Bool or Error
*/
func (s *SmartContract) ActionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	actionArr, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to check action: %s", err)
	}
	return actionArr != nil, nil
}
