package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// initliaze contract
type SmartContract struct {
	contractapi.Contract
}

// initlialize file stuct
type FileData struct {
	ID        string   `json:"ID"`        // id of the file
	Author    string   `json:"Author"`    // requester wallet id
	Path      string   `json:"Path"`      // ipfs file path
	State     int      `json:"State"`     // state of the file
	Period    int      `json:"Period"`    // deployment period of the file
	RentPrice int      `json:"RentPrice"` // price of the space for the file "paid to provider"
	Price     int      `json:"Price"`     // price of the file for clients "paid to author"
	Provider  string   `json:"Provider"`  // provider wallet address
	Clients   []string `json:"Clients"`   // clients wallet addresses
	TxID      string   `json:"TxID"`      // transcation id
	Timestamp string   `json:"Timestamp"` // timestamp of transaction
}

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
 * @param {String} path The path of the ipfs hash that saved on providers pc
 * @param {Int} state The state of the work progress
 * @param {Int} period The period of how long file will be saved on providers pc
 * @param {Int} rentPrice The price of the renting space for file from providers pc
 * @param {Int} price The price of the file
 * @param {String} provider The provider  wallet.id
 */
func (s *SmartContract) CreateFile(ctx contractapi.TransactionContextInterface, id, author, path string, state, period, rentPrice, price int, provider string) error {
	// check for file existance
	if exists, err := s.FileExists(ctx, id); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("The file %s already exists", id)
	}
	params := []string{"BalanceOf", author}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}
	response := ctx.GetStub().InvokeChaincode("token", queryArgs, "mychannel")
	if response.Status != 200 {
		return fmt.Errorf("Failed to query token chaincode Got Error: %s", string(response.Payload))
	}
	toInt, err := strconv.Atoi(string(response.Payload))
	if err != nil {
		return fmt.Errorf("Error occured while parsing user balance: %s", err)
	} else if toInt < rentPrice {
		return fmt.Errorf("User amount is not enough to handle this transaction")
	}

	file := FileData{
		ID:        id,
		Author:    author,
		Path:      path,
		State:     state,
		Period:    period,
		RentPrice: rentPrice,
		Price:     price,
		Provider:  provider}

	if fileJSON, err := json.Marshal(file); err != nil {
		return err
	} else {
		param := []string{"Transfer", author, "admin", strconv.Itoa(rentPrice)}
		invokeArgs := make([][]byte, len(param))
		for i, arg := range param {
			invokeArgs[i] = []byte(arg)
		}
		resp := ctx.GetStub().InvokeChaincode("token", invokeArgs, "mychannel")
		if resp.Status != 200 {
			return fmt.Errorf("Failed to query token chaincode Got Error: %s", resp.Payload)
		}
		respBool, err := strconv.ParseBool(string(resp.Payload))
		if !respBool {
			return fmt.Errorf("Failed to transcat token to admin")
		} else if err != nil {
			return fmt.Errorf("Failed to transcat token got Error %s", err)
		}

		return ctx.GetStub().PutState(id, fileJSON)
	}

}

/**
 * function: CreateFile
 *
 * @param {Context} ctx the transaction context
 * @param {String} id the string id  of the file
 * @param {String} author The author of the file aka 'Requestor wallet.id'
 * @param {Int} state the state of the file
 */
func (s *SmartContract) UpdateFileProgress(ctx contractapi.TransactionContextInterface, id, provider string, state int) (bool, error) {
	var err error
	if exists, err := s.FileExists(ctx, id); err != nil {
		return false, err
	} else if exists {
		fileJSON, _ := ctx.GetStub().GetState(id)
		var tempFile FileData
		if err = json.Unmarshal(fileJSON, &tempFile); err != nil {
			return false, err
		}
		if tempFile.Provider != provider {

		}
		tempFile.State = state
		updatedFile, err := json.Marshal(&tempFile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error occured while marshiling json: %s", err))
			return false, err
		}
		if state == 200 {
			param := []string{"Transfer", "admin", provider, strconv.Itoa(tempFile.RentPrice)}
			invokeArgs := make([][]byte, len(param))
			for i, arg := range param {
				invokeArgs[i] = []byte(arg)
			}
			if resp := ctx.GetStub().InvokeChaincode("token", invokeArgs, "mychannel"); resp.Status != 200 {
				return false, fmt.Errorf("Failed to transfer token to provider from admin")
			}

		}
		return true, ctx.GetStub().PutState(id, updatedFile)
	}

	return false, err
}

/* check file exists
 *this function strictly called inside chaincode
 */
func (s *SmartContract) FileExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	fileJSON, err := ctx.GetStub().GetState(id)

	if err != nil {
		return false, fmt.Errorf("failde to read from world state: %v", err)
	}

	return fileJSON != nil, nil
}

/*
	CancelFileProgress
 * @param {Context}	ctx the ransaction context
 * @param {String} id of the file
 * @param {String} author the author of the file aka 'Requestor wallet.id'
 * @param {Int} state the state of the file
*/
func (s *SmartContract) CancelFileProgress(ctx contractapi.TransactionContextInterface, id, author string, state int) (bool, error) {
	exists, err := s.FileExists(ctx, id)
	if err != nil {
		return false, err
	} else if exists {
		fileJSON, _ := ctx.GetStub().GetState(id)
		var tempFile FileData
		if err = json.Unmarshal(fileJSON, &tempFile); err != nil {
			return false, err
		}
		if tempFile.Author != author {
			return false, fmt.Errorf("plese check your account, you are not authorized  to call this contract")
		}

		tempFile.State = state
		canceledFile, err := json.Marshal(&tempFile)
		if err != nil {
			fmt.Println(fmt.Sprintf("error occured while marshiling json %s", err))
			return false, err
		}
		if err = ctx.GetStub().PutState(id, canceledFile); nil == err {
			return true, nil
		}
	}
	return false, nil
}

/**
 * OrderFileFromAouthor
 *
 * @param {Context} ctx the transaction\
 * @param {String} id of the file
 * @param {String} clinet the client wallet.id
 */

func (s *SmartContract) OrderFileFromAuthor(ctx contractapi.TransactionContextInterface, id, client string) bool {
	exists, err := s.FileExists(ctx, id)
	if err != nil {
		// return false, &OrderFile{}
		return false
	} else if exists {
		fileJSON, _ := ctx.GetStub().GetState(id)
		var tempFile FileData
		if err = json.Unmarshal(fileJSON, &tempFile); err != nil {
			// return false, &OrderFile{}
			return false
		}
		// invoke token chaincode function
		// check the balance of the user
		params := []string{"BalanceOf", client}
		queryArgs := make([][]byte, len(params))
		for i, arg := range params {
			queryArgs[i] = []byte(arg)
		}
		response := ctx.GetStub().InvokeChaincode("token", queryArgs, "mychannel")
		if response.Status != 200 {
			//return false, &OrderFile{}
			return false
		}
		toInt, _ := strconv.Atoi(string(response.Payload))
		if tempFile.Price > toInt {
			//return false, &OrderFile{}
			return false
		}

		transferParam := []string{"Transfer", client, tempFile.Author, strconv.Itoa(tempFile.Price)}
		transferArgs := make([][]byte, len(transferParam))
		for j, args := range transferParam {
			transferArgs[j] = []byte(args)
		}
		resp := ctx.GetStub().InvokeChaincode("token", transferArgs, "mychannel")
		if resp.Status != 200 {
			return false
		}
		if respBool, err := strconv.ParseBool(string(resp.Payload)); !respBool || err != nil {
			return false
		}
		// var result OrderFile
		// result.Author = tempFile.Author
		// result.ID = tempFile.ID
		// result.Path = tempFile.Path
		// result.Price = tempFile.Price
		return true

	}

	return false
}

/**
 * GetaLLFiles
 *
 * @param {Context} ctx the transaction\
 */
func (s *SmartContract) GetAllFiles(ctx contractapi.TransactionContext) (error, []*FileData) {
	//range query with empty string for startKey and EndKey does an
	// open-ended query of all files in the chaincode namespace

	resultIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return err, nil
	}

	defer resultIterator.Close()

	var files []*FileData
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return err, nil
		}

		var file FileData
		if err = json.Unmarshal(queryResponse.Value, &file); err != nil {
			return err, nil
		}
		files = append(files, &file)
	}
	return nil, files
}

/**
 * DeleteFile
 *
 * @param {Context} ctx the transaction
 * @param {String} id the id of the file
 * @param {String} author the author of the file
 */
func (s *SmartContract) DeleteFile(ctx contractapi.TransactionContextInterface, id, author string) (error, bool) {
	if exists, err := s.FileExists(ctx, id); err != nil || !exists {
		return err, false
	}

	file, _ := ctx.GetStub().GetState(id)
	var fileJson FileData
	if err := json.Unmarshal(file, &fileJson); err != nil {
		return err, false
	}
	if &author != &fileJson.Author {
		return fmt.Errorf("author wallet address does not match"), false
	}

	return ctx.GetStub().DelState(id), true

}
