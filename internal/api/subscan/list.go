package subscan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/api"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type NetworkPallet struct {
	Network string
	Pallet  []string
	Err     error
}

func NetworkPalletList(ctx context.Context, networkNode []string) []NetworkPallet {
	concurrency := 2
	if conf.Conf.APIKey != "" {
		if c, err := APILimit(ctx); err != nil {
			log.Printf("Failed to get apikey concurrency limit. err: %s\n", err)
		} else {
			if c != 0 {
				concurrency = c
			}
		}
	}
	log.Printf("current concurrency: %d\n", concurrency)
	limitCh := make(chan struct{}, concurrency)
	palletCh := make(chan NetworkPallet, concurrency)
	go func() {
		wg := new(sync.WaitGroup)
	BEGIN:
		for _, network := range networkNode {
			select {
			case <-ctx.Done():
				break BEGIN
			case limitCh <- struct{}{}:
			}
			wg.Add(1)
			go func(nw string) {
				result, err := list(ctx, nw)
				palletCh <- NetworkPallet{Network: nw, Pallet: result, Err: err}
				<-limitCh
				wg.Done()
			}(network)
		}
		close(limitCh)
		wg.Wait()
		close(palletCh)
	}()

	data := make([]NetworkPallet, 0, len(networkNode))
	statusCh, doneCh := utils.ProgressDisplay(len(networkNode))
	var doneIdx = 0
	for p := range palletCh {
		doneIdx++
		statusCh <- doneIdx
		data = append(data, p)
	}
	close(statusCh)
	<-doneCh
	return data
}

type ListResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			Modules string `json:"modules"`
		} `json:"list"`
	} `json:"data"`
}

func list(ctx context.Context, network string) ([]string, error) {
	var metadataURL = fmt.Sprintf("https://%s.api.subscan.io/api/scan/", network) + "runtime/list"
	var (
		req        *http.Request
		rsp        *http.Response
		rspData    ListResp
		retryCount = 3
		err        error
	)
	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, metadataURL, nil); err != nil {
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
