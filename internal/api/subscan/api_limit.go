package subscan

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

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

func sendRequest[T any](ctx context.Context, host string, reqBody io.Reader) (*Resp[T], error) {
	var (
		req        *http.Request
		rsp        *http.Response
		rspData    = new(Resp[T])
		retryCount = 3
		err        error
	)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, host, reqBody); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if conf.Conf.APIKey != "" {
		req.Header.Set("X-API-Key", conf.Conf.APIKey)
	}
	sendReq := func() (bool, time.Duration, error) {
		if rsp, err = api.HTTPCli.Do(req); err != nil {
			return true, 0, err
		}
		defer rsp.Body.Close()
		if rsp.StatusCode != 200 {
			_, _ = io.Copy(ioutil.Discard, rsp.Body)
			if rsp.StatusCode == http.StatusTooManyRequests {
				delay, _ := strconv.Atoi(rsp.Header.Get("Retry-After"))
				return true, time.Second * time.Duration(delay), errors.New("API rate limit exceeded")
			}
			return true, 0, errors.New(rsp.Status)
		}
		if err = json.NewDecoder(rsp.Body).Decode(&rspData); err != nil {
			return false, 0, err
		}
		if rspData.Code != 0 {
			return false, 0, errors.New(rspData.Message)
		}
		return false, 0, nil
	}

	for {
		if retry, delay, err := sendReq(); err != nil {
			if retry && retryCount > 0 {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
				}
				if delay == 0 {
					retryCount--
				}
				continue
			}
			return nil, err
		}
		break
	}
	return rspData, nil
}
