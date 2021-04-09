package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/drive/base"
	Crypto "github.com/drive/crypto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// initliaze contract
type SmartContract struct {
	contractapi.Contract
}

const allowancePrefix = "allowance~ccid~user"
const likePrefix = "like~ccid~user~txId"
const dislikePrefix = "dislike~ccid~user~txId"
const downloadCount = "ccid~user~txId"

/**
Create Content
@param string ccid
@param string author
@param string approveTo
@returns
@memeberof Drive
*/
func (s *SmartContract) CreateFile(ctx contractapi.TransactionContextInterface, author, ipfsHash, data string) (interface{}, error) {

	var cd base.Content
	var err error
	// check for file existance
	if err = json.Unmarshal([]byte(data), &cd); err != nil {
		return nil, fmt.Errorf(base.JSONParseError)
	}
	hashSha1 := Crypto.EncodeToSha256(ipfsHash)
	if exists, err := s.FileExists(ctx, hashSha1); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("%s %s", base.FileExistsError, hashSha1)
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	err = ctx.GetStub().PutState(ipfsHash, []byte(author))
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}
	err = ctx.GetStub().PutState(hashSha1, []byte(ipfsHash))
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}

	details := &base.DetailsTx{
		From:   author,
		To:     "Drive",
		Action: "Create",
		Value:  hashSha1,
	}

	dtl, err := json.Marshal(details)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), dtl)
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}
	res := &base.Response{
		Success:   true,
		Fcn:       "CreateFile",
		TxID:      ctx.GetStub().GetTxID(),
		Value:     hashSha1,
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf(base.JSONParseError)
	}

	return string(content), nil

}

/**
Approve Content
@param string ccid
@param string author
@param string approveTo
@returns
@memeberof Drive
*/
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, ccidcode, author, spenderAdr string) (interface{}, error) {
	ccidByte, err := ctx.GetStub().GetState(ccidcode)
	if err != nil {
		return nil, fmt.Errorf(base.GetstateError)
	}
	ownerByte, err := ctx.GetStub().GetState(string(ccidByte))
	if err != nil {
		return nil, fmt.Errorf(base.GetstateError)
	}
	if string(ownerByte) != author {
		return nil, fmt.Errorf("%s: %s, %s", base.OwnerError, string(ownerByte), author)
	}
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{ccidcode, spenderAdr})
	if err != nil {
		return nil, fmt.Errorf(base.KeyCreationError)
	}
	err = ctx.GetStub().PutState(allowanceKey, []byte(spenderAdr))
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	res := &base.Response{
		Success:   true,
		Fcn:       "Approve",
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf(base.JSONParseError)
	}
	details := &base.DetailsTx{
		From:   author,
		To:     "Drive",
		Action: "Approve",
		Value:  fmt.Sprintf("%s to %s ", ccidcode, spenderAdr),
	}

	dtl, err := json.Marshal(details)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), dtl)
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}
	return string(content), nil

}

/**
Like Content Counter
@param string ccid
@param string walletid
@param []int args [contentID, userID]
@returns
@memeberof Drive
*/
func (s *SmartContract) Allowance(ctx contractapi.TransactionContextInterface, ccidcode, spender string) (bool, error) {

	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{ccidcode, spender})
	if err != nil {
		return false, fmt.Errorf(base.KeyCreationError)
	}
	allowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return false, fmt.Errorf(base.GetstateError)
	}

	if allowanceBytes == nil {
		return false, fmt.Errorf(base.EmptyAllowance)
	}

	return true, nil
}

/**
Like Content Counter
@param string ccid
@param string walletid
@returns
@memeberof Drive
*/
func (s *SmartContract) LikeContent(ctx contractapi.TransactionContextInterface, ccid, walletid string, args []string) (interface{}, error) {

	exists, err := s.FileExists(ctx, ccid)
	if err != nil {
		return nil, fmt.Errorf(base.CheckFileError)
	} else if !exists {
		return nil, fmt.Errorf(base.EmptyFile)
	}
	txID := ctx.GetStub().GetTxID()
	contentLikeKey, err := ctx.GetStub().CreateCompositeKey(likePrefix, []string{ccid, walletid, txID})
	if err != nil {
		return nil, fmt.Errorf(base.KeyCreationError)
	}
	err = ctx.GetStub().PutState(contentLikeKey, []byte{0x00})
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	res := &base.Response{
		Success:   true,
		Fcn:       "LikeContent",
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	details := &base.DetailsTx{
		From:   walletid,
		To:     "Drive",
		Action: "Like",
		Value:  ccid,
	}
	// set event
	likeEvent := &base.Event{UserID: args[0], ContentID: args[1], Timestamp: txTime}
	likeEventJSON, err := json.Marshal(likeEvent)
	if err != nil {
		return nil, fmt.Errorf(base.JSONParseError)
	}
	err = ctx.GetStub().SetEvent("UserLikes", likeEventJSON)
	if err != nil {
		return nil, fmt.Errorf(base.EventError)
	}

	dtl, err := json.Marshal(details)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), dtl)
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}
	return string(content), nil

}

/**
Download Content Counter
@param string ccid
@param string walletid
@param []int args [contentID, userID]
@returns
@memeberof Drive
*/
func (s *SmartContract) CountDownloads(ctx contractapi.TransactionContextInterface, ccid, walletid string, args []string) (interface{}, error) {

	// check for file existance
	exists, err := s.FileExists(ctx, ccid)
	if err != nil {
		return nil, fmt.Errorf(base.CheckFileError)
	} else if !exists {
		return nil, fmt.Errorf(base.EmptyFile)
	}
	// get txID
	txID := ctx.GetStub().GetTxID()
	downloadCount, err := ctx.GetStub().CreateCompositeKey(downloadCount, []string{ccid, walletid, txID})
	if err != nil {
		return nil, fmt.Errorf(base.KeyCreationError)
	}

	err = ctx.GetStub().PutState(downloadCount, []byte{0x00})
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}
	// set response
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	res := &base.Response{
		Success:   true,
		Fcn:       "CountDownloads",
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf(base.JSONParseError)
	}
	details := &base.DetailsTx{
		From:   walletid,
		To:     "Drive",
		Action: "Download",
		Value:  ccid,
	}

	// set event
	downloadEvent := &base.Event{UserID: args[0], ContentID: args[1], Timestamp: txTime}
	downloadEventJSON, err := json.Marshal(downloadEvent)
	if err != nil {
		return nil, fmt.Errorf(base.JSONParseError)
	}
	// emit event
	err = ctx.GetStub().SetEvent("UserDownloads", downloadEventJSON)
	if err != nil {
		return nil, fmt.Errorf(base.EventError)
	}
	// marshal json data
	dtl, err := json.Marshal(details)
	err = ctx.GetStub().PutState(ctx.GetStub().GetTxID(), dtl)
	if err != nil {
		return nil, fmt.Errorf(base.PutStateError)
	}
	return string(content), nil
}

/* check file exists
 *this function strictly called inside chaincode
 */
func (s *SmartContract) FileExists(ctx contractapi.TransactionContextInterface, ccid string) (bool, error) {
	ipfsHash, err := ctx.GetStub().GetState(ccid)

	if err != nil {
		return false, fmt.Errorf("failde to read from world state: %v", err)
	}

	return ipfsHash != nil, nil
}

func (s *SmartContract) GetFile(ctx contractapi.TransactionContextInterface, ccid, spender string) (interface{}, error) {
	if exists, err := s.FileExists(ctx, ccid); err != nil {
		return nil, fmt.Errorf("error checking File, %s", err)
	} else if !exists {
		return nil, fmt.Errorf("error getting file doesnt exists")
	}

	//check allowance
	if ok, _ := s.Allowance(ctx, ccid, spender); ok {
		ipfsHash, err := ctx.GetStub().GetState(ccid)
		if err != nil {
			return nil, err
		}
		if ipfsHash != nil {
			res := &base.Response{
				Success: true,
				Fcn:     "GetFile",
				Value:   string(ipfsHash),
			}
			content, err := json.Marshal(res)
			if err != nil {
				return nil, err
			}

			return string(content), nil
		}
		return nil, fmt.Errorf("Ipfs hash is empty")
	}
	return nil, fmt.Errorf("You do not have allowance for this file")
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
	res := &base.Response{
		Success: true,
		Fcn:     "GetTotalLikes",
		Value:   strconv.Itoa(finalVal),
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

	res := &base.Response{
		Success: true,
		Fcn:     "GetTotalDownloads",
		Value:   strconv.Itoa(finalVal),
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return string(content), nil
}
