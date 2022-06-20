package match

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/modood/table"
	"github.com/subscan-explorer/network-runtime-check/internal/runtime"
)

var networkNode = []string{"polkadot", "kusama", "darwinia", "acala", "acala-testnet", "alephzero", "astar",
	"altair", "bifrost", "calamari", "centrifuge", "chainx", "clover", "clv", "clover-testnet",
	"composable", "crab-parachain", "crust", "maxwell", "shadow", "dali", "crab",
	"dbc", "dock", "dolphin", "edgeware", "encointer", "equilibrium", "integritee", "interlay", "karura", "kintsugi",
	"kulupu", "khala", "kilt-testnet", "spiritnet", "litmus", "moonbase",
	"moonbeam", "moonriver", "nodle", "pangolin", "pangolin-parachain", "pangoro", "parallel", "parallel-heiko",
	"picasso", "pioneer", "polkadex", "polymesh", "polymesh-testnet", "quartz", "reef",
	"rococo", "sakura", "shibuya", "shiden", "sora", "subgame", "stafi", "statemine", "statemint", "turing", "uniarts",
	"westend", "zeitgeist"}

type Pallet struct {
	Network string `json:"network"`
	Pallet  string `json:"Pallet"`
}

func NetworkPalletMatch(ctx context.Context, pallet string) {
	concurrency := 2
	if runtime.APIKey != "" {
		if c, err := runtime.APILimit(ctx); err != nil {
			log.Printf("Failed to get apikey concurrency limit. err: %s\n", err)
		} else {
			if c != 0 {
				concurrency = c
			}
		}
	}
	log.Printf("current concurrency: %d\n", concurrency)
	limitCh := make(chan struct{}, concurrency)
	palletCh := make(chan Pallet, concurrency)
	go func() {
		wg := new(sync.WaitGroup)
	BEGIN:
		for _, network := range networkNode {
			select {
			case <-ctx.Done():
				break BEGIN
			case limitCh <- struct{}{}:
			}
			go func(nw string) {
				wg.Add(1)
				result := matchPallet(ctx, nw, pallet)
				palletCh <- Pallet{Network: nw, Pallet: strings.Join(result, " | ")}
				<-limitCh
				wg.Done()
			}(network)
		}
		close(limitCh)
		wg.Wait()
		close(palletCh)
	}()

	list := make([]Pallet, 0, len(networkNode))
	statusCh := make(chan string, 5)
	doneCh := make(chan struct{})
	go func() {
		for ch := range statusCh {
			fmt.Printf("\r%s", ch)
		}
		fmt.Printf("\rProcessing Complete!\n")
		close(doneCh)
	}()
	var doneIdx = 0
	for p := range palletCh {
		doneIdx++
		statusCh <- fmt.Sprintf("Processing: %d/%d", doneIdx, len(networkNode))
		list = append(list, p)
	}
	close(statusCh)
	<-doneCh
	table.Output(list)
}

func matchPallet(ctx context.Context, network string, pallet string) []string {
	list, err := runtime.List(ctx, network)
	if err != nil {
		// treat errors as output
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
		modelSet[strings.TrimSpace(strings.ToLower(model))] = struct{}{}
	}
	var result = make([]string, 0, len(pallet))
	for _, model := range list {
		if _, ok := modelSet[strings.TrimSpace(strings.ToLower(model))]; ok {
			result = append(result, model)
		}
	}
	return result
}
