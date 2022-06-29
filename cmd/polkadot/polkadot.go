package polkadot

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/subscan-explorer/network-runtime-check/internal/api/github/substrate"
	"github.com/subscan-explorer/network-runtime-check/internal/output"
	"github.com/subscan-explorer/network-runtime-check/internal/ws"
)

func NewPolkadotCmd() *cobra.Command {
	polkadotCmd := &cobra.Command{
		Use:   "polkadot",
		Short: "polkadot websocket metadata",
		Long:  "Obtain metadata related data through polkadot websocket",
	}
	palletCmd := &cobra.Command{Use: "pallet", Short: "substrate pallet"}
	palletCmd.AddCommand(newMatchCmd(), newCompareCmd())
	polkadotCmd.AddCommand(palletCmd)
	return polkadotCmd
}

func newMatchCmd() *cobra.Command {
	matchCmd := &cobra.Command{
		Use:   "match",
		Short: "Pallets supported by poladot runtime",
		Long:  "Check polkadot for all network supported pallets",
		Run: func(cmd *cobra.Command, args []string) {
			pe := strings.TrimSpace(cmd.Flag("pallet").Value.String())
			var palletList []string
			if len(pe) != 0 {
				palletList = strings.Split(pe, ",")
			}
			matchPallet(cmd.Context(), palletList, cmd.Flag("output").Value.String())
		},
	}
	matchCmd.PersistentFlags().StringP("pallet", "p", "", "Find supported pallets, multiple separated by ',' \n eg: -p System,Babe")
	matchCmd.PersistentFlags().StringP("output", "o", "", "output to file path")
	return matchCmd
}

func matchPallet(ctx context.Context, pallet []string, path string) {
	palletList := ws.NetworkPalletList(ctx)
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

func newCompareCmd() *cobra.Command {
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "polkadot network comparison with substrate standard pallet",
		Long:  "polkadot network comparison with substrate standard pallet",
		Run: func(cmd *cobra.Command, args []string) {
			comparePallet(cmd.Context(), cmd.Flag("output").Value.String())
		},
	}
	compareCmd.PersistentFlags().StringP("output", "o", "", "output to file path")
	return compareCmd
}

func comparePallet(ctx context.Context, path string) {
	pallet, err := substrate.PalletList(ctx)
	if err != nil {
		log.Printf("failed to get substrate standard pallet. err: %s\n", err.Error())
		return
	}
	if len(pallet) == 0 {
		log.Println("get the substrate pallet is empty")
		return
	}
	palletList := ws.NetworkPalletList(ctx)
	if len(palletList) == 0 {
		log.Println("get the network pallet is empty")
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
