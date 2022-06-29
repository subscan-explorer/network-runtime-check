package subscan

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/api/github/substrate"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/output"
)

func NewSubscanCmd() *cobra.Command {
	subscanCmd := &cobra.Command{
		Use:   "subscan",
		Short: "subscan metadata",
		Long:  "Obtain metadata related data through subscan",
	}
	palletCmd := &cobra.Command{Use: "pallet", Short: "substrate pallet"}
	palletCmd.AddCommand(newCompareCmd(), newMatchCmd())
	subscanCmd.AddCommand(palletCmd)
	return subscanCmd
}

func newMatchCmd() *cobra.Command {
	matchCmd := &cobra.Command{
		Use:   "match",
		Short: "Pallets supported by subscan runtime",
		Long:  "Check subscan for all network supported pallets",
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

func newCompareCmd() *cobra.Command {
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "subscan network comparison with substrate standard pallet",
		Long:  "subscan network comparison with substrate standard pallet",
		Run: func(cmd *cobra.Command, args []string) {
			nw := strings.TrimSpace(cmd.Flag("network").Value.String())
			var networkNodes []string
			if len(nw) == 0 {
				// default network
				networkNodes = conf.Conf.Network
			} else {
				networkNodes = strings.Split(nw, ",")
			}

			networkComparePallet(cmd.Context(), networkNodes, cmd.Flag("output").Value.String())
		},
	}
	compareCmd.PersistentFlags().StringP("network", "w", "", "multiple separated by ',' \n eg: -w polkadot")
	compareCmd.PersistentFlags().StringP("output", "o", "", "output to file path")
	return compareCmd
}

func networkComparePallet(ctx context.Context, network []string, path string) {
	pallet, err := substrate.PalletList(ctx)
	if err != nil {
		log.Printf("failed to get substrate standard pallet. err: %s\n", err.Error())
		return
	}
	if len(pallet) == 0 {
		log.Println("get the substrate pallet is empty")
		return
	}
	palletList := subscan.NetworkPalletList(ctx, network)

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

func palletMatch(ctx context.Context, network, pallet []string, path string) {
	palletList := subscan.NetworkPalletList(ctx, network)
	if len(palletList) == 0 {
		return
	}
	var instance output.FormatCharter
	if path != "" {
		instance = output.NewFileOutput(path)
	} else {
		instance = output.NewStdout()
	}
	if err := instance.FormatChart(pallet, palletList); err != nil {
		log.Printf("output err: %s", err.Error())
		return
	}
}
