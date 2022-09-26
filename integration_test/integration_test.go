//stm: #integration
package integration

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/venus-shared/types"

	"github.com/filecoin-project/go-jsonrpc"
	api2 "github.com/filecoin-project/venus-wallet/api"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/venus-wallet/cli/helper"

	"github.com/urfave/cli/v2"
)

var (
	inst             *WalletInst
	client           api2.IFullAPI
	clientCloser     jsonrpc.ClientCloser
	defaultWalletPwd = "default-wallet-pwd"
)

func setup() error {
	var err error
	var url string
	if inst, err = NewWalletInst(); err != nil {
		return err
	}
	if url, err = inst.Start(); err != nil {
		return err
	}

	log.Infof("Wallet instance listen on endpoint: %s", url)

	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.String("repo", inst.repoDir, "")

	client, clientCloser, err = helper.GetFullAPI(cli.NewContext(nil, flags, nil))
	return err
}

func shutDown() error {
	defer func() {
		if err := os.RemoveAll(inst.repoDir); err != nil {
			log.Errorf("remove repo dir:%s failed:%s", inst.repoDir, err.Error())
		}
	}()

	clientCloser()
	return inst.StopAndWait()
}

func TestMain(m *testing.M) {
	//stm: @VENUSWALLET_NODE_NEW_NODE_CLIENT_001
	if err := setup(); err != nil {
		panic(fmt.Sprintf("setup ingeration test failed:%v", err))
	}

	exitCode := m.Run()

	if err := shutDown(); err != nil {
		panic(fmt.Sprintf("shutdown ingeration test failed:%v", err))
	}

	os.Exit(exitCode)
}

func TestWallet(t *testing.T) {
	//stm: @VENUSWALLET_STORAGE_KEYMIX_SET_PASSWORD_001, @VENUSWALLET_STORAGE_WALLET_SET_PASSWORD_001
	t.Run("wallet setPwd/unlock", testWalletSetPassword)

	//stm: @VENUSWALLET_STORAGE_SQLITE_KEY_STORE_PUT_001, @VENUSWALLET_STORAGE_SQLITE_KEY_STORE_HAS_001,
	//stm: @VENUSWALLET_STORAGE_SQLITE_KEY_STORE_LIST_001,  @VENUSWALLET_STORAGE_SQLITE_KEY_STORE_DELETE_001, @VENUSWALLET_STORAGE_SQLITE_KEY_STORE_GET_001
	//stm: @VENUSWALLET_STORAGE_WALLET_WALLET_NEW_001, @VENUSWALLET_STORAGE_WALLET_WALLET_LIST_001
	t.Run("wallet address", testWalletAddress)
}

func testWalletSetPassword(t *testing.T) {
	ctx := context.TODO()
	if err := client.SetPassword(ctx, defaultWalletPwd); err != nil {
		require.Contains(t, err.Error(), "already have")
	}

	if err := client.Unlock(ctx, defaultWalletPwd); err != nil {
		require.Contains(t, err.Error(), "already unlock")
	}
}

func testWalletAddress(t *testing.T) {
	ctx := context.TODO()

	var newAddrs = make(map[address.Address]struct{})
	var err error

	blsAddr, err := client.WalletNew(ctx, types.KTBLS)
	require.NoError(t, err)
	newAddrs[blsAddr] = struct{}{}

	secpAddr, err := client.WalletNew(ctx, types.KTSecp256k1)
	require.NoError(t, err)
	newAddrs[secpAddr] = struct{}{}

	addrList, err := client.WalletList(ctx)
	require.NoError(t, err)

	var totalAddrs = make(map[address.Address]struct{})
	for _, a := range addrList {
		totalAddrs[a] = struct{}{}
	}

	for addr := range newAddrs {
		_, isok := totalAddrs[addr]
		require.True(t, isok, true)

		isok, err = client.WalletHas(ctx, addr)
		require.NoError(t, err)
		require.True(t, isok)
	}

	require.NoError(t, client.Lock(ctx, defaultWalletPwd))

	// call unlock for covering `store.sqlite.sqliteStorage.Get`
	require.NoError(t, client.Unlock(ctx, defaultWalletPwd))

	// for covering  `store.sqlite.sqliteStorage.Delete`
	deleteAddr, err := client.WalletNew(ctx, types.KTBLS)
	require.NoError(t, err)
	require.NoError(t, client.WalletDelete(ctx, deleteAddr))
	// after delete this `deleteAddr`, `WalletHas` should returns us a `false`
	has, err := client.WalletHas(ctx, deleteAddr)
	require.NoError(t, err)
	require.False(t, has)
}
