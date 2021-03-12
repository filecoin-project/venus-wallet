package wallet

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/crypto"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"golang.org/x/xerrors"
)

type ILocalWallet interface {
	IWallet
	storage.IWalletLock
}

// remote wallet api
type IWallet interface {
	WalletNew(context.Context, core.KeyType) (core.Address, error)
	WalletHas(ctx context.Context, address core.Address) (bool, error)
	WalletList(ctx context.Context) ([]core.Address, error)
	WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error)
	WalletExport(ctx context.Context, addr core.Address) (*core.KeyInfo, error)
	WalletImport(context.Context, *core.KeyInfo) (core.Address, error)
	WalletDelete(context.Context, core.Address) error
}

var _ IWallet = &wallet{}

// wallet implementation
type wallet struct {
	ws storage.KeyStore
	mw storage.KeyMiddleware
}

func NewWallet(ks storage.KeyStore, mw storage.KeyMiddleware) ILocalWallet {
	return &wallet{ws: ks, mw: mw}
}
func (w *wallet) SetPassword(ctx context.Context, password string) error {
	return w.mw.SetPassword(ctx, password)
}
func (w *wallet) Unlock(ctx context.Context, password string) error {
	return w.mw.Unlock(ctx, password)
}
func (w *wallet) Lock(ctx context.Context, password string) error {
	return w.mw.Lock(ctx, password)
}

func (w *wallet) WalletNew(ctx context.Context, kt core.KeyType) (core.Address, error) {
	if err := w.mw.Next(); err != nil {
		return core.NilAddress, err
	}
	prv, err := crypto.GeneratePrivateKey(core.KeyType2Sign(kt))
	if err != nil {
		return core.NilAddress, err
	}
	addr, err := prv.Address()
	if err != nil {
		return core.NilAddress, err
	}
	ckey, err := w.mw.Encrypt(prv)
	if err != nil {
		return core.NilAddress, err
	}
	err = w.ws.Put(ckey)
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
	if err := w.mw.Next(); err != nil {
		return nil, err
	}
	var (
		owner core.Address
		err   error
		data  []byte
	)
	if meta.Type == core.MTChainMsg {
		if len(meta.Extra) == 0 {
			return nil, xerrors.New("msg type must contain extra data")
		}
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
	key, err := w.ws.Get(owner)
	if err != nil {
		return nil, err
	}
	pkey, err := w.mw.Decrypt(key)
	if err != nil {
		return nil, err
	}
	return pkey.Sign(data)
}

func (w *wallet) WalletExport(ctx context.Context, addr core.Address) (*core.KeyInfo, error) {
	if err := w.mw.Next(); err != nil {
		return nil, err
	}
	key, err := w.ws.Get(addr)
	if err != nil {
		return nil, err
	}
	pkey, err := w.mw.Decrypt(key)
	if err != nil {
		return nil, err
	}
	return pkey.ToKeyInfo(), nil
}
func (w *wallet) WalletImport(ctx context.Context, ki *core.KeyInfo) (core.Address, error) {
	if err := w.mw.Next(); err != nil {
		return core.NilAddress, err
	}
	pk, err := crypto.NewKeyFromKeyInfo(ki)
	if err != nil {
		return core.NilAddress, err
	}
	addr, err := pk.Address()
	if err != nil {
		return core.NilAddress, err
	}
	exist, err := w.ws.Has(addr)
	if err != nil {
		return core.NilAddress, err
	}
	if exist {
		return addr, nil
	}
	key, err := w.mw.Encrypt(pk)
	if err != nil {
		return core.NilAddress, err
	}
	err = w.ws.Put(key)
	if err != nil {
		return core.NilAddress, err
	}
	return addr, nil
}
func (w *wallet) WalletDelete(ctx context.Context, addr core.Address) error {
	if err := w.mw.Next(); err != nil {
		return err
	}
	return w.ws.Delete(addr)
}
