package base

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
