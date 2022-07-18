package pallet

import (
	"context"
	"log"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/api/github/substrate"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/filter"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
	"github.com/subscan-explorer/network-runtime-check/internal/output"
	"github.com/subscan-explorer/network-runtime-check/internal/ws"
)

func NewPalletCmd() *cobra.Command {
	palletCmd := &cobra.Command{Use: "pallet", Short: "substrate pallet"}
	palletCmd.AddCommand(newCompareCmd(), newMatchCmd())
	return palletCmd
}

func newMatchCmd() *cobra.Command {
	matchCmd := &cobra.Command{
		Use:   "match",
		Short: "Pallets supported by subscan runtime",
		Long:  "Check subscan for all network supported pallets",
		Run: func(cmd *cobra.Command, args []string) {
			nw := strings.TrimSpace(cmd.Flag("network").Value.String())
			var networkName []string
			var websocketAddr []string
			if len(nw) == 0 {
				// default network
				networkName = conf.Conf.Network
			} else {
				for _, n := range strings.Split(nw, ",") {
					if len(n) == 0 {
						continue
					}
					if strings.HasPrefix(n, "wss://") || strings.HasPrefix(n, "ws://") {
						websocketAddr = append(websocketAddr, n)
					} else {
						networkName = append(networkName, n)
					}
				}
			}
			pe := strings.TrimSpace(cmd.Flag("pallet").Value.String())
			ep := strings.TrimSpace(cmd.Flag("exclude-pallet").Value.String())
			var fls []filter.Filter
			if len(pe) != 0 {
				fls = append(fls, filter.NewExist(strings.Split(pe, ",")))
			}
			if len(ep) != 0 {
				fls = append(fls, filter.NewExclude(strings.Split(ep, ",")))
			}
			palletMatch(cmd.Context(), networkName, websocketAddr, fls, cmd.Flag("output").Value.String())
		},
	}
	matchCmd.PersistentFlags().StringP("network", "w", "", "multiple separated by ',' \n eg: -w polkadot")
	matchCmd.PersistentFlags().StringP("pallet", "p", "", "Find supported pallets, multiple separated by ',' \n eg: -p System,Babe")
	matchCmd.PersistentFlags().StringP("exclude-pallet", "e", "", "Exclude supported pallets, multiple separated by ',' \n eg: -p System,Babe")
	matchCmd.PersistentFlags().StringP("output", "o", "", "output to file path")
	return matchCmd
}

func newCompareCmd() *cobra.Command {
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "subscan network comparison with substrate standard pallet",
		Long:  "subscan network comparison with substrate standard pallet",
		Run: func(cmd *cobra.Command, args []string) {
			nw := strings.TrimSpace(cmd.Flag("network").Value.String())
			var networkName []string
			var websocketAddr []string
			if len(nw) == 0 {
				// default network
				networkName = conf.Conf.Network
			} else {
				for _, n := range strings.Split(nw, ",") {
					if len(n) == 0 {
						continue
					}
					if strings.HasPrefix(n, "wss://") || strings.HasPrefix(n, "ws://") {
						websocketAddr = append(websocketAddr, n)
					} else {
						networkName = append(networkName, n)
					}
				}
			}

			networkComparePallet(cmd.Context(), networkName, websocketAddr, cmd.Flag("output").Value.String())
		},
	}
	compareCmd.PersistentFlags().StringP("network", "w", "", "network name or network websocket addr,multiple separated by ',' \n eg: -w polkadot")
	compareCmd.PersistentFlags().StringP("output", "o", "", "output to file path")
	return compareCmd
}

func networkComparePallet(ctx context.Context, networkName, websocketAddr []string, path string) {
	pallet, err := substrate.PalletList(ctx)
	if err != nil {
		log.Printf("failed to get substrate standard pallet. err: %s\n", err.Error())
		return
	}
	if len(pallet) == 0 {
		log.Println("get the substrate pallet is empty")
		return
	}
	var palletList []model.NetworkData[string]
	if len(networkName) != 0 {
		log.Printf("get subscan network pallet")
		palletList = append(palletList, subscan.NetworkPalletList(ctx, networkName)...)
	}
	if len(websocketAddr) != 0 {
		log.Printf("get websocket network pallet")
		palletList = append(palletList, ws.NetworkPalletList(ctx, wsAddrFormat(websocketAddr))...)
	}
	if len(palletList) == 0 {
		return
	}
	var instance output.FormatCompareCharter
	if path != "" {
		instance = output.NewFileOutput(path)
	} else {
		instance = output.NewStdout()
	}
	if err = instance.FormatCompareChart(pallet, palletList); err != nil {
		log.Printf("output err: %s", err.Error())
		return
	}
}

func palletMatch(ctx context.Context, networkName, websocketAddr []string, fls []filter.Filter, path string) {
	var palletList []model.NetworkData[string]
	if len(networkName) != 0 {
		log.Printf("get subscan network pallet")
		palletList = append(palletList, subscan.NetworkPalletList(ctx, networkName)...)
	}
	if len(websocketAddr) != 0 {
		log.Printf("get websocket network pallet")
		palletList = append(palletList, ws.NetworkPalletList(ctx, wsAddrFormat(websocketAddr))...)
	}

	if len(palletList) == 0 {
		return
	}
	var instance output.FormatCharter
	if path != "" {
		instance = output.NewFileOutput(path)
	} else {
		instance = output.NewStdout()
	}
	for _, f := range fls {
		palletList = f.FilterPallet(palletList)
	}
	if err := instance.FormatChart(palletList); err != nil {
		log.Printf("output err: %s", err.Error())
		return
	}
}

func wsAddrFormat(ws []string) map[string]string {
	socketAddr := make(map[string]string)
	for _, addr := range ws {
		if u, err := url.Parse(addr); err == nil {
			if len(u.Host) != 0 {
				socketAddr[u.Host] = addr
				continue
			}
		}
		socketAddr[addr] = addr
	}
	return socketAddr
}
