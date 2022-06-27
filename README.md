# network-runtime-check

### Build 
##### build binary
`make build`  
##### build docker image
`make image`

### Running Help
`./bin/runtime-check -h`  

### Configure the configuration file
path `conf/config.yaml`

### Pallet match
#### Shows all pallets supported by the Network Runtime

`-w` query network, default all  
`-p` matching pallet, default all  
`-o` output to file path

#### Example
`./bin/runtime-check pallet match`  

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet match`

##### output
| Network  | Pallet                                        | 
|----------|-----------------------------------------------|
| polkadot | System Scheduler ... Preimage  Babe XcmPallet |
| kusama   | System Babe  ... Timestamp Indices Balances   |
| ...      | ...                                           |


#### Check if the network runtime supports a pallet
`./bin/runtime-check pallet match -p System,Babe`

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet match -p System,Babe`  


##### output
| Network   | Pallet       | 
|-----------|--------------|
| polkadot  | System  Babe |
| kusama    | System  Babe |
| acala     | System       |
| darwinia  | System  Babe |
| alephzero | System       |
| altair    | System       |
| ...       | ...          |


### Pallet compare
#### Network comparison with substrate standard pallet

`-w` query network, default all  
`-o` output to file path

#### Example
`./bin/runtime-check pallet compare`

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet compare`

##### output
|         | statemint | stafi | sora |
|---------|-----------|-------|------|
| System  | O         | O     | O    |
| Utility | O         | O     | O    |
| Babe    | X         | O     | O    | 
| ...     | ...       | ...   | ...  |

#### Example
`./bin/runtime-check pallet compare -w stafi,sora`

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet compare -w stafi,sora`

##### output
|         | stafi | sora |
|---------|-------|------|
| System  | O     | O    |
| Utility | O     | O    |
| Babe    | O     | O    | 
| ...     | ...   |      |
