package filter

import (
	"github.com/subscan-explorer/network-runtime-check/internal/model"
)

type Filter interface {
	FilterPallet([]model.NetworkData[string]) []model.NetworkData[string]
}

type Default struct{}

func (Default) FilterPallet(list []model.NetworkData[string]) []model.NetworkData[string] {
	return list
}
