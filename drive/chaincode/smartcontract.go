package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// initliaze contract
type SmartContract struct {
	contractapi.Contract
}

//initlialize file stuct
type FileData struct {
	Author    string `json:"Author"`    // requester wallet id
	State     int    `json:"State"`     // state of the file
	TxID      string `json:"TxID"`      // transcation id
	Timestamp string `json:"Timestamp"` // timestamp of transaction
}

const allowancePrefix = "allowance~ccid~user"
const likePrefix = "like~ccid~user~txId"
const dislikePrefix = "dislike~ccid~user~txId"
const downloadCount = "ccid~user~txId"

type OrderFile struct {
	ID     string `json:"ID"`
	Author string `json:"Author"`
	Path   string `json:"Path"`
	Price  int    `json:"Price"`
}

// initialize response
type Response struct {
	Fcn       string               `json:"Fcn,omitempty"`       // function name
	Success   bool                 `json:"Success,omitempty"`   // true if success
	TxID      string               `json:"TxID,omitempty"`      // transction id
	Timestamp *timestamp.Timestamp `json:"Timestamp,omitempty"` // timestamp of the transaction
	Value     int                  `json:"Value,omitempty"`     // value of dislike of like, count
}

// txDetails struct
// Tx Details struct
type DetailsTx struct {
	From   string `json:"From"`
	To     string `json:"To"`
	Action string `json:"Action"`
	Value  string `json:"Value"`
}

/**
 * function: CreateFile
 *
 * @param {Context} ctx the transaction context
 * @param {String} id the string id  of the file
 * @param {String} author The author of the file aka 'Requestor wallet.id'
 */
func (s *SmartContract) CreateFile(ctx contractapi.TransactionContextInterface, ccid, author string) (interface{}, error) {
	// check for file existance
	if exists, err := s.FileExists(ctx, ccid); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("The file %s already exists", ccid)
	}
	err := ctx.GetStub().PutState(ccid, []byte(author))
	if err != nil {
		return nil, err
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	res := &Response{
		Success:   true,
		Fcn:       "CreateFile",
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	details := &DetailsTx{
		From:   author,
		To:     "Drive",
		Action: "Create",
		Value:  ccid,
	}

	dtl, err := json.Marshal(details)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), dtl)
	if err != nil {
		return nil, err
	}
	return string(content), nil

}

func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, ccid, author, spenderAdr string) (interface{}, error) {
	ownerByte, err := ctx.GetStub().GetState(ccid)
	if err != nil {
		return nil, err
	}
	if string(ownerByte) != author {
		return nil, fmt.Errorf("owner are wrong address %s, %s", string(ownerByte), author)
	}
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{ccid, spenderAdr})
	if err != nil {
		return nil, fmt.Errorf("failed to create the composite key for prefix")
	}
	err = ctx.GetStub().PutState(allowanceKey, []byte(spenderAdr))
	if err != nil {
		return nil, fmt.Errorf("error to update state of the smart contract")
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	res := &Response{
		Success:   true,
		Fcn:       "Approve",
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	details := &DetailsTx{
		From:   author,
		To:     "Drive",
		Action: "Approve",
		Value:  spenderAdr,
	}

	dtl, err := json.Marshal(details)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), dtl)
	if err != nil {
		return nil, err
	}
	return string(content), nil

}

func (s *SmartContract) Allowance(ctx contractapi.TransactionContextInterface, ccid, spender string) (interface{}, error) {

	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{ccid, spender})
	if err != nil {
		return nil, fmt.Errorf("error creating composite key")
	}
	allowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read allowance for ")
	}

	if allowanceBytes == nil {
		return nil, fmt.Errorf("allowance is empty")
	}
	return string(allowanceBytes), nil
}

func (s *SmartContract) LikeContent(ctx contractapi.TransactionContextInterface, ccid, walletid string) (interface{}, error) {

	exists, err := s.FileExists(ctx, ccid)
	if err != nil {
		return nil, fmt.Errorf("error checking file")
	} else if !exists {
		return nil, fmt.Errorf("error getting file doesnt exists")
	}
	txID := ctx.GetStub().GetTxID()
	contentLikeKey, err := ctx.GetStub().CreateCompositeKey(likePrefix, []string{ccid, walletid, txID})
	if err != nil {
		return nil, fmt.Errorf("getting error while creating like key")
	}
	err = ctx.GetStub().PutState(contentLikeKey, []byte{0x00})
	if err != nil {
		return nil, fmt.Errorf("error while writing data")
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	res := &Response{
		Success:   true,
		Fcn:       "LikeContent",
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	details := &DetailsTx{
		From:   walletid,
		To:     "Drive",
		Action: "Like",
		Value:  ccid,
	}

	dtl, err := json.Marshal(details)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), dtl)
	if err != nil {
		return nil, err
	}
	return string(content), nil

}

func (s *SmartContract) CountDownloads(ctx contractapi.TransactionContextInterface, ccid, walletid string) (interface{}, error) {

	exists, err := s.FileExists(ctx, ccid)
	if err != nil {
		return nil, fmt.Errorf("error checking file")
	} else if !exists {
		return nil, fmt.Errorf("error getting file doesnt exists")
	}
	txID := ctx.GetStub().GetTxID()
	downloadCount, err := ctx.GetStub().CreateCompositeKey(downloadCount, []string{ccid, walletid, txID})
	if err != nil {
		return nil, fmt.Errorf("getting error while creating like key")
	}
	err = ctx.GetStub().PutState(downloadCount, []byte{0x00})
	if err != nil {
		return nil, fmt.Errorf("error while writing data")
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	res := &Response{
		Success:   true,
		Fcn:       "CountDownloads",
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	details := &DetailsTx{
		From:   walletid,
		To:     "Drive",
		Action: "Download",
		Value:  ccid,
	}

	dtl, err := json.Marshal(details)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), dtl)
	if err != nil {
		return nil, err
	}
	return string(content), nil
}

/* check file exists
 *this function strictly called inside chaincode
 */
func (s *SmartContract) FileExists(ctx contractapi.TransactionContextInterface, ccid string) (bool, error) {
	authorAdr, err := ctx.GetStub().GetState(ccid)

	if err != nil {
		return false, fmt.Errorf("failde to read from world state: %v", err)
	}

	return authorAdr != nil, nil
}

func (s *SmartContract) GetTotalLikes(ctx contractapi.TransactionContextInterface, ccid string) (interface{}, error) {

	//get all deltas for the variable

	deltaResultIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey(likePrefix, []string{ccid})
	if deltaErr != nil {
		return nil, fmt.Errorf("error occured while getting file")
	}

	defer deltaResultIterator.Close()

	if !deltaResultIterator.HasNext() {
		return nil, fmt.Errorf("error getting file empty")
	}

	var finalVal int
	var i int

	for i = 0; deltaResultIterator.HasNext(); i++ {
		//get the next row
		_, nextErr := deltaResultIterator.Next()
		if nextErr != nil {
			return nil, fmt.Errorf(nextErr.Error())
		}

		finalVal += 1
	}
	res := &Response{
		Success: true,
		Fcn:     "GetTotalLikes",
		Value:   finalVal,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return string(content), nil
}

func (s *SmartContract) GetTotalDownloads(ctx contractapi.TransactionContextInterface, ccid string) (interface{}, error) {

	//get all deltas for the variable

	deltaResultIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey(downloadCount, []string{ccid})
	if deltaErr != nil {
		return nil, fmt.Errorf("error occured while getting file")
	}

	defer deltaResultIterator.Close()

	if !deltaResultIterator.HasNext() {
		return nil, fmt.Errorf("error getting file empty")
	}

	var finalVal int
	var i int

	for i = 0; deltaResultIterator.HasNext(); i++ {
		//get the next row
		_, nextErr := deltaResultIterator.Next()
		if nextErr != nil {
			return nil, fmt.Errorf(nextErr.Error())
		}

		finalVal += 1
	}

	res := &Response{
		Success: true,
		Fcn:     "GetTotalDownloads",
		Value:   finalVal,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return string(content), nil
}
