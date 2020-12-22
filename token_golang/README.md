
# ERC20 CON Token Smart Contract:

 The ERC-20 CON Token smart contract used mainly in CONUN Private Blockchain as main asset. All rewards and payment are involved in with CON Token.


#### Constants

- `name` : `CONUN`
- `symbol` : `CON`
- `decimal` : (int)18


#### Structs

- `Event`: `event provides an organized struct for emmitng events`
- `Response` : `An organized Json struct response`
- `Fcn` : `function response struct`
- `Info` : `GetDetails function Response struct` 


#### Methods
- `Init` : `(wallet address owner) => (JSON Object, error)`
- `GetDetails` : `Query => (JSON Object, error)` 
- `BalanceOf` : `Query (wallet address) => (int)`
- `Transfer` : `(wallet address _from, wallet address _to, int amount) => (JSON Object, error)`
- `Mint` : `(int amont) => (JSON Object, error)`
- `Burn` : `(int amount) => (JSON Object, error)`
- `TransferFrom` : `(wallet address from, wallet address to, amount int) => (bool)`
- `Allowance` : `(wallet address owner, wallet address spender) => (bool)`
- `Approve` : `(wallet address user, amount int) => (bool)`


#### Inner Methods

- `transferHelper` : `transfer function extension`


#### Events

- `Transfer(wallet address from, wallet address to, amount int)`
  




