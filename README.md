# network-runtime-check

### Build 
##### build binary
`make build`  
##### build docker image
`make image`

### Running Help
`./runtime-check -h`  

### Configure the configuration file
path `conf/config.yaml`

### Subscab
subscan supported networks

#### Pallet match
##### Shows all pallets supported by the Network Runtime

`-w` query network, default all  
`-p` matching pallet, default all  
`-o` output to file path

##### Example
`./runtime-check subscan pallet match`  

`docker run --name runtime-check --rm runtime-check bin/runtime-check subscan pallet match`

##### output
| Network  | Pallet                                        | 
|----------|-----------------------------------------------|
| polkadot | System Scheduler ... Preimage  Babe XcmPallet |
| kusama   | System Babe  ... Timestamp Indices Balances   |
| ...      | ...                                           |


##### Check if the network runtime supports a pallet
`./runtime-check subscan pallet match -p System,Babe`

`docker run --name runtime-check --rm runtime-check bin/runtime-check subscan pallet match -p System,Babe`  


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


#### Pallet compare
##### Network comparison with substrate standard pallet

`-w` query network, default all  
`-o` output to file path

##### Example
`./runtime-check subscan pallet compare`

`docker run --name runtime-check --rm runtime-check bin/runtime-check subscan pallet compare`

###### output
|         | statemint | stafi | sora |
|---------|-----------|-------|------|
| System  | O         | O     | O    |
| Utility | O         | O     | O    |
| Babe    | X         | O     | O    | 
| ...     | ...       | ...   | ...  |

##### Example
`./runtime-check subscan pallet compare -w stafi,sora`

`docker run --name runtime-check --rm runtime-check bin/runtime-check subscan pallet compare -w stafi,sora`

###### output
|         | stafi | sora |
|---------|-------|------|
| System  | O     | O    |
| Utility | O     | O    |
| Babe    | O     | O    | 
| ...     | ...   |      |


### Polkadot
polkadot supported networks

#### Pallet match
##### Shows all pallets supported by the Network Runtime

`-p` matching pallet, default all  
`-o` output to file path

##### Example
`./runtime-check polkadot pallet match`

`docker run --name runtime-check --rm runtime-check bin/runtime-check polkadot pallet match`

##### output
| Network  | Pallet                                        | 
|----------|-----------------------------------------------|
| polkadot | System Scheduler ... Preimage  Babe XcmPallet |
| kusama   | System Babe  ... Timestamp Indices Balances   |
| ...      | ...                                           |


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


#### Pallet compare
##### Network comparison with substrate standard pallet

`-o` output to file path

##### Example
`./runtime-check polkadot pallet compare`

`docker run --name runtime-check --rm runtime-check bin/runtime-check polkadot pallet compare`

###### output
|         | statemint | stafi | sora |
|---------|-----------|-------|------|
| System  | O         | O     | O    |
| Utility | O         | O     | O    |
| Babe    | X         | O     | O    | 
| ...     | ...       | ...   | ...  |

