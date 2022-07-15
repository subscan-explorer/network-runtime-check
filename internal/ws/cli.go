package ws

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/itering/scale.go/types"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

var dialer = websocket.Dialer{
	Proxy:            http.ProxyFromEnvironment,
	HandshakeTimeout: 15 * time.Second,
}

type Endpoint struct {
	conn  *websocket.Conn
	req   map[int]chan Resp
	reqMx sync.RWMutex
	close chan struct{}
}

func NewEndpoint(ctx context.Context, host string) (*Endpoint, error) {
	var (
		ep  = new(Endpoint)
		err error
	)
	if ep.conn, _, err = dialer.DialContext(ctx, host, nil); err != nil {
		return nil, err
	}
	ep.req = make(map[int]chan Resp, 128)
	ep.close = make(chan struct{})

	go ep.read()
	return ep, nil
}

func (e *Endpoint) read() {
	defer close(e.close)
	for {
		tp, msg, err := e.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) || errors.Is(err, io.EOF) {
				return
			}
			//log.Println(err.Error())
			return
		}

		switch tp {
		case websocket.TextMessage:
			var rsp Resp
			if err = json.Unmarshal(msg, &rsp); err != nil {
				log.Printf("failed to unmarshal data. data: [%s]", string(msg))
				continue
			}
			// match ID
			if rsp.ID != 0 {
				e.reqMx.RLock()
				req, ok := e.req[rsp.ID]
				e.reqMx.RUnlock()
				if ok {
					req <- rsp
					continue
				}
			}
			// match method
			// default other
			continue
		case websocket.CloseMessage:
			// close
			return
		default:
			// skip
			continue
		}
	}
}

func (e *Endpoint) Close() {
	err := e.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err == nil {
		select {
		case <-e.close:
		case <-time.After(time.Second):
		}
	} else {
		log.Println("close", err.Error())
	}
	_ = e.conn.Close()
}

func (e *Endpoint) request(ctx context.Context, param Param) (*Resp, error) {
	var ch = make(chan Resp, 1)
	e.reqMx.Lock()
	e.req[param.ID] = ch
	e.reqMx.Unlock()
	defer func() {
		e.reqMx.Lock()
		delete(e.req, param.ID)
		e.reqMx.Unlock()
	}()
	if err := e.conn.WriteJSON(param); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-e.close:
		return nil, errors.New("connection is close")
	case rsp := <-ch:
		return &rsp, nil
	}
}

func (e *Endpoint) GetMetadata(ctx context.Context) (*types.MetadataStruct, error) {
	var (
		param = newMetadataParam()
		rsp   *Resp
		meta  string
		err   error
	)
	if rsp, err = e.request(ctx, param); err != nil {
		return nil, err
	}
	if meta, err = rsp.ResultString(); err != nil {
		return nil, err
	}
	return utils.DecodeMetadata(meta)
}

func (e *Endpoint) GetStorage(ctx context.Context, blockHash string) (string, error) {
	var (
		param = newStateStorageParam("", blockHash)
		rsp   *Resp
		err   error
	)
	if rsp, err = e.request(ctx, param); err != nil {
		return "", err
	}
	return rsp.ResultString()
}

func (e *Endpoint) GetFinalizedHeadBlock(ctx context.Context) (string, error) {
	var (
		param = newChainGetFinalizedHeadParam()
		rsp   *Resp
		err   error
	)
	if rsp, err = e.request(ctx, param); err != nil {
		return "", err
	}
	return rsp.ResultString()
}

func (e *Endpoint) GetBlockHash(ctx context.Context, num int64) (string, error) {
	var (
		param = newChainGetBlockHashParam(num)
		rsp   *Resp
		err   error
	)
	if rsp, err = e.request(ctx, param); err != nil {
		return "", err
	}
	return rsp.ResultString()
}

func (e *Endpoint) GetBlock(ctx context.Context, hash string) (*BlockInfo, error) {
	var (
		param = newChainGetBlockParam(hash)
		rsp   *Resp
		err   error
	)
	if rsp, err = e.request(ctx, param); err != nil {
		return nil, err
	}
	return rsp.BlockResult()
}

func (e *Endpoint) GetRuntimeVersion(ctx context.Context, hash string) (*RuntimeVersion, error) {
	var (
		param = newChainGetRuntimeVersionParam(hash)
		rsp   *Resp
		err   error
	)
	if rsp, err = e.request(ctx, param); err != nil {
		return nil, err
	}
	return rsp.RuntimeVersion()
}
