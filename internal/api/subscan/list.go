package subscan

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

func NetworkPalletList(ctx context.Context, networkNode []string) []model.NetworkData[string] {
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
	palletCh := make(chan model.NetworkData[string], concurrency)
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
				_, result, err := runtimeList(ctx, nw)
				palletCh <- model.NetworkData[string]{Network: nw, Data: result, Err: err}
				<-limitCh
				wg.Done()
			}(network)
		}
		close(limitCh)
		wg.Wait()
		close(palletCh)
	}()

	data := make([]model.NetworkData[string], 0, len(networkNode))
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

type Resp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type Modules struct {
	List []struct {
		SpecVersion int    `json:"spec_version"`
		Modules     string `json:"modules"`
	} `json:"list"`
}

func runtimeList(ctx context.Context, network string) (int, []string, error) {
	var listURL = fmt.Sprintf("https://%s.api.subscan.io/api/scan/", network) + "runtime/list"
	var (
		rspData *Resp[Modules]
		err     error
	)
	if rspData, err = sendRequest[Modules](ctx, listURL, nil); err != nil {
		return 0, nil, err
	}

	if len(rspData.Data.List) == 0 {
		return 0, nil, nil
	}
	// API保证了第一个值为最新spec
	result := rspData.Data.List[0]
	return result.SpecVersion, strings.Split(result.Modules, "|"), nil
}
