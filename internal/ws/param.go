package ws

import "sync/atomic"

var (
	msgID  = new(int64)
	UniqID = func() int {
		return int(atomic.AddInt64(msgID, 1))
	}
)

type Param struct {
	ID      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func newParam(id int, method string, param interface{}) Param {
	return Param{
		ID:      id,
		Jsonrpc: "2.0",
		Method:  method,
		Params:  param,
	}
}

func newMetadataParam() Param {
	return newParam(UniqID(), "state_getMetadata", nil)
}

func newChainGetBlockParam(hash string) Param {
	return newParam(UniqID(), "chain_getBlock", []string{hash})
}

func newChainGetRuntimeVersionParam(hash string) Param {
	return newParam(UniqID(), "chain_getRuntimeVersion", []string{hash})
}

func newChainGetBlockHashParam(num int64) Param {
	return newParam(UniqID(), "chain_getBlockHash", []int64{num})
}

func newChainGetFinalizedHeadParam() Param {
	return newParam(UniqID(), "chain_getFinalizedHead", nil)
}

func newStateStorageParam(storageKey, hash string) Param {
	param := []string{storageKey}
	if len(hash) != 0 {
		param = append(param, hash)
	}
	return newParam(UniqID(), "state_getStorageAt", param)
}
