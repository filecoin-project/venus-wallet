package wallet

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/crypto"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"golang.org/x/xerrors"
)

type wallet struct {
	ws storage.KeyStore
}

func NewWallet(ks storage.KeyStore) api.IWallet {
	return &wallet{ws: ks}
}
func (w *wallet) WalletNew(ctx context.Context, kt core.KeyType) (core.Address, error) {
	prv, err := crypto.GeneratePrivateKey(core.KeyType2Sign(kt))
	if err != nil {
		return core.NilAddress, err
	}
	addr, err := prv.Address()
	if err != nil {
		return core.NilAddress, err
	}
	err = w.ws.Put(prv)
	if err != nil {
		return core.NilAddress, err
	}
	return addr, nil

}
func (w *wallet) WalletHas(ctx context.Context, address core.Address) (bool, error) {
	return w.ws.Has(address)
}
func (w *wallet) WalletList(ctx context.Context) ([]core.Address, error) {
	return w.ws.List()
}
func (w *wallet) WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) {
	var (
		owner core.Address
		err   error
		data  []byte
	)
	if meta.Type == core.MTChainMsg {
		msg, err := core.DecodeMessage(meta.Extra)
		if err != nil {
			return nil, err
		}
		owner = msg.From
		data = msg.Cid().Bytes()
	} else {
		_, toSign, err := core.GetSignBytes(toSign, meta)
		if err != nil {
			return nil, xerrors.Errorf("get sign bytes failed:%w", err)
		}
		owner = signer
		data = toSign
	}
	pk, err := w.ws.Get(owner)
	if err != nil {
		return nil, err
	}
	return pk.Sign(data)
}

func (w *wallet) WalletExport(ctx context.Context, addr core.Address) (*core.KeyInfo, error) {
	key, err := w.ws.Get(addr)
	if err != nil {
		return nil, err
	}
	return key.ToKeyInfo(), nil
}
func (w *wallet) WalletImport(ctx context.Context, ki *core.KeyInfo) (core.Address, error) {
	pk, err := crypto.NewKeyFromKeyInfo(ki)
	if err != nil {
		return core.NilAddress, err
	}
	addr, err := pk.Address()
	if err != nil {
		return core.NilAddress, err
	}
	err = w.ws.Put(pk)
	if err != nil {
		return core.NilAddress, err
	}
	return addr, nil
}
func (w *wallet) WalletDelete(ctx context.Context, addr core.Address) error {
	return w.ws.Delete(addr)
}

var _ api.IWallet = &wallet{}
