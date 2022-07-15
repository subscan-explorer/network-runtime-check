package ws

import (
	"context"
	"sync"
	"time"

	"github.com/itering/scale.go/types"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type networkPallet struct {
	network  string
	metadata *types.MetadataStruct
	Err      error
}

func NetworkPalletList(ctx context.Context, node map[string]string) []model.NetworkData[string] {
	if len(node) == 0 {
		return nil
	}
	palletCh := make(chan networkPallet, len(node))
	go func() {
		wg := new(sync.WaitGroup)
		for name, addr := range node {
			wg.Add(1)
			go func(network string, addr string) {
				pallet := networkPallet{network: network}
				pallet.metadata, pallet.Err = GetMetadataInfo(ctx, addr)
				palletCh <- pallet
				wg.Done()
			}(name, addr)
		}
		wg.Wait()
		close(palletCh)
	}()

	data := make([]model.NetworkData[string], 0, len(node))
	statusCh, doneCh := utils.ProgressDisplay(len(node))
	var doneIdx = 0
	for p := range palletCh {
		doneIdx++
		statusCh <- doneIdx
		pl := model.NetworkData[string]{
			Network: p.network,
		}
		if p.Err != nil {
			pl.Err = p.Err
		} else {
			pl.Data = make([]string, 0, len(p.metadata.Metadata.Modules))
			for _, m := range p.metadata.Metadata.Modules {
				pl.Data = append(pl.Data, m.Name)
			}
		}
		data = append(data, pl)
	}
	close(statusCh)
	<-doneCh
	return data
}

func GetMetadataInfo(ctx context.Context, addr string) (*types.MetadataStruct, error) {
	var (
		ep             *Endpoint
		subCtx, cancel = context.WithTimeout(ctx, time.Second*30)
		err            error
	)
	defer cancel()
	if ep, err = NewEndpoint(subCtx, addr); err != nil {
		return nil, err
	}
	defer ep.Close()

	return ep.GetMetadata(subCtx)
}
