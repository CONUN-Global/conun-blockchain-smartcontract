# Wrapchain Smart Contract:

 Wrapchain Smart Contract For storing user activites that happened in the blockchain



#### Structs

- `Action`: `Action details struct`
- `UserArray` : `user Actions struct`


#### Methods
- `ActionWrite` : `(wallet address, int actionId, string ccid, txId) => (JSON Object, error)`
- `ActionExists` : `(int actionId) => (JSON Object, error)` 
- `GetActionById` : `(int actionId) => (JSON Object, error)`
- `GetUserActions` : `(wallet address) => (JSON Object, error)`


