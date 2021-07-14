package main

import (
	"context"
	"log"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/api"
	"github.com/filecoin-project/venus-wallet/build"
	"github.com/filecoin-project/venus-wallet/cmd"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/filemgr"
	"github.com/filecoin-project/venus-wallet/middleware"
	"github.com/filecoin-project/venus-wallet/version"
	"github.com/mitchellh/go-homedir"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// this package for debug
func main() {
	absoluteTmp := "~/.venus_wallet"
	ctx, _ := tag.New(context.Background(), tag.Insert(middleware.Version, version.BuildVersion))
	dir, err := homedir.Expand(absoluteTmp)
	if err != nil {
		log.Printf("could not expand repo location error:%s", err)
	} else {
		log.Printf("venus repo: %s", dir)
	}
	apiListen := "/ip4/0.0.0.0/tcp/5678"
	op := &filemgr.OverrideParams{
		API: apiListen,
	}
	r, err := filemgr.NewFS(absoluteTmp, op)
	if err != nil {
		log.Fatalf("opening fs repo: %s", err)
	}
	core.WalletStrategyLevel = r.Config().Strategy.Level
	secret, err := r.APISecret()
	if err != nil {
		log.Fatalf("read secret failed: %s", err)
	}
	var fullAPI api.IFullAPI
	stop, err := build.New(ctx,
		build.Override(build.SetNet, func() {
			address.CurrentNetwork = address.Testnet
		}),
		build.FullAPIOpt(&fullAPI),
		build.WalletOpt(r, ""),
		build.CommonOpt(secret),
		build.Override(new(build.NetworkName), build.NetworkName("main net")),
	)
	if err != nil {
		log.Fatalf("initializing node: %s", err)
	}

	// Register all metric views
	if err = view.Register(
		middleware.DefaultViews...,
	); err != nil {
		log.Fatalf("Cannot register the view: %s", err)
	}

	// Set the metric to one so it is published to the exporter
	stats.Record(ctx, middleware.VenusInfo.M(1))

	endpoint, err := r.APIEndpoint()
	if err != nil {
		log.Fatalf("getting api endpoint: %s", err)
	}

	log.Println(endpoint, stop)
	log.Println("Pre-preparation completed")
	// TODO: properly parse api endpoint (or make it a URL)
	// Use serveRPC method to perform local CLI debugging
	err = cmd.ServeRPC(fullAPI, stop, endpoint)
	if err != nil {
		log.Fatal(err)
	}
}
