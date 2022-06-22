package cmd

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

func NewCompare() *cobra.Command {
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

	var instance output.Output
	if path != "" {
		instance = output.NewFileOutput(path)
	} else {
		instance = output.NewStdout()
	}
	if err = instance.Output(pallet, palletList); err != nil {
		log.Printf("output err: %s", err.Error())
		return
	}

}
