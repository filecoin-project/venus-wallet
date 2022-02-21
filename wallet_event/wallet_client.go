package wallet_event

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	types2 "github.com/filecoin-project/venus/venus-shared/types"
	types "github.com/filecoin-project/venus/venus-shared/types/gateway"
	"github.com/ipfs-force-community/venus-common-utils/apiinfo"
)

type WalletRegisterClient struct {
	ResponseWalletEvent func(ctx context.Context, resp *types.ResponseEvent) error
	ListenWalletEvent   func(ctx context.Context, policy *types.WalletRegisterPolicy) (chan *types.RequestEvent, error)
	SupportNewAccount   func(ctx context.Context, channelId types2.UUID, account string) error
	AddNewAddress       func(ctx context.Context, channelId types2.UUID, newAddrs []address.Address) error
	RemoveAddress       func(ctx context.Context, channelId types2.UUID, newAddrs []address.Address) error
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
