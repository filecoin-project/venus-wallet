package main

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/build"
	"github.com/ipfs-force-community/venus-wallet/filemgr"
	"github.com/ipfs-force-community/venus-wallet/middleware"
	"github.com/ipfs-force-community/venus-wallet/version"
	"github.com/mitchellh/go-homedir"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"io/ioutil"
	"os"
	"os/user"
	"testing"
)

func TestSetup(t *testing.T) {
	absoluteTmp, err := ioutil.TempDir("", "venus-wallet-")
	defer os.RemoveAll(absoluteTmp)
	if err != nil {
		t.Fatal(err)
	}
	ctx, _ := tag.New(context.Background(), tag.Insert(middleware.Version, version.BuildVersion))
	dir, err := homedir.Expand(absoluteTmp)
	if err != nil {
		t.Logf("could not expand repo location error:%s", err)
	} else {
		t.Logf("venus repo: %s", dir)
	}

	apiListen := "/ip4/0.0.0.0/tcp/5678"

	op := &filemgr.OverrideParams{
		API: apiListen,
	}
	// true to debug local
	if false {
		user, err := user.Current()
		if err != nil {
			t.Fatal(err)
		}
		absoluteTmp = user.HomeDir
	}
	r, err := filemgr.NewFS(absoluteTmp, op)
	if err != nil {
		t.Fatalf("opening fs repo: %s", err)
	}
	secret, err := r.APISecret()
	if err != nil {
		t.Fatalf("read secret failed: %s", err)
	}
	var fullAPI api.IFullAPI
	stop, err := build.New(ctx,
		build.Override(build.SetNet, func() {
			address.CurrentNetwork = address.Mainnet
		}),
		build.FullAPIOpt(&fullAPI),
		build.WalletOpt(r.Config()),
		build.CommonOpt(secret),
		build.Override(new(build.NetworkName), build.NetworkName("main net")),
	)
	if err != nil {
		t.Fatalf("initializing node: %s", err)
	}

	// Register all metric views
	if err = view.Register(
		middleware.DefaultViews...,
	); err != nil {
		t.Fatalf("Cannot register the view: %s", err)
	}

	// Set the metric to one so it is published to the exporter
	stats.Record(ctx, middleware.VenusInfo.M(1))

	endpoint, err := r.APIEndpoint()
	if err != nil {
		t.Fatalf("getting api endpoint: %s", err)
	}
	t.Log(endpoint.String(), stop)
	t.Log("Pre-preparation completed")

	// TODO: properly parse api endpoint (or make it a URL)
	// Use serveRPC method to perform local CLI debugging
	// ServeRPC(fullAPI, stop, endpoint)
}
