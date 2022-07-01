package ws

import (
	"context"
	"sync"
	"time"

	scalecodec "github.com/itering/scale.go"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type networkPallet struct {
	network string
	codec   *scalecodec.MetadataDecoder
	Err     error
}

func NetworkPalletList(ctx context.Context, node map[string]string) []subscan.NetworkPallet {
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
				pallet.codec, pallet.Err = getMetadataModules(ctx, addr)
				palletCh <- pallet
				wg.Done()
			}(name, addr)
		}
		wg.Wait()
		close(palletCh)
	}()

	data := make([]subscan.NetworkPallet, 0, len(node))
	statusCh, doneCh := utils.ProgressDisplay(len(node))
	var doneIdx = 0
	for p := range palletCh {
		doneIdx++
		statusCh <- doneIdx
		pl := subscan.NetworkPallet{
			Network: p.network,
		}
		if p.Err != nil {
			pl.Err = p.Err
		} else {
			pl.Err = p.codec.Process()
			if pl.Err == nil {
				pl.Pallet = make([]string, 0, len(p.codec.Metadata.Metadata.Modules))
				for _, m := range p.codec.Metadata.Metadata.Modules {
					pl.Pallet = append(pl.Pallet, m.Name)
				}
			}
		}
		data = append(data, pl)
	}
	close(statusCh)
	<-doneCh
	return data
}

func getMetadataModules(ctx context.Context, addr string) (*scalecodec.MetadataDecoder, error) {
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
