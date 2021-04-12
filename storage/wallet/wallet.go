package wallet

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/crypto"
	"github.com/ipfs-force-community/venus-wallet/errcode"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"github.com/ipfs-force-community/venus-wallet/storage/strategy"
	"golang.org/x/xerrors"
	"sync"
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
	keyCache map[string]crypto.PrivateKey
	ws       storage.KeyStore
	mw       storage.KeyMiddleware
	verify   strategy.IStrategyVerify
	m        sync.RWMutex
}

func NewWallet(ks storage.KeyStore, mw storage.KeyMiddleware, verify strategy.ILocalStrategy) ILocalWallet {
	return &wallet{
		ws:       ks,
		mw:       mw,
		verify:   verify,
		keyCache: make(map[string]crypto.PrivateKey),
	}
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
func (w *wallet) LockState(ctx context.Context) bool {
	return w.mw.LockState(ctx)
}
func (w *wallet) WalletNew(ctx context.Context, kt core.KeyType) (core.Address, error) {
	if err := w.mw.Next(); err != nil {
		return core.NilAddress, err
	}
	err := w.mw.CheckToken(ctx)
	if err != nil {
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
	if !w.verify.ContainWallet(ctx, address) {
		return false, errcode.ErrWithoutPermission
	}
	return w.ws.Has(address)
}

func (w *wallet) WalletList(ctx context.Context) ([]core.Address, error) {
	addrScope, err := w.verify.ScopeWallet(ctx)
	if err != nil {
		return nil, err
	}
	if addrScope.Root {
		return w.ws.List()
	}
	if len(addrScope.Addresses) == 0 {
		return addrScope.Addresses, nil
	}
	addrs, err := w.ws.List()
	if err != nil {
		return nil, err
	}
	linq.From(addrs).Intersect(linq.From(addrScope.Addresses)).ToSlice(&addrs)
	return addrs, nil
}

func (w *wallet) WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) {
	if err := w.mw.Next(); err != nil {
		return nil, err
	}
	var (
		owner core.Address
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
		if signer.String() != owner.String() {
			return nil, xerrors.New("singer does not match from in MSG")
		}
		data = msg.Cid().Bytes()
		if err = w.verify.Verify(ctx, signer, meta.Type, msg); err != nil {
			return nil, err
		}
	} else {
		_, toSign, err := core.GetSignBytes(toSign, meta)
		if err != nil {
			return nil, xerrors.Errorf("get sign bytes failed:%w", err)
		}
		owner = signer
		data = toSign
		if err = w.verify.Verify(ctx, signer, meta.Type, nil); err != nil {
			return nil, err
		}
	}
	prvKey := w.cacheKey(owner)
	if prvKey == nil {
		key, err := w.ws.Get(owner)
		if err != nil {
			return nil, err
		}
		prvKey, err = w.mw.Decrypt(key)
		if err != nil {
			return nil, err
		}
		w.pushCache(owner, prvKey)
	}
	return prvKey.Sign(data)
}

func (w *wallet) WalletExport(ctx context.Context, addr core.Address) (*core.KeyInfo, error) {
	if err := w.mw.Next(); err != nil {
		return nil, err
	}
	if !w.verify.ContainWallet(ctx, addr) {
		return nil, errcode.ErrWithoutPermission
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
	err := w.mw.CheckToken(ctx)
	if err != nil {
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
	if !w.verify.ContainWallet(ctx, addr) {
		return errcode.ErrWithoutPermission
	}
	err := w.ws.Delete(addr)
	if err != nil {
		return err
	}
	w.pullCache(addr)
	return nil
}

func (w *wallet) pushCache(address core.Address, prv crypto.PrivateKey) {
	w.m.Lock()
	defer w.m.Unlock()
	w.keyCache[address.String()] = prv
}

func (w *wallet) pullCache(address core.Address) {
	w.m.Lock()
	defer w.m.Unlock()
	delete(w.keyCache, address.String())
}
func (w *wallet) cacheKey(address core.Address) crypto.PrivateKey {
	w.m.RLock()
	defer w.m.RUnlock()
	return w.keyCache[address.String()]
}
