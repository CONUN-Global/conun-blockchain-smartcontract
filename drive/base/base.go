package base

import "github.com/golang/protobuf/ptypes/timestamp"

// errro responses
const (
	JSONParseError       = "1001 JSON Parsing Error"
	FileExistsError      = "1002 Content Already Exists on Blockchain"
	PutStateError        = "1003 Error saving to world state"
	GetstateError        = "1004 Error getting Data from world state"
	OwnerError           = "1005 Error owner are wrong Address"
	KeyCreationError     = "1006 Error unique Key creation"
	EmptyAllowance       = "1007 Allowance is Empty"
	CheckFileError       = "1008 Error Checking File"
	EventError           = "1009 Error Setting Event"
	EmptyFile            = "1010 Error File does not exists"
	AllownaceRightError  = "1011 Allowance is not difened for this wallet"
	NumberError          = "1012 Number is not integer"
	InvokeChaincodeError = "1013 Invoke Error"
	WrongAmount          = "1014 Wrong Amount"
)

// Content struct
type Content struct {
	Author   string `json:"author"`
	IpfsHash string `json:"ipfshash"`
	Price    string `json:"price"`
	Status   string `json:"status"`
}

//user action struct
type Action struct {
	Ccid      string `json:"ccid"`
	Wallet    string `json:"wallet"`
	UserID    string `json:"user_id"`
	ContentID string `json:"content_id"`
}

// Response struct
type Response struct {
	Fcn       string               `json:"Fcn,omitempty"`       // function name
	Success   bool                 `json:"Success,omitempty"`   // true if success
	TxID      string               `json:"TxID,omitempty"`      // transaction id
	Timestamp *timestamp.Timestamp `json:"Timestmap,omitempty"` // timestmap of the transaction
	Value     string               `json:"Value,omitempty"`     // value of content
}

// txDetails struct
type DetailsTx struct {
	From   string `json:"From"`
	To     string `json:"To"`
	Action string `json:"Action"`
	Value  string `json:"Value"`
}

// Event
type Event struct {
	UserID    string               `json:"UserID"`
	ContentID string               `json:"ContentID"`
	Timestamp *timestamp.Timestamp `json:"Timestmap"`
}
