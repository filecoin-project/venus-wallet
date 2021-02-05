package db_proc

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	mtype "github.com/ipfs-force-community/venus-wallet/chain/types"
	"github.com/ipfs-force-community/venus-wallet/node/config"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type MockDbProc struct{}

func (m MockDbProc) WalletPut(keyType types.KeyType) (address.Address, error) {
	return address.Undef, nil
}

func (m MockDbProc) WalletHas(a address.Address) (bool, error) { return false, nil }

func (m MockDbProc) WalletList() ([]address.Address, error) { return nil, nil }

func (m MockDbProc) WalletNonceSet(a address.Address, u uint64) (uint64, error) { return 0, nil }

func (m MockDbProc) WalletQuery(addr address.Address) (*Key, error) { return nil, nil }

func (m MockDbProc) WalletExport(addr address.Address) (*types.KeyInfo, error) { return nil, nil }

func (m MockDbProc) WalletImport(info *types.KeyInfo) (address.Address, error) {
	return address.Undef, nil
}

func (m MockDbProc) WalletSign(ctx context.Context, a address.Address, bytes []byte, meta api.MsgMeta) (*crypto.Signature, error) {
	return nil, nil
}

func (m MockDbProc) WalletDel(a address.Address) (bool, error) { return false, nil }

func (m MockDbProc) MessageAdd(ctx context.Context, u uint64, message *types.Message) (*types.SignedMessage, error) {
	return nil, nil
}

func (m MockDbProc) MessageResult(c cid.Cid, u uint64) error { return nil }

func (m MockDbProc) ChainState(uint64s []uint64, u uint64, a address.Address) ([]SignedMsg, error) {
	return nil, nil
}

func (m MockDbProc) MessageDel(a address.Address, u uint64, u2 uint64) (bool, error) {
	return false, nil
}

func (m MockDbProc) MessageQuery(a address.Address, u uint64, u2 uint64) ([]mtype.SignedMsg, error) {
	return nil, nil
}

func (m MockDbProc) UnchainMessageQuery(addr address.Address) ([]cid.Cid, error) {
	return nil, nil
}

func (m MockDbProc) QueryNoPubMsg() ([]SignedMsg, error) { return nil, nil }

var _ DbProcInterface = (*MockDbProc)(nil)

func NewMockDbProc(cfg *config.DbCfg) (DbProcInterface, error) {
	loger.Info("init nil db...")
	return &MockDbProc{}, nil
}
