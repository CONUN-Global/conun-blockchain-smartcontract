
# CONUN Drive Dapp Smart Contract:

 CONUN Drive Smart Contract for CONUN Drive dapp



#### Structs

- `FileData`: `File details struct`
- `Response` : `an organized Json struct response`


#### Methods
- `CreateFile` : `(author wallet address, ipfsHash string) => (JSON Object, error)`
- `Approve` : `(privateCode string, author,spenderAdr wallet address) => (JSON Object, error)` 
- `Allowance` : `Query (priveCode string, spender wallet address) => (boolean, error)`
- `LikeContent` : `(privateCode string, wallet address, args[userID int, contentID int]int) => (JSON Object, error)`
- `CountDownloads` : `(privateCode string, wallet address, args[userID int, contentID int]int) => (JSON Object, error)`
- `GetFile` : `(privateCode string) => (JSON Object, error)`
- `GetTotalLikes` : `(privateCode string) => (JSON Object, error)`
- `GetTotalDownloads`:`(privateCode string) => (JSON Object, error)`


  




