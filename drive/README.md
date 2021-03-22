
# CONUN Drive Dapp Smart Contract:

 CONUN Drive Smart Contract for CONUN Drive dapp



#### Structs

- `FileData`: `File details struct`
- `Response` : `an organized Json struct response`
- `OrderFile` : `Order function response`


#### Methods
- `CreateFile` : `(ipfsHash string, author wallet address) => (JSON Object, error)`
- `Approve` : `(ipfsHash string, author,spenderAdr wallet address) => (JSON Object, error)` 
- `Allowance` : `Query (ipfsHash string, spender wallet address) => (JSON Object, error)`
- `LikeContent` : `(ipfsHash string, wallet address) => (JSON Object, error)`
- `CountDownloads` : `(ipfsHash string, wallet address) => (JSON Object, error)`
- `GetTotatLikes` : `(ipfsHash string) => (JSON Object, error)`
- `GetTotalDownloads`:`(ipfsHash string) => (JSON Object, error)`


  




