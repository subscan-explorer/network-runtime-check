package subscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

func NetworkMetadataList(ctx context.Context, node []string) []model.NetworkData[model.Metadata] {
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
	eventCh := make(chan model.NetworkData[model.Metadata], concurrency)
	go func() {
		wg := new(sync.WaitGroup)
	BEGIN:
		for _, n := range node {
			select {
			case <-ctx.Done():
				break BEGIN
			case limitCh <- struct{}{}:
			}
			wg.Add(1)
			go func(addr string) {
				var (
					ne model.NetworkData[model.Metadata]
					sv int
				)
				ne.Network = addr
				if sv, _, ne.Err = runtimeList(ctx, addr); ne.Err == nil {
					ne.Data, ne.Err = metadata(ctx, addr, sv)
				}
				eventCh <- ne
				<-limitCh
				wg.Done()
			}(n)
		}
		close(limitCh)
		wg.Wait()
		close(eventCh)
	}()
	data := make([]model.NetworkData[model.Metadata], 0, len(node))
	statusCh, doneCh := utils.ProgressDisplay(len(node))
	var doneIdx = 0
	for p := range eventCh {
		data = append(data, p)
		doneIdx++
		statusCh <- doneIdx

	}
	close(statusCh)
	<-doneCh
	return data
}

type Info struct {
	Info struct {
		Metadata []model.Metadata `json:"metadata"`
	} `json:"info"`
}

func metadata(ctx context.Context, network string, spec int) ([]model.Metadata, error) {
	var metadataURL = fmt.Sprintf("https://%s.api.subscan.io/api/scan/", network) + "runtime/metadata"
	var (
		param = new(struct {
			Spec int `json:"spec"`
		})
		rspData *Resp[Info]
		err     error
	)
	param.Spec = spec
	reqBody, _ := json.Marshal(param)

	if rspData, err = sendRequest[Info](ctx, metadataURL, bytes.NewBuffer(reqBody)); err != nil {
		return nil, err
	}
	return rspData.Data.Info.Metadata, nil
}
