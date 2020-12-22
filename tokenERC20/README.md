# ERC-20 Token Smart Contract

ERC-20 Token smart contract written in javascript. 

#### Methods
- `balanceOf` : `query (wallet address) => (int)`
- `transfer` : `(wallet address from, wallet address to, int amount) => (bool)`
- `mint` : `(wallet address owner, amount int) => (bool)`
- `burn` : `(wallet address owner, amount int) => (bool)`


#### Events
- `transfer(wallet address _from, wallet address _to, amount int)`


#### Build Environment
- Node.js v12.20.0
- Npm 6.14.8
