package subscan

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/api"
)

func APILimit(ctx context.Context) (int, error) {
	var metadataURL = "https://kusama.api.subscan.io/api/now"
	var (
		req        *http.Request
		rsp        *http.Response
		rateLimit  int
		retryCount = 3
		err        error
	)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, metadataURL, nil); err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if conf.Conf.APIKey != "" {
		req.Header.Set("X-API-Key", conf.Conf.APIKey)
	}
	sendReq := func() (bool, error) {
		if rsp, err = api.HTTPCli.Do(req); err != nil {
			return true, err
		}
		defer rsp.Body.Close()
		_, _ = io.Copy(ioutil.Discard, rsp.Body)
		rateLimit, _ = strconv.Atoi(rsp.Header.Get("ratelimit-limit"))
		if rateLimit != 0 {
			return false, nil
		}
		return true, nil
	}

	for {
		if retry, err := sendReq(); err != nil {
			if retry && retryCount > 0 {
				select {
				case <-ctx.Done():
					return rateLimit, ctx.Err()
				default:
				}
				retryCount--
				continue
			}
			return rateLimit, err
		}
		break
	}
	return rateLimit, nil
}
