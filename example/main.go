package main

import (
	"context"
	"github.com/filecoin-project/venus-wallet/api/remotecli"
	"github.com/filecoin-project/venus-wallet/api/remotecli/httpparse"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

type RemoteWallet struct {
	wallet.IWallet
	Cancel func()
}

func SetupRemoteWallet(info string) (*RemoteWallet, error) {
	ai, err := httpparse.ParseApiInfo(info)
	if err != nil {
		return nil, err
	}
	url, err := ai.DialArgs()
	if err != nil {
		return nil, err
	}
	wapi, closer, err := remotecli.NewWalletRPC(context.Background(), url, ai.AuthHeader())
	if err != nil {
		return nil, xerrors.Errorf("creating jsonrpc client: %w", err)
	}
	return &RemoteWallet{
		IWallet: wapi,
		Cancel:  closer,
	}, nil
}

func (w *RemoteWallet) Get() wallet.IWallet {
	if w == nil {
		return nil
	}
	return w
}

// How to access remote wallet
func main() {
	// env to prepare
	// mock production environment
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		os.Remove(path.Join(dir, "example", "pid.tmp"))
		os.Remove(path.Join(dir, "example", "remote-token.tmp"))
	}()

	cmd := exec.Cmd{
		Path:   path.Join(dir, "example", "wallet-setup.sh"),
		Args:   []string{"./wallet-setup.sh", "."},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	_ = cmd.Start()
	_ = cmd.Wait()

	tb, err := ioutil.ReadFile(path.Join(dir, "example", "remote-token.tmp"))
	if err != nil {
		log.Fatal(err)
	}
	token := strings.TrimSpace(string(tb))
	pb, err := ioutil.ReadFile(path.Join(dir, "example", "pid.tmp"))
	if err != nil {
		log.Fatal(err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(pb)))
	if err != nil {
		log.Fatal(err)
	}

	// remote wallet setup
	remoteWallet, err := SetupRemoteWallet(token)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := remoteWallet.WalletNew(context.Background(), core.KTSecp256k1)
	if err != nil {
		log.Fatalf("remote wallet new address error:%s", err)
	}
	log.Println("new address ", addr.String())
	addrs, err := remoteWallet.WalletList(context.Background())
	if err != nil {
		log.Fatalf("remote wallet list addresses error:%s", err)
	}
	for _, v := range addrs {
		log.Println(v.String())
	}
	exist, err := remoteWallet.WalletHas(context.Background(), addr)
	if err != nil {
		log.Fatalf("remote wallet check address exist error:%s", err)
	}
	log.Printf("addr:%s exist:%v", addr.String(), exist)
	remoteWallet.Cancel()
	err = syscall.Kill(pid, 9)
	if err != nil {
		log.Fatal(err)
	}
}
