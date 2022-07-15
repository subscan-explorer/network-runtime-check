package filter

import (
	"strings"

	"github.com/subscan-explorer/network-runtime-check/internal/model"
)

func NewExclude(pallet []string) *Exclude {
	return &Exclude{Pallet: pallet}
}

type Exclude struct {
	Pallet []string
}

func (e Exclude) FilterPallet(list []model.NetworkData[string]) []model.NetworkData[string] {
	if len(e.Pallet) == 0 {
		return list
	}
	palletSet := make(map[string]struct{})
	for _, p := range e.Pallet {
		palletSet[strings.ToLower(p)] = struct{}{}
	}

	result := make([]model.NetworkData[string], 0, len(list))
	for _, item := range list {
		if len(item.Data) == 0 {
			continue
		}
		exist := false
		for _, p := range item.Data {
			if _, ok := palletSet[strings.ToLower(p)]; ok {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		result = append(result, item)
	}
	return result
}
