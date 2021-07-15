package chaincode

import (
	"encoding/json"

	"github.com/bridge/bridge"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Details struct {
	Id     string `json:"id"`
	User   string `json:"user"`
	Amount string `json:"amount"`
}

type TxDetails struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Action string `json:"action"`
	Value  string `json:"value"`
}

const DepositPrefix = "depostix~prefix"
const WithdrawPrefix = "withdraw~prefix"

const TokenContract = "token"

// deposit
func (s *SmartContract) MintAndTransfer(ctx contractapi.TransactionContextInterface, data string) (interface{}, error) {

	var dataJson Details

	err := json.Unmarshal([]byte(data), &dataJson)
	if err != nil {
		return nil, err
	}

	_, err = bridge.Bridge(ctx, "MintAndTransfer", dataJson.User, dataJson.Amount)
	if err != nil {
		return nil, err
	}

	response := &TxDetails{
		From:   "Bridge",
		To:     dataJson.User,
		Action: "Mint",
		Value:  dataJson.Amount,
	}

	resp, err := json.Marshal(response)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), resp)

	// call the conos contract token

	return string(resp), nil
}

func (s *SmartContract) BurnFrom(ctx contractapi.TransactionContextInterface, data string) (interface{}, error) {
	var dataJson Details

	err := json.Unmarshal([]byte(data), &dataJson)
	if err != nil {
		return nil, err
	}

	_, err = bridge.Bridge(ctx, "BurnFrom", dataJson.User, dataJson.Amount)
	if err != nil {
		return nil, err
	}

	response := &TxDetails{
		From:   dataJson.User,
		To:     "0x0",
		Action: "BurnFrom",
		Value:  dataJson.Amount,
	}

	resp, err := json.Marshal(response)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), resp)

	return string(resp), nil
}
