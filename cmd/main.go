package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"network-runtime-check/internal/match"
)

func main() {
	var pallet = flag.String("pallet", "", "Find supported pallets, multiple separated by ',' \n eg: -pallet=System,Babe")
	flag.Parse()

	var ctx, cancel = context.WithCancel(context.Background())
	go notify(cancel)
	match.NetworkPalletMatch(ctx, *pallet)
	cancel()
}

func notify(cancel context.CancelFunc) {
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Kill, os.Interrupt, syscall.SIGTERM)
	s := <-sign
	log.Printf("receive signal %s, exit...\n", s.String())
	cancel()
}
