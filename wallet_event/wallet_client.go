package wallet_event

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/google/uuid"
	"github.com/ipfs-force-community/venus-gateway/types"
	"github.com/ipfs-force-community/venus-gateway/walletevent"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"net/http"
	"net/url"
)

type WalletRegisterClient struct {
	ResponseWalletEvent func(ctx context.Context, resp *types.ResponseEvent) error
	ListenWalletEvent   func(ctx context.Context, policy *walletevent.WalletRegisterPolicy) (chan *types.RequestEvent, error)
	SupportNewAccount   func(ctx context.Context, channelId uuid.UUID, account string) error
}

func NewWalletRegisterClient(ctx context.Context, url, token string) (*WalletRegisterClient, jsonrpc.ClientCloser, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	addr, err := dialArgs(url)
	if err != nil {
		return nil, nil, err
	}
	walletEventClient := &WalletRegisterClient{}
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin", []interface{}{walletEventClient}, headers)
	if err != nil {
		return nil, nil, err
	}
	return walletEventClient, closer, nil
}

func dialArgs(maddr string) (string, error) {
	ma, err := multiaddr.NewMultiaddr(maddr)
	if err == nil {
		_, addr, err := manet.DialArgs(ma)
		if err != nil {
			return "", err
		}

		return "ws://" + addr + "/rpc/v0", nil
	}

	_, err = url.Parse(maddr)
	if err != nil {
		return "", err
	}
	return maddr + "/rpc/v0", nil
}
