package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bridge/base"
	"github.com/bridge/bridge"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Details struct {
	Id        string `json:"id"`
	Key       string `json:"key"`
	User      string `json:"user"`
	Amount    string `json:"amount"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type TxDetails struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Action string `json:"action"`
	Value  string `json:"value"`
}

type Event struct {
	Id     string `json:"id"`
	User   string `json:"user"`
	Amount string `json:"amount"`
}

const DepositPrefix = "depostix~prefix"
const WithdrawPrefix = "withdraw~prefix"
const TokenContract = "token"

var IdState = make(map[string]bool)

// deposit
func (s *SmartContract) MintAndTransfer(ctx contractapi.TransactionContextInterface, data string) (interface{}, error) {

	var dataJson Details

	err := json.Unmarshal([]byte(data), &dataJson)
	if err != nil {
		return nil, err
	}
	if _, exists := IdState[dataJson.Id]; exists {
		return nil, fmt.Errorf("key Id is already exists")
	}

	hash := sha256.New()
	hash.Write([]byte(dataJson.Key))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	if c := strings.Compare(dataJson.Id, mdStr); c < 0 {
		return nil, fmt.Errorf("keys are not matching")
	}

	_, err = bridge.Bridge(ctx, "MintAndTransfer", dataJson.User, dataJson.Amount, dataJson.Message, dataJson.Signature)
	if err != nil {
		return nil, err
	}

	IdState[dataJson.Id] = true

	response := &TxDetails{
		From:   "Bridge",
		To:     dataJson.User,
		Action: "Mint",
		Value:  dataJson.Amount,
	}

	// set event
	mintEevent := &Event{Id: dataJson.Id, User: dataJson.User, Amount: dataJson.Amount}
	mintEeventJSON, err := json.Marshal(mintEevent)
	if err != nil {
		return nil, fmt.Errorf(base.JSONParseError)
	}
	err = ctx.GetStub().SetEvent("MintAndTransfer", mintEeventJSON)
	if err != nil {
		return nil, fmt.Errorf(base.EventError)
	}

	resp, _ := json.Marshal(response)
	_ = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), resp)

	return string(resp), nil
}

func (s *SmartContract) BurnFrom(ctx contractapi.TransactionContextInterface, data string) (interface{}, error) {
	var dataJson Details

	err := json.Unmarshal([]byte(data), &dataJson)
	if err != nil {
		return nil, err
	}

	if _, exists := IdState[dataJson.Id]; exists {
		return nil, fmt.Errorf("key Id is already exists")
	}

	_, err = bridge.Bridge(ctx, "BurnFrom", dataJson.User, dataJson.Amount, dataJson.Message, dataJson.Signature)
	if err != nil {
		return nil, err
	}

	response := &TxDetails{
		From:   dataJson.User,
		To:     "0x0",
		Action: "BurnFrom",
		Value:  dataJson.Amount,
	}

	// set event
	burnEevent := &Event{Id: dataJson.Id, User: dataJson.User, Amount: dataJson.Amount}
	burnEeventtJSON, err := json.Marshal(burnEevent)
	if err != nil {
		return nil, fmt.Errorf(base.JSONParseError)
	}
	err = ctx.GetStub().SetEvent("BurnFrom", burnEeventtJSON)
	if err != nil {
		return nil, fmt.Errorf(base.EventError)
	}

	resp, _ := json.Marshal(response)
	_ = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), resp)

	return string(resp), nil
}
