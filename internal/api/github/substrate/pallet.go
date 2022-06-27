package substrate

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/subscan-explorer/network-runtime-check/internal/api"
)

type palletResp struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func PalletList(ctx context.Context) ([]string, error) {
	const palletURL = "https://api.github.com/repos/paritytech/substrate/contents/bin/node/runtime/src/lib.rs?ref=master"
	var (
		req        *http.Request
		rspData    palletResp
		retryCount = 3
		err        error
	)
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, palletURL, nil); err != nil {
		return nil, err
	}

	sendReq := func() (bool, error) {
		var rsp *http.Response
		if rsp, err = api.HTTPCli.Do(req); err != nil {
			return true, err
		}
		defer rsp.Body.Close()
		if rsp.StatusCode != 200 {
			_, _ = io.Copy(ioutil.Discard, rsp.Body)
			if rsp.StatusCode >= 500 {
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
	var content []byte
	if content, err = base64.StdEncoding.DecodeString(rspData.Content); err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.New("failed to get substrate pallet")
	}
	list := strings.Split(string(content), "\n")
	for i := 0; i < len(list); i++ {
		if strings.HasPrefix(list[i], "construct_runtime!(") {
			list = list[i+1:]
			break
		}
	}

	for i := 0; i < len(list); i++ {
		if strings.TrimSpace(list[i]) == "{" {
			list = list[i+1:]
			break
		}
	}
	var result []string
	for _, str := range list {
		s := strings.TrimSpace(str)
		if strings.HasPrefix(s, "//") {
			continue
		}
		if s == "}" {
			break
		}
		p := strings.Split(s, ":")
		if len(p) > 1 {
			result = append(result, p[0])
		}
	}
	return result, nil
}
