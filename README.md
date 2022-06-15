# network-runtime-check

### Build 
make build

### Running
./bin/runtime-check -h

#### 展示网络支持的所有 pallet
./bin/runtime-check
docker run --name runtime-check --rm runtime-check 


#### 检查网络是否支持某个 pallet
./bin/runtime-check -pallet=System,Babe
$ docker run --name runtime-check --rm runtime-check bin/runtime-check -pallet=System,Babe

#### 添加APIKey
$ SUBSCAN_API_KEY={{token}} ./bin/runtime-check -pallet=System,Babe
$ docker run --name runtime-check --rm -e SUBSCAN_API_KEY={{token}} runtime-check bin/runtime-check -pallet=System,Babe
