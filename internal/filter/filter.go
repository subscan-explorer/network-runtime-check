package filter

import (
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
)

type Filter interface {
	FilterPallet([]subscan.NetworkPallet) []subscan.NetworkPallet
}

type Default struct{}

func (Default) FilterPallet(list []subscan.NetworkPallet) []subscan.NetworkPallet {
	return list
}
