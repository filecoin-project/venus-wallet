package wallet_event

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/google/uuid"
	"github.com/ipfs-force-community/venus-common-utils/apiinfo"
	"github.com/ipfs-force-community/venus-gateway/types"
	"github.com/ipfs-force-community/venus-gateway/walletevent"
)

type WalletRegisterClient struct {
	ResponseWalletEvent func(ctx context.Context, resp *types.ResponseEvent) error
	ListenWalletEvent   func(ctx context.Context, policy *walletevent.WalletRegisterPolicy) (chan *types.RequestEvent, error)
	SupportNewAccount   func(ctx context.Context, channelId uuid.UUID, account string) error
	AddNewAddress       func(ctx context.Context, channelId uuid.UUID, newAddrs []address.Address) error
	RemoveAddress       func(ctx context.Context, channelId uuid.UUID, newAddrs []address.Address) error
}

func NewWalletRegisterClient(ctx context.Context, url, token string) (*WalletRegisterClient, jsonrpc.ClientCloser, error) {
	apiInfo := apiinfo.NewAPIInfo(url, token)
	addr, err := apiInfo.DialArgs("v0")
	if err != nil {
		return nil, nil, err
	}
	walletEventClient := &WalletRegisterClient{}
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Gateway", []interface{}{walletEventClient}, apiInfo.AuthHeader())
	if err != nil {
		return nil, nil, err
	}
	return walletEventClient, closer, nil
}
