package utils

import (
	"sync"

	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
)

var metadataDecodeLock sync.Mutex

func DecodeMetadata(meta string) (*types.MetadataStruct, error) {
	metadataDecodeLock.Lock()
	defer metadataDecodeLock.Unlock()
	md := new(scalecodec.MetadataDecoder)
	md.Init(HexToBytes(meta))
	if err := md.Process(); err != nil {
		return nil, err
	}
	return &md.Metadata, nil
}

func TransformMetadata(meta *types.MetadataStruct) (result []model.Metadata) {
	if meta == nil {
		return nil
	}
	for _, m := range meta.Metadata.Modules {
		md := model.Metadata{}
		md.Name = m.Name
		md.Prefix = m.Prefix
		md.Events = make([]model.MetadataEvent, 0, len(m.Events))
		for _, event := range m.Events {
			md.Events = append(md.Events, model.MetadataEvent{
				Lookup:       event.Lookup,
				Name:         event.Name,
				Args:         event.Args,
				ArgsTypeName: event.ArgsTypeName,
			})
		}
		for _, call := range m.Calls {
			c := model.MetadataCalls{
				Lookup: call.Lookup,
				Name:   call.Name,
			}
			for _, arg := range call.Args {
				c.Args = append(c.Args, model.MetadataModuleCallArgument{
					Name:     arg.Name,
					Type:     arg.Type,
					TypeName: arg.TypeName,
				})
			}
			md.Calls = append(md.Calls, c)
		}
		result = append(result, md)
	}
	return
}
