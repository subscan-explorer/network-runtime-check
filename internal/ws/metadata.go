package ws

import (
	"context"
	"sync"

	"github.com/itering/scale.go/types"
	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type metadata struct {
	network  string
	metadata *types.MetadataStruct
	err      error
}

func GetMetadata(ctx context.Context, node []conf.ParamRule) []model.NetworkData[model.Metadata] {
	if len(node) == 0 {
		return nil
	}
	palletCh := make(chan metadata, len(node))
	go func() {
		wg := new(sync.WaitGroup)
		for _, n := range node {
			wg.Add(1)
			go func(n conf.ParamRule) {
				meta := metadata{}
				if len(n.Domain) != 0 {
					meta.network = n.Domain
				} else {
					meta.network = n.WsAddr
				}
				meta.metadata, meta.err = GetMetadataInfo(ctx, n.WsAddr)
				palletCh <- meta
				wg.Done()
			}(n)
		}
		wg.Wait()
		close(palletCh)
	}()

	data := make([]model.NetworkData[model.Metadata], 0, len(node))
	statusCh, doneCh := utils.ProgressDisplay(len(node))
	var doneIdx = 0
	for p := range palletCh {
		doneIdx++
		statusCh <- doneIdx
		meta := model.NetworkData[model.Metadata]{
			Network: p.network,
		}
		if p.err != nil {
			meta.Err = p.err
		} else {
			meta.Data = utils.TransformMetadata(p.metadata)
		}
		data = append(data, meta)
	}
	close(statusCh)
	<-doneCh
	return data
}
