package match

import (
	"context"
	"log"
	"strings"
	"sync/atomic"

	"network-runtime-check/internal/rumtime"

	"github.com/modood/table"
)

var networkNode = []string{"polkadot", "kusama", "darwinia", "acala", "acala-testnet", "alephzero", "astar",
	"altair", "bifrost", "calamari", "centrifuge", "chainx", "clover", "clv", "clover-testnet",
	"composable", "crab-parachain", "crust", "maxwell", "shadow", "dali", "crab",
	"dbc", "dock", "dolphin", "edgeware", "encointer", "equilibrium", "integritee", "interlay", "karura", "kintsugi",
	"kulupu", "khala", "kilt-testnet", "spiritnet", "litmus", "moonbase",
	"moonbeam", "moonriver", "nodle", "pangolin", "pangolin-parachain", "pangoro", "parallel", "parallel-heiko",
	"picasso", "pioneer", "polkadex", "polymesh", "polymesh-testnet", "plasm", "quartz", "reef",
	"rococo", "sakura", "shibuya", "shiden", "sora", "subgame", "stafi", "statemine", "statemint", "turing", "uniarts",
	"westend", "zeitgeist"}

type networkPallet struct {
	network string
	pallet  []string
}

type Pallet struct {
	Network string `json:"network"`
	Pallet  string `json:"Pallet"`
}

func NetworkPalletMatch(ctx context.Context, pallet string) {
	concurrency := 1
	if rumtime.APIKey != "" {
		concurrency = 5
	}
	log.Printf("concurrency: %d\n", concurrency)
	networkCh := make(chan string, len(networkNode))
	palletCh := make(chan networkPallet, concurrency)
	for _, network := range networkNode {
		networkCh <- network
	}
	close(networkCh)
	closeCount := new(int64)

	for i := 0; i < concurrency; i++ {
		go func() {
			// work 全部关闭后关闭channel
			defer func() {
				if atomic.AddInt64(closeCount, 1) == int64(concurrency) {
					close(palletCh)
				}
			}()
			for {
				select {
				case <-ctx.Done():
					return
				case nw, ok := <-networkCh:
					if !ok {
						return
					}
					result := matchPallet(ctx, nw, pallet)
					palletCh <- networkPallet{network: nw, pallet: result}
				}
			}
		}()
	}

	list := make([]Pallet, 0, len(networkNode))
	for p := range palletCh {
		list = append(list, Pallet{Network: p.network, Pallet: strings.Join(p.pallet, " | ")})
	}
	table.Output(list)
}

func matchPallet(ctx context.Context, network string, pallet string) []string {
	list, err := rumtime.List(ctx, network)
	if err != nil {
		// 把错误当做输出
		if strings.Contains(err.Error(), "timeout") {
			return []string{"timeout"}
		}
		return []string{err.Error()}
	}
	if len(pallet) == 0 {
		return list
	}
	var modelSet = make(map[string]struct{}, len(pallet))
	for _, model := range strings.Split(pallet, ",") {
		modelSet[model] = struct{}{}
	}
	var result = make([]string, 0, len(pallet))
	for _, model := range list {
		if _, ok := modelSet[model]; ok {
			result = append(result, model)
		}
	}
	return result
}
