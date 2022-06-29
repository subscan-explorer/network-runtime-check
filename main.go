package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/subscan-explorer/network-runtime-check/cmd/polkadot"
	"github.com/subscan-explorer/network-runtime-check/cmd/subscan"
	"github.com/subscan-explorer/network-runtime-check/conf"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		confPath    string
		rootCmd     = cobra.Command{Use: "network-runtime-check",
			CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true}}
	)
	go notify(cancel)
	rootCmd.PersistentFlags().StringVarP(&confPath, "config", "f", "conf/config.yaml", "configuration file path")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		conf.InitConf(cmd.Context(), cmd.Flag("config").Value.String())
	}

	rootCmd.AddCommand(subscan.NewSubscanCmd(), polkadot.NewPolkadotCmd())
	rootCmd.SetContext(ctx)
	_ = rootCmd.Execute()
	cancel()
}

func notify(cancel context.CancelFunc) {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Kill, os.Interrupt, syscall.SIGTERM)
	s := <-sign
	log.Printf("receive signal %s, exit...\n", s.String())
	cancel()
}
