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

`-w` query subscan network name, support websocket address, default all subscan network name  
`-p` matching pallet, default all  
`-e` Exclude supported pallets, default empty
`-o` output to file path

##### Example

`./runtime-check pallet match`

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet match`

##### output

| Network  | Pallet                                        | 
|----------|-----------------------------------------------|
| polkadot | System Scheduler ... Preimage  Babe XcmPallet |
| kusama   | System Babe ... Timestamp Indices Balances    |
| ...      | ...                                           |

##### Example

`./runtime-check pallet match -w stafi,sora,wss://astar.api.onfinality.io/public-ws`

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet match -w stafi,sora,wss://astar.api.onfinality.io/public-ws`

##### output

| Network                 | Pallet                                        | 
|-------------------------|-----------------------------------------------|
| polkadot                | System Scheduler ... Preimage  Babe XcmPallet |
| kusama                  | System Babe ... Timestamp Indices Balances    |
| astar.api.onfinality.io | System  Utility ... Identity  Timestamp       |

##### Check if the network runtime supports a pallet

`./runtime-check pallet match -p System,Babe`

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

##### Exclude supported pallets

`./runtime-check pallet match -e babe,timestamp -p preimage,xcmpallet`

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet match -e babe,timestamp -p preimage,xcmpallet`

##### output

| Network  | Pallet              | 
|----------|---------------------|
| polkadot | Preimage  XcmPallet |
| kusama   | Preimage  XcmPallet |
| acala    | Preimage            |
| ...      | ...                 |

#### Pallet compare

##### Network comparison with substrate standard pallet

`-w` query subscan network name, support websocket address, default all subscan network name   
`-o` output to file path

##### Example

`./runtime-check pallet compare`

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet compare`

###### output

|         | statemint | stafi | sora |
|---------|-----------|-------|------|
| System  | O         | O     | O    |
| Utility | O         | O     | O    |
| Babe    | X         | O     | O    | 
| ...     | ...       | ...   | ...  |

##### Example

`./runtime-check pallet compare -w stafi,sora,wss://astar.api.onfinality.io/public-ws`

`docker run --name runtime-check --rm runtime-check bin/runtime-check pallet compare -w stafi,sora,wss://astar.api.onfinality.io/public-ws`

###### output

|         | stafi | sora | astar.api.onfinality.io |
|---------|-------|------|-------------------------|
| System  | O     | O    | O                       |
| Utility | O     | O    | O                       |
| Babe    | O     | O    | X                       |
| ...     | ...   |      |                         |

#### Param

##### Check whether the extrinsic and event in the pallet conform to the parameter definition

`-r` rule configuration file path
`-o` output to file path

##### Config file rule

```yaml
rule:
  - name: common
    pallet:
      - name: Balances
        event:
          - name: Transfer
            param: [ AccountId,AccountId,Balance ]
        extrinsic:
          - name: Transfer
            param: [ Address, Balance ]
      - name: System
        event:
          - name: ExtrinsicSuccess
            param: [ DispatchInfo ]
          - name: ExtrinsicFailed
            param: [ DispatchError,DispatchInfo ]
  - name: pmk
    pallet:
      - name: PredictionMarkets
        event:
          - name: MarketDisputed
            param: [ "MarketIdOf","MarketStatus","MarketDispute<AccountId, BlockNumber>" ]
          - name: MarketRejected
            param: [ MarketIdOf ]
        extrinsic:
          - name: buy_complete_set
            param: [ MarketIdOf,BalanceOf ]
      - name: ParachainSystem
        event:
          - name: DownwardMessagesProcessed
            param: [ Weight,relay_chain::Hash ]

network:
  - name: polkadot # name key
    domain: polkadot # subscan domain
    wsAddr: wss://rpc.polkadot.io/ # websocket addr  Priority use
    rule_inherit: [ common,pmk ]    # inherit rules from `rule` field
    pallet:
      - name: Treasury
        event:
          - name: Deposit # event id
            param: [ BalanceOf ]
        extrinsic:
          - name: Awarded # extrinsic id
            param: [ ProposalIndex, BalanceOf, AccountId ]
```

##### Example

`./runtime-check param -r rule.yml`

`docker run --name runtime-check --rm runtime-check bin/runtime-check param -r rule.yml`

###### output

| Network  | Pallet   | Event    | Check | Note |
|----------|----------|----------|-------|------|
| polkadot | Balances | Transfer | O     |      |
| khala    | Balances | Transfer | O     |      |
| ...      | ...      |          |       |      |

