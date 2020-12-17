package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define key names for options
const totalSupplyKey = "totalSupply"

// Define objectType names for prefix
const allowancePrefix = "allowance"

// name of the token
const namePrefix = "namePrefix"

// symbol of the token
const symbolPrefix = "SYB"

//decimal of the token
const decimalPrefix = "decimal"

//owner of the contract
const ownerPrefix = "owner"

// SmartContract provides functions for transferring tokens between accounts
type SmartContract struct {
	contractapi.Contract
}

// event provides an organized struct for emitting events
type event struct {
	from  string
	to    string
	value int
}

// respnse struct
type Response struct {
	Success   bool                 `json:"Success"`
	Func      *Fcn                 `json:"Func,omitempty"`
	TxID      string               `json:"TxID"`
	Timestamp *timestamp.Timestamp `json:"Timestamp"`
}

// function response struct
type Fcn struct {
	Minter string `json:"Minter,omitempty"`
	From   string `json:"From,omitempty"`
	To     string `json:"To,omitempty"`
	Amount int    `json:"Amount,omitempty"`
	Total  int    `json:"Total,omitempty"`
}

// info response struct
type Info struct {
	Owner     string `json:"Owner"`
	TokenName string `json:"TokenName"`
	Symbol    string `json:"Symbol"`
	Decimal   string `json:"Decimal"`
}

/*
	Init declares chaincode details

	@param {Context} ctx the transaction context
	@param {string} contract owner address

	Return success interface or error
*/
func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface, owner string) (interface{}, error) {

	exists, err := ctx.GetStub().GetState(ownerPrefix)
	if err != nil || exists != nil {
		return nil, fmt.Errorf("Contract already initalized by %s error:%s", string(exists), err)
	}
	err = ctx.GetStub().PutState(namePrefix, []byte("CONUN"))
	err = ctx.GetStub().PutState(symbolPrefix, []byte("CON"))
	err = ctx.GetStub().PutState(decimalPrefix, []byte(strconv.Itoa(18)))
	err = ctx.GetStub().PutState(ownerPrefix, []byte(owner))
	if err != nil {
		return nil, fmt.Errorf("error setting values %s", err)
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	res := &Response{
		Success:   true,
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return string(content), nil
}

/*
	Mint creates new tokens and adds them to contract owners account balance
	 // this function triggers a Transfer event

	@param {Context} ctx the transaction context
	@param {string} the contract owner address

	Return success interface or error
*/
func (s *SmartContract) Mint(ctx contractapi.TransactionContextInterface, amount int) (interface{}, error) {

	// retrieve contract owner address
	minterByte, err := ctx.GetStub().GetState(ownerPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed while getting minterAddress %s", err)
	} else if minterByte == nil {
		return nil, fmt.Errorf("Contract is not initialized yet")
	}
	minter := string(minterByte)
	// check if contract caller is contract owner
	ownerID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get user Address %s", err)
	}

	if verify, err := addressHelper(ownerID, minter); err != nil || !verify {
		return nil, fmt.Errorf("failed to Mint  Sender is not valid to Mint: %s", err)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("mint amount must be a positive integer")
	}

	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return nil, fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
	}

	var currentBalance int

	// If minter current balance doesn't yet exist, we'll create it with a current balance of 0
	if currentBalanceBytes == nil {
		currentBalance = 0
	} else {
		currentBalance, _ = strconv.Atoi(string(currentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
	}

	updatedBalance := currentBalance + amount

	err = ctx.GetStub().PutState(minter, []byte(strconv.Itoa(updatedBalance)))
	if err != nil {
		return nil, err
	}

	// Update the totalSupply
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve total token supply: %v", err)
	}

	var totalSupply int

	// If no tokens have been minted, initialize the totalSupply
	if totalSupplyBytes == nil {
		totalSupply = 0
	} else {
		totalSupply, _ = strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	// Add the mint amount to the total supply and update the state
	totalSupply += amount
	err = ctx.GetStub().PutState(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	if err != nil {
		return nil, err
	}

	// Emit the Transfer event
	transferEvent := event{"0x0", minter, amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)

	txTime, _ := ctx.GetStub().GetTxTimestamp()
	mintResp := &Fcn{
		Minter: minter,
		Amount: amount,
		Total:  updatedBalance,
	}
	res := &Response{
		Success:   true,
		Func:      mintResp,
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	context, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return string(context), nil
}

/*
	Burn redeems tokens the contract owner's account balance
	// this function triggers a Transfer event

	@param {Context} ctx the transaction context
	@param {string} the contract owner address

	Return success interface or error
*/
func (s *SmartContract) Burn(ctx contractapi.TransactionContextInterface, amount int) (interface{}, error) {

	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to burn new tokens
	// retrieve contract owner address
	minterByte, err := ctx.GetStub().GetState(ownerPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed while getting minterAddress %s", err)
	} else if minterByte == nil {
		return nil, fmt.Errorf("Contract is not initialized yet")
	}
	minter := string(minterByte)
	// check if contract caller is contract owner
	ownerID, err := ctx.GetClientIdentity().GetID()

	if verify, err := addressHelper(ownerID, minter); err != nil || !verify {
		return nil, fmt.Errorf("failed to Burn  Sender is not valid to Burn: %s", err)
	}

	if amount <= 0 {
		return nil, errors.New("burn amount must be a positive integer")
	}

	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return nil, fmt.Errorf("failed to read minter account %s from world state: %v", minter, err)
	}

	var currentBalance int

	// Check if minter current balance exists
	if currentBalanceBytes == nil {
		return nil, errors.New("The balance does not exist")
	}

	currentBalance, _ = strconv.Atoi(string(currentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	updatedBalance := currentBalance - amount

	err = ctx.GetStub().PutState(minter, []byte(strconv.Itoa(updatedBalance)))
	if err != nil {
		return nil, err
	}

	// Update the totalSupply
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve total token supply: %v", err)
	}

	// If no tokens have been minted, throw error
	if totalSupplyBytes == nil {
		return nil, errors.New("totalSupply does not exist")
	}

	totalSupply, _ := strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.

	// Subtract the burn amount to the total supply and update the state
	totalSupply -= amount
	err = ctx.GetStub().PutState(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	if err != nil {
		return nil, err
	}

	// Emit the Transfer event
	transferEvent := event{minter, "0x0", amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	mintResp := &Fcn{
		Minter: minter,
		Amount: amount,
		Total:  updatedBalance,
	}
	resp := &Response{
		Success:   true,
		Func:      mintResp,
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json %s", err)
	}
	return string(content), nil
}

/*
   Transfer transfers tokens from client account to recipient account
   // recipient account must be a valid clientID
   // this function triggers a Transfer event

   @param {Context} ctx the transcation context
   @param {string} client account address
   @param {string} recipient account address

   Returns success interface or error
*/
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, from, recipient string, amount int) (interface{}, error) {

	caller, err := ctx.GetClientIdentity().GetID()
	if verify, err := addressHelper(caller, from); err != nil || !verify {
		return nil, fmt.Errorf("failed to Transfer  Sender is not valid to Transfer: %s", err)
	}

	err = transferHelper(ctx, from, recipient, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer: %v", err)
	}

	// Emit the Transfer event
	transferEvent := event{from, recipient, amount}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to set event: %v", err)
	}
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	mintResp := &Fcn{
		From:   from,
		To:     recipient,
		Amount: amount,
	}
	resp := &Response{
		Success:   true,
		Func:      mintResp,
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: txTime,
	}
	content, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json %s", err)
	}

	return string(content), nil
}

// BalanceOf returns the balance of the given account
func (s *SmartContract) BalanceOf(ctx contractapi.TransactionContextInterface, account string) (int, error) {
	balanceBytes, err := ctx.GetStub().GetState(account)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		return 0, nil
	}

	balance, _ := strconv.Atoi(string(balanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	return balance, nil
}

/*
	ClientAccountBalance returns the balance of the requesting client's account

	@oaram {Context} ctx the transaction context

	Returns int value of the balance or error

*/
func (s *SmartContract) ClientAccountBalance(ctx contractapi.TransactionContextInterface) (int, error) {

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to get client id: %v", err)
	}

	balanceBytes, err := ctx.GetStub().GetState(clientID)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if balanceBytes == nil {
		return 0, fmt.Errorf("the account %s does not exist", clientID)
	}

	balance, _ := strconv.Atoi(string(balanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	return balance, nil
}

/*
	ClientAccountID returns the id of the requesting client's account
	// in this implementation, the client account ID is the clientId itself
	// users can use this function to get their own account id, which they can then give to otherss as the payment address

	@param {Context} ctx the transaction context

	Returns string user adress or error
*/
func (s *SmartContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {

	// Get ID of submitting client identity
	clientAccountID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client id: %v", err)
	}

	return clientAccountID, nil
}

/*
	TotalSupply returns the totoal token supply

	@param {Context} ctx the transaction context

	Return int the total supply or error
*/
func (s *SmartContract) TotalSupply(ctx contractapi.TransactionContextInterface) (int, error) {

	// Retrieve total supply of tokens from state of smart contract
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve total token supply: %v", err)
	}

	var totalSupply int

	// If no tokens have been minted, return 0
	if totalSupplyBytes == nil {
		totalSupply = 0
	} else {
		totalSupply, _ = strconv.Atoi(string(totalSupplyBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	log.Printf("TotalSupply: %d tokens", totalSupply)

	return totalSupply, nil
}

func (s *SmartContract) GetDetails(ctx contractapi.TransactionContextInterface) (interface{}, error) {
	deployer, err := ctx.GetStub().GetState(ownerPrefix)
	tokenName, err := ctx.GetStub().GetState(namePrefix)
	symbol, err := ctx.GetStub().GetState(symbolPrefix)
	decimal, err := ctx.GetStub().GetState(decimalPrefix)

	if err != nil {
		return nil, err
	}
	if decimal == nil || tokenName == nil || symbol == nil || deployer == nil {
		return nil, fmt.Errorf("Init is not declared %s,%s,%s", string(decimal), string(tokenName), string(deployer))
	}

	res := &Info{
		Owner:     string(deployer),
		TokenName: string(tokenName),
		Symbol:    string(symbol),
		Decimal:   string(decimal),
	}
	content, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return string(content), nil
}

/*
	Approve allows the spender to withdraw from the calling client's token account
	// the spender can withdraw multiple times if neccessary, up to the value amount
	// this function triggers an Approval event

	@param {Context} ctx the transaction context
	@param {string} spender the spender address
	@param {int} value the amount to approve

	Return success interface or error
*/
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, spender string, value int) error {

	// Get ID of submitting client identity
	owner, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, spender})
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Update the state of the smart contract by adding the allowanceKey and value
	err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(value)))
	if err != nil {
		return fmt.Errorf("failed to update state of smart contract for key %s: %v", allowanceKey, err)
	}

	// Emit the Approval event
	approvalEvent := event{owner, spender, value}
	approvalEventJSON, err := json.Marshal(approvalEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Approval", approvalEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("client %s approved a withdrawal allowance of %d for spender %s", owner, value, spender)

	return nil
}

/*
	Allowance returns the amount still available for the spender to withdraw from the owner

	@param {Context} ctx the transaction context
	@param {string} owner the owner address
	@param {spender} spender the spender address

	Returns int amount or error
*/
func (s *SmartContract) Allowance(ctx contractapi.TransactionContextInterface, owner string, spender string) (int, error) {

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, spender})
	if err != nil {
		return 0, fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Read the allowance amount from the world state
	allowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read allowance for %s from world state: %v", allowanceKey, err)
	}

	var allowance int

	// If no current allowance, set allowance to 0
	if allowanceBytes == nil {
		allowance = 0
	} else {
		allowance, err = strconv.Atoi(string(allowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.
	}

	log.Printf("The allowance left for spender %s to withdraw from owner %s: %d", spender, owner, allowance)

	return allowance, nil
}

/*
	TransferFrom transfers the value amount from the "from" address to the "to" address
	// this function triggers a Transfer event

	@param {string} from the from client address
	@param {string} to the to client address
	@param {int} value the amount to transfer

	Returns success interface or error
*/
func (s *SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, value int) error {

	// Get ID of submitting client identity
	spender, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	// Create allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{from, spender})
	if err != nil {
		return fmt.Errorf("failed to create the composite key for prefix %s: %v", allowancePrefix, err)
	}

	// Retrieve the allowance of the spender
	currentAllowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return fmt.Errorf("failed to retrieve the allowance for %s from world state: %v", allowanceKey, err)
	}

	var currentAllowance int
	currentAllowance, _ = strconv.Atoi(string(currentAllowanceBytes)) // Error handling not needed since Itoa() was used when setting the totalSupply, guaranteeing it was an integer.

	// Check if transferred value is less than allowance
	if currentAllowance < value {
		return fmt.Errorf("spender does not have enough allowance for transfer")
	}

	// Initiate the transfer
	err = transferHelper(ctx, from, to, value)
	if err != nil {
		return fmt.Errorf("failed to transfer: %v", err)
	}

	// Decrease the allowance
	updatedAllowance := currentAllowance - value
	err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(updatedAllowance)))
	if err != nil {
		return err
	}

	// Emit the Transfer event
	transferEvent := event{from, to, value}
	transferEventJSON, err := json.Marshal(transferEvent)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}
	err = ctx.GetStub().SetEvent("Transfer", transferEventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	log.Printf("spender %s allowance updated from %d to %d", spender, currentAllowance, updatedAllowance)

	return nil
}

/*
	Helper functions
	//transferHelper is a helper function that transfers tokens from the "from" address to the "to" address
	//dependant functions include Transfer and TransferFrom

	@param {Context} ctx the transaction context
	@oaram {string} from client address
	@param {string} to the recipient address
	@param {int} value the amount to transfer

	Returns error
*/
func transferHelper(ctx contractapi.TransactionContextInterface, from string, to string, value int) error {

	if value < 0 { // transfer of 0 is allowed in ERC-20, so just validate against negative amounts
		return fmt.Errorf("transfer amount cannot be negative")
	}

	fromCurrentBalanceBytes, err := ctx.GetStub().GetState(from)
	if err != nil {
		return fmt.Errorf("failed to read client account %s from world state: %v", from, err)
	}

	if fromCurrentBalanceBytes == nil {
		return fmt.Errorf("client account %s has no balance", from)
	}

	fromCurrentBalance, _ := strconv.Atoi(string(fromCurrentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.

	if fromCurrentBalance < value {
		return fmt.Errorf("client account %s has insufficient funds", from)
	}

	toCurrentBalanceBytes, err := ctx.GetStub().GetState(to)
	if err != nil {
		return fmt.Errorf("failed to read recipient account %s from world state: %v", to, err)
	}

	var toCurrentBalance int
	// If recipient current balance doesn't yet exist, we'll create it with a current balance of 0
	if toCurrentBalanceBytes == nil {
		toCurrentBalance = 0
	} else {
		toCurrentBalance, _ = strconv.Atoi(string(toCurrentBalanceBytes)) // Error handling not needed since Itoa() was used when setting the account balance, guaranteeing it was an integer.
	}

	fromUpdatedBalance := fromCurrentBalance - value
	toUpdatedBalance := toCurrentBalance + value

	err = ctx.GetStub().PutState(from, []byte(strconv.Itoa(fromUpdatedBalance)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(to, []byte(strconv.Itoa(toUpdatedBalance)))
	if err != nil {
		return err
	}

	log.Printf("client %s balance updated from %d to %d", from, fromCurrentBalance, fromUpdatedBalance)
	log.Printf("recipient %s balance updated from %d to %d", to, toCurrentBalance, toUpdatedBalance)

	return nil
}

func addressHelper(encodedAdr, client string) (bool, error) {

	decodedAdr, err := base64.StdEncoding.DecodeString(encodedAdr)
	if err != nil {
		return false, err
	} else if strings.Contains(string(decodedAdr), client) {
		return true, nil
	}
	return false, nil

}
