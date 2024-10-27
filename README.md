# divSvault

A CLI tool for navigating Hashicorp vault KV secrets under KV2 secret engine. 


## Usage
```
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=********************

divSvault -secretEngine="secrets" -output-format=text 
Search Secret: 
✔ example
foo: bar
```