
# ERC20Token:

 The ERC-20 token smart contract demonstrates how to create and transfer Erc20 based tokens using account bassed model. There is an account for each participant that holds a balance tokens.

 * The <i>Init</i> transaction initializes Token details: name,symbol,decimal, owner.    
   * Init(String: owner) --> Return Json or Error
   * GetDetails() --> Return Json or Error
   * BalanceOf(String: userAccount) --> Return Int balance or Error.
  * A <i>Transfer</i> transaction debits the caller's account and creadits another account
    * Transfer(string: From, String: To, Int: Amount) --> Return Json or Error.
  * A <i>Mint</i> transaction creates tokens in the contract owner account.
    * Mint(Int:Amount) Return Json or Error.
  * A <i>Burn</i> transaction burnes tokens from the contract owner accout
    * Burn(Int:Amount) Return Json or Error.
  

<br/>
<br/>



