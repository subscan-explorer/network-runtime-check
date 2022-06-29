package ws

import "sync/atomic"

type Param struct {
	ID      int           `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Resp struct {
	ID      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error"`
}

func (r Resp) ResultVal() interface{} {
	return r.Result
}

func newParam(id int, method string, param []interface{}) Param {
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

var (
	msgID  = new(int64)
	UniqID = func() int {
		return int(atomic.AddInt64(msgID, 1))
	}
)
