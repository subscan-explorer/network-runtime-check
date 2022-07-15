package ws

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SubParams struct {
	Result interface{} `json:"result"`
}

type Resp struct {
	ID      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Params  *SubParams  `json:"params,omitempty"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error"`
}

type ChainNewHeadLog struct {
	Logs []string `json:"logs"`
}

type ChainNewHeadResult struct {
	ExtrinsicsRoot string          `json:"extrinsicsRoot"`
	Number         string          `json:"number"`
	ParentHash     string          `json:"parentHash"`
	StateRoot      string          `json:"stateRoot"`
	Digest         ChainNewHeadLog `json:"digest"`
}

type RuntimeVersion struct {
	Apis             [][]interface{} `json:"apis"`
	AuthoringVersion int             `json:"authoringVersion"`
	ImplName         string          `json:"implName"`
	ImplVersion      int             `json:"implVersion"`
	SpecName         string          `json:"specName"`
	SpecVersion      int             `json:"specVersion"`
}

type BlockInfo struct {
	Block struct {
		Extrinsics []string           `json:"extrinsics"`
		Header     ChainNewHeadResult `json:"header"`
	} `json:"block"`
}

func (r Resp) ResultString() (string, error) {
	if err := r.checkErr(); err != nil {
		return "", err
	}
	if rs, ok := r.Result.(string); ok {
		return rs, nil
	}
	return "", fmt.Errorf("not the expected type, type: %T", r.Result)
}

func (r Resp) checkErr() error {
	if r.Error != nil {
		return errors.New(r.Error.Message)
	}
	return nil
}

func (r Resp) BlockResult() (*BlockInfo, error) {
	if err := r.checkErr(); err != nil {
		return nil, err
	}
	if r.Result == nil {
		return nil, errors.New("result is empty")
	}
	if result, ok := r.Result.(map[string]interface{}); ok {
		if len(result) == 0 {
			return nil, errors.New("result is empty")
		}
		rs := new(BlockInfo)
		m, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		return rs, json.Unmarshal(m, rs)
	}
	return nil, fmt.Errorf("not the expected type, type: %T", r.Result)
}

func (r Resp) RuntimeVersion() (*RuntimeVersion, error) {
	if err := r.checkErr(); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if r.Params != nil {
		result = r.Params.Result.(map[string]interface{})
	} else {
		result = r.Result.(map[string]interface{})
	}
	if len(result) == 0 {
		return nil, errors.New("result is empty")
	}
	rs := new(RuntimeVersion)
	m, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return rs, json.Unmarshal(m, rs)
}

type ChainEvent struct {
	EventID      string       `json:"event_id"`
	EventIdx     int          `json:"event_idx"`
	ExtrinsicIdx int          `json:"extrinsic_idx"`
	ModuleID     string       `json:"module_id"`
	Params       []EventParam `json:"params"`
	Phase        int          `json:"phase"`
}

type EventParam struct {
	Type     string      `json:"type"`
	TypeName string      `json:"type_name,omitempty"`
	Value    interface{} `json:"value"`
}

type ChainExtrinsic struct {
	CallCode           string           `json:"call_code"`
	CallModule         string           `json:"call_module"`
	CallModuleFunction string           `json:"call_module_function"`
	Params             []ExtrinsicParam `json:"params"`
	Nonce              int              `json:"nonce"`
	Era                string           `json:"era"`
}

type ExtrinsicParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	TypeName string      `json:"type_name"`
	Value    interface{} `json:"value"`
}
