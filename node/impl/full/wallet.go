package full

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ipfs-force-community/venus-wallet/node/impl/force/db_proc"
	api2 "github.com/filecoin-project/lotus/api"
	"go.uber.org/fx"

	"github.com/filecoin-project/lotus/chain/types"
)

type WalletAPI struct {
	fx.In

	Db db_proc.DbProcInterface
}

func (a *WalletAPI) WalletNew(ctx context.Context, keyType types.KeyType) (address.Address, error) {
	return a.Db.WalletPut(keyType)
}

func (a *WalletAPI) WalletHas(ctx context.Context, addr address.Address) (bool, error) {
	return a.Db.WalletHas(addr)
}

func (a *WalletAPI) WalletList(ctx context.Context) ([]address.Address, error) {
	return a.Db.WalletList()
}

func (a *WalletAPI) WalletSign(ctx context.Context, addr address.Address, toSign []byte, meta api2.MsgMeta) (*crypto.Signature, error) {
	return a.Db.WalletSign(ctx, addr, toSign, meta)
}

func (a *WalletAPI) WalletExport(ctx context.Context, addr address.Address) (*types.KeyInfo, error) {
	return a.Db.WalletExport(addr)
}

func (a *WalletAPI) WalletImport(ctx context.Context, keyInfo *types.KeyInfo) (address.Address, error) {
	return a.Db.WalletImport(keyInfo)
}

func (a *WalletAPI) WalletDelete(ctx context.Context, addr address.Address) error {
	_, err := a.Db.WalletDel(addr)
	return err
}
