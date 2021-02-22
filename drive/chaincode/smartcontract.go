package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// initliaze contract
type SmartContract struct {
	contractapi.Contract
}

// initlialize file stuct
// type FileData struct {
// 	ID        string   `json:"ID"`        // id of the file
// 	Author    string   `json:"Author"`    // requester wallet id
// 	Path      string   `json:"Path"`      // ipfs file path
// 	State     int      `json:"State"`     // state of the file
// 	Period    int      `json:"Period"`    // deployment period of the file
// 	RentPrice int      `json:"RentPrice"` // price of the space for the file "paid to provider"
// 	Price     int      `json:"Price"`     // price of the file for clients "paid to author"
// 	Provider  string   `json:"Provider"`  // provider wallet address
// 	Clients   []string `json:"Clients"`   // clients wallet addresses
// 	TxID      string   `json:"TxID"`      // transcation id
// 	Timestamp string   `json:"Timestamp"` // timestamp of transaction
// }

const allowancePrefix = "allowance"
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
	Fcn       string `json:"Fcn"`       // function name
	ID        string `json:"ID"`        // file id
	Success   bool   `json:"Success"`   // true if success
	TxID      string `json:"TxID"`      // transction id
	Timestamp string `json:"Timestamp"` // timestamp of the transaction
}

/**
 * function: CreateFile
 *
 * @param {Context} ctx the transaction context
 * @param {String} id the string id  of the file
 * @param {String} author The author of the file aka 'Requestor wallet.id'
 */
func (s *SmartContract) CreateFile(ctx contractapi.TransactionContextInterface, ccid, author string) error {
	// check for file existance
	if exists, err := s.FileExists(ctx, ccid); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("The file %s already exists", ccid)
	}

	return ctx.GetStub().PutState(ccid, []byte(author))

}

func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, ccid, author, spenderAdr string) error {
	ownerByte, err := ctx.GetStub().GetState(ccid)
	if err != nil {
		return err
	}
	if string(ownerByte) != author {
		return fmt.Errorf("owner are wrong address")
	}
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{ccid, spenderAdr})
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix")
	}
	err = ctx.GetStub().PutState(allowanceKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("error to update state of the smart contract")
	}
	return nil

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

	return txID, nil

}

func (s *SmartContract) DislikeContent(ctx contractapi.TransactionContextInterface, ccid, walletid string) (interface{}, error) {

	exists, err := s.FileExists(ctx, ccid)
	if err != nil {
		return nil, fmt.Errorf("error checking file")
	} else if !exists {
		return nil, fmt.Errorf("error getting file doesnt exists")
	}
	txID := ctx.GetStub().GetTxID()
	contentDislikeKey, err := ctx.GetStub().CreateCompositeKey(dislikePrefix, []string{ccid, walletid, txID})
	if err != nil {
		return nil, fmt.Errorf("getting error while creating like key")
	}
	err = ctx.GetStub().PutState(contentDislikeKey, []byte{0x00})
	if err != nil {
		return nil, fmt.Errorf("error while writing data")
	}

	return txID, nil
}

func (s *SmartContract) DownloadCount(ctx contractapi.TransactionContextInterface, ccid, walletid string) (interface{}, error) {

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

	return txID, nil
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

func (s *SmartContract) getTotatLikes(ctx contractapi.TransactionContextInterface, ccid string) (interface{}, error) {

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

	return finalVal, nil
}

func (s *SmartContract) getTotalDislikes(ctx contractapi.TransactionContextInterface, ccid string) (interface{}, error) {

	//get all deltas for the variable

	deltaResultIterator, deltaErr := ctx.GetStub().GetStateByPartialCompositeKey(dislikePrefix, []string{ccid})
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

	return finalVal, nil
}

func (s *SmartContract) getTotalDownloads(ctx contractapi.TransactionContextInterface, ccid string) (interface{}, error) {

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

	return finalVal, nil
}
