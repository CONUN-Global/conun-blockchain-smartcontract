
# CONUN Drive Dapp Smart Contract:

 CONUN Drive Smart Contract for CONUN Drive dapp



#### Structs

- `FileData`: `File details struct`
- `Response` : `an organized Json struct response`
- `OrderFile` : `Order function response`


#### Methods
- `CreateFile` : `(author wallet address, argsp[ipfsHash string, priveCode string(optional)]) => (JSON Object, error)`
- `Approve` : `(priveCode string, author,spenderAdr wallet address) => (JSON Object, error)` 
- `Allowance` : `Query (priveCode string, spender wallet address) => (boolean, error)`
- `LikeContent` : `(ipfsHash string, wallet address) => (JSON Object, error)`
- `CountDownloads` : `(ipfsHash string, wallet address) => (JSON Object, error)`
- `GetTotalLikes` : `(ipfsHash string) => (JSON Object, error)`
- `GetTotalDownloads`:`(ipfsHash string) => (JSON Object, error)`


  




