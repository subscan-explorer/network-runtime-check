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
	var palletSet map[string]struct{}
	if len(e.Pallet) != 0 {
		palletSet = make(map[string]struct{})
		for _, p := range e.Pallet {
			palletSet[strings.ToLower(p)] = struct{}{}
		}
	}

	for i := 0; i < len(list); i++ {
		if len(list[i].Pallet) == 0 {
			continue
		}
		pallet := make([]string, 0, len(e.Pallet))
		for _, p := range list[i].Pallet {
			if _, ok := palletSet[strings.ToLower(p)]; ok {
				pallet = append(pallet, p)
			}
		}
		list[i].Pallet = pallet
	}
	return list
}
