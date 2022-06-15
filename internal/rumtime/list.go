package rumtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ListResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			Modules string `json:"modules"`
		} `json:"list"`
	} `json:"data"`
}

func List(ctx context.Context, network string) ([]string, error) {
	var metadataUrl = fmt.Sprintf(HostUrl, network) + "runtime/list"
	var (
		req        *http.Request
		rsp        *http.Response
		rspData    ListResp
		retryCount = 3
		err        error
	)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, metadataUrl, nil); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if APIKey != "" {
		req.Header.Set("X-API-Key", APIKey)
	}
	sendReq := func() (bool, time.Duration, error) {
		if rsp, err = HttpCli.Do(req); err != nil {
			return true, 0, err
		}
		defer rsp.Body.Close()
		if rsp.StatusCode != 200 {
			_, _ = io.Copy(ioutil.Discard, rsp.Body)
			if rsp.StatusCode == http.StatusTooManyRequests {
				delay, _ := strconv.Atoi(rsp.Header.Get("Retry-After"))
				return true, time.Second * time.Duration(delay), errors.New("API rate limit exceeded")
			}
			return false, 0, errors.New(rsp.Status)
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

	if len(rspData.Data.List) == 0 {
		return nil, nil
	}
	// API保证了第一个值为最新spec
	return strings.Split(rspData.Data.List[0].Modules, "|"), nil
}
