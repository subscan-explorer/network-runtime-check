package balances

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/output"
	"github.com/subscan-explorer/network-runtime-check/internal/ws"
)

func NewParamCmd() *cobra.Command {
	matchCmd := &cobra.Command{
		Use:   "param",
		Short: "Check the extrinsic and event in the pallet",
		Long:  "Check whether the extrinsic and event in the pallet conform to the parameter definition",
		Run: func(cmd *cobra.Command, args []string) {
			events := conf.Conf.ParamRule
			if len(events) == 0 {
				log.Fatalln("The rule list in the configuration file is required")
			}
			var out output.FormatEventCharter = output.NewStdout()
			if path := cmd.Flag("output").Value.String(); len(path) > 0 {
				out = output.NewFileOutput(path)
			}
			paramCheck(cmd.Context(), events, out)
		},
	}
	matchCmd.PersistentFlags().StringP("output", "o", "", "output to file path")
	return matchCmd
}

func paramCheck(ctx context.Context, node []conf.ParamRule, charter output.FormatEventCharter) {
	var (
		domainNode []string
		wsNode     []conf.ParamRule
	)

	for _, e := range node {
		if strings.HasPrefix(e.WsAddr, "ws://") || strings.HasPrefix(e.WsAddr, "wss://") {
			wsNode = append(wsNode, e)
			continue
		}
		if len(e.Domain) != 0 {
			domainNode = append(domainNode, strings.ToLower(e.Domain))
		}
	}

	metadata := subscan.NetworkMetadataList(ctx, domainNode)
	metadata = append(metadata, ws.GetMetadata(ctx, wsNode)...)

	if err := charter.FormatEventChart(metadata, node); err != nil {
		log.Printf("output err: %s", err.Error())
		return
	}
}
