# network-runtime-check

### Build 
make build  
make image

### Running Help
`./bin/runtime-check -h`  

#### Shows all pallets supported by the Network Runtime
`./bin/runtime-check`  

`docker run --name runtime-check --rm runtime-check`


##### output
| Network  | Pallet             | 
|----------|--------------------|
| polkadot | System &#124; Scheduler &#124; ... &#124; Preimage &#124; Babe &#124; XcmPallet |
| kusama |  System &#124; Babe &#124; ... &#124; Timestamp &#124; Indices &#124; Balances  |
| ...      | ...                |


#### Check if the network runtime supports a pallet
`./bin/runtime-check -pallet=System,Babe`

`docker run --name runtime-check --rm runtime-check bin/runtime-check -pallet=System,Babe`  


##### output
| Network  | Pallet             | 
|----------|--------------------|
| polkadot | System &#124; Babe |
| kusama   | System &#124; Babe |
| acala    | System             |
| darwinia | System &#124; Babe |
| alephzero| System             |
| altair   | System             |
| ...      | ...                |


#### Add APIKey to speed up
`SUBSCAN_API_KEY={{token}} ./bin/runtime-check -pallet=System,Babe`

`docker run --name runtime-check --rm -e SUBSCAN_API_KEY={{token}} runtime-check bin/runtime-check -pallet=System,Babe`

