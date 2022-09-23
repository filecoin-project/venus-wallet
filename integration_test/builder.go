package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/filecoin-project/venus-wallet/api"
	"github.com/filecoin-project/venus-wallet/build"
	"github.com/filecoin-project/venus-wallet/cmd"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/filemgr"
	"github.com/filecoin-project/venus-wallet/middleware"
	"github.com/filecoin-project/venus-wallet/version"
	logging "github.com/ipfs/go-log/v2"
	"github.com/multiformats/go-multiaddr"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var log = logging.Logger("wallet_instance")

type WalletInst struct {
	repo    filemgr.Repo
	repoDir string

	stopChan chan error
	sigChan  chan os.Signal
}

func (inst *WalletInst) Start() (string, error) {
	secret, err := inst.repo.APISecret()
	if err != nil {
		return "", nil
	}

	var fullAPI api.IFullAPI
	var appStopFn build.StopFunc

	ctx, _ := tag.New(context.Background(), tag.Insert(middleware.Version, version.BuildVersion))

	if appStopFn, err = build.New(ctx, build.FullAPIOpt(&fullAPI),
		build.WalletOpt(inst.repo, ""),
		build.CommonOpt(secret)); err != nil {
		return "", err
	}

	// Register all metric views
	if err = view.Register(
		middleware.DefaultViews...,
	); err != nil {
		return "", fmt.Errorf("can't register the view: %v", err)
	}
	stats.Record(ctx, middleware.VenusInfo.M(1))

	endPoint, err := inst.repo.APIEndpoint()

	ma, err := multiaddr.NewMultiaddr(endPoint)
	if err != nil {
		return "", fmt.Errorf("new multi-address failed:%w", err)
	}

	url, err := ToURL(ma)
	if err != nil {
		return "", fmt.Errorf("convert multi-addr:%s to url failed:%w", endPoint, err)
	}

	go func() {
		err := cmd.ServeRPC(fullAPI, appStopFn, endPoint, inst.sigChan)
		inst.stopChan <- err
	}()

	return url.String(), inst.checkService()
}

func (inst *WalletInst) checkService() error {
	select {
	case err := <-inst.stopChan:
		return err
	case <-time.After(time.Second):
		log.Info("waiting for service shutdown for 1 seconds")
	}
	return nil
}

func (inst *WalletInst) StopAndWait() error {
	inst.sigChan <- syscall.SIGINT
	for {
		if err := inst.checkService(); err != nil {
			// server close is not an error
			if strings.ContainsAny(err.Error(), "server closed") {
				return nil
			}
			return err
		}
	}
}

func NewWalletInst() (*WalletInst, error) {
	dir, err := ioutil.TempDir("", "venus_wallet_")
	if err != nil {
		return nil, err

	}
	repo, err := filemgr.NewFS(dir, nil)
	if err != nil {
		return nil, err
	}
	core.WalletStrategyLevel = repo.Config().Strategy.Level
	return &WalletInst{
		repo:     repo,
		sigChan:  make(chan os.Signal, 1),
		stopChan: make(chan error, 1),
		repoDir:  dir}, nil
}
