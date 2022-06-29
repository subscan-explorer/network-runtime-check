package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/subscan-explorer/network-runtime-check/internal/api"
)

type fileInfo struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func GetFileContent(ctx context.Context, addrURL string) ([]byte, error) {
	var (
		req        *http.Request
		rspData    fileInfo
		retryCount = 3
		err        error
	)
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, addrURL, nil); err != nil {
		return nil, err
	}

	sendReq := func() (bool, error) {
		var rsp *http.Response
		if rsp, err = api.HTTPCli.Do(req); err != nil {
			return true, err
		}
		defer rsp.Body.Close()
		if rsp.StatusCode != http.StatusOK {
			_, _ = io.Copy(ioutil.Discard, rsp.Body)
			if rsp.StatusCode >= http.StatusInternalServerError {
				return true, errors.New(rsp.Status)
			}
			return false, errors.New(rsp.Status)
		}
		if err = json.NewDecoder(rsp.Body).Decode(&rspData); err != nil {
			return false, err
		}
		return false, nil
	}
	for {
		if retry, err := sendReq(); err != nil {
			if retry && retryCount > 0 {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
				}
				retryCount--
				continue
			}
			return nil, err
		}
		break
	}
	return base64.StdEncoding.DecodeString(rspData.Content)
}
