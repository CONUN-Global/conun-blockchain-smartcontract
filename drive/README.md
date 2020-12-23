
# CONUN Drive Dapp Smart Contract:

 CONUN Drive Smart Contract for CONUN Drive dapp



#### Structs

- `FileData`: `File details struct`
- `Response` : `an organized Json struct response`
- `OrderFile` : `Order function response`


#### Methods
- `CreateFile` : `(wallet address, fileData) => (JSON Object, error)`
- `UpdateFileProgress` : `(string fileID, wallet address, int state) => (JSON Object, error)` 
- `FileExists` : `Query (string fileID) => (bool)`
- `CancelFileProgress` : `(string fileID, wallet address, int state) => (JSON Object, error)`
- `OrderFileFromAuthor` : `(wallet address, string fileID) => (JSON Object, error)`
- `GetAllFiles` : `Query => (JSON Object, error)`
- `DeleteFile` : `(string fileID, wallet address) => (JSON Object, error)`


  




