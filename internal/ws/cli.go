package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	scalecodec "github.com/itering/scale.go"
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

func (e *Endpoint) GetMetadata(ctx context.Context) (*scalecodec.MetadataDecoder, error) {
	var (
		param = newMetadataParam()
		ch    = make(chan Resp, 1)
		err   error
	)
	// register req
	e.reqMx.Lock()
	e.req[param.ID] = ch
	e.reqMx.Unlock()
	defer func() {
		e.reqMx.Lock()
		delete(e.req, param.ID)
		e.reqMx.Unlock()
	}()
	if err = e.conn.WriteJSON(param); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-e.close:
		return nil, errors.New("connection is close")
	case rsp := <-ch:
		if data, ok := rsp.Result.(string); ok {
			md := new(scalecodec.MetadataDecoder)
			md.Init(utils.HexToBytes(data))
			return md, nil
		}
		return nil, fmt.Errorf("not the expected string type, type: %T", rsp.Result)
	}
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
