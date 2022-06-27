package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/spf13/cobra"
	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

func NewMatch() *cobra.Command {
	matchCmd := &cobra.Command{
		Use:   "match",
		Short: "Pallets supported by subscan network runtime",
		Long:  "Check if the runtime of all supported networks subscan supports a pallet",
		Run: func(cmd *cobra.Command, args []string) {
			nw := strings.TrimSpace(cmd.Flag("network").Value.String())
			var networkNodes []string
			if len(nw) == 0 {
				// default network
				networkNodes = conf.Conf.Network
			} else {
				networkNodes = strings.Split(nw, ",")
			}
			pe := strings.TrimSpace(cmd.Flag("pallet").Value.String())
			var palletList []string
			if len(pe) != 0 {
				palletList = strings.Split(pe, ",")
			}
			palletMatch(cmd.Context(), networkNodes, palletList, cmd.Flag("output").Value.String())
		},
	}
	matchCmd.PersistentFlags().StringP("network", "w", "", "multiple separated by ',' \n eg: -w polkadot")
	matchCmd.PersistentFlags().StringP("pallet", "p", "", "Find supported pallets, multiple separated by ',' \n eg: -p System,Babe")
	matchCmd.PersistentFlags().StringP("output", "o", "", "output to file path")
	return matchCmd
}

func palletMatch(ctx context.Context, network, pallet []string, output string) {
	palletList := subscan.NetworkPalletList(ctx, network)
	if len(palletList) == 0 {
		return
	}
	var palletSet map[string]struct{}
	if len(pallet) != 0 {
		palletSet = make(map[string]struct{})
		for _, p := range pallet {
			palletSet[strings.ToLower(p)] = struct{}{}
		}
	}
	var tableData [][]string
	for _, np := range palletList {
		var support []string
		support = append(support, np.Network)
		if np.Err != nil {
			support = append(support, utils.ErrorReduction(np.Err))
			continue
		}
		if palletSet != nil {
			var ps = make([]string, 0, len(pallet))
			for _, p := range np.Pallet {
				if _, ok := palletSet[strings.ToLower(p)]; ok {
					ps = append(ps, p)
				}
			}
			support = append(support, strings.Join(ps, "  "))
		} else {
			support = append(support, strings.Join(np.Pallet, "  "))
		}
		tableData = append(tableData, support)
	}
	if len(tableData) == 0 {
		return
	}
	tb := gotabulate.Create(tableData)
	tb.SetHeaders([]string{"Network", "Pallet"})
	tb.SetAlign("left")
	if output != "" {
		_ = os.WriteFile(output, []byte(tb.Render("grid")), os.FileMode(0644))
		return
	}
	// adaptive window size
	width := utils.TerminalWidth() - utils.MaxLenArrString(network) - 15
	if width < 0 {
		fmt.Println(tb.Render("grid"))
		return
	}
	tb.SetMaxCellSize(width)
	tb.SetWrapStrings(true)
	fmt.Println(tb.Render("grid"))
}
