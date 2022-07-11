package filter

import (
	"strings"

	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
)

func NewExist(pallet []string) *Exist {
	return &Exist{Pallet: pallet}
}

type Exist struct {
	Pallet []string
}

func (e Exist) FilterPallet(list []subscan.NetworkPallet) []subscan.NetworkPallet {
	if len(e.Pallet) == 0 {
		return list
	}
	palletSet := make(map[string]struct{})
	for _, p := range e.Pallet {
		palletSet[strings.ToLower(p)] = struct{}{}
	}
	result := make([]subscan.NetworkPallet, 0, len(list))
	for _, item := range list {
		pallet := make([]string, 0, len(e.Pallet))
		for _, p := range item.Pallet {
			if _, ok := palletSet[strings.ToLower(p)]; ok {
				pallet = append(pallet, p)
			}
		}
		if len(pallet) == 0 {
			continue
		}
		item.Pallet = pallet
		result = append(result, item)
	}
	return result
}
