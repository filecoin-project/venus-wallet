package wallet

import (
	"context"
	"sync"

	"github.com/ahmetb/go-linq/v3"
	"github.com/asaskevich/EventBus"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/crypto"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
	"github.com/filecoin-project/venus-wallet/errcode"
	"github.com/filecoin-project/venus-wallet/storage"
	"github.com/filecoin-project/venus-wallet/storage/strategy"
)

var log = logging.Logger("wallet")

type ILocalWallet interface {
	IWallet
	storage.IWalletLock
}

// IWallet remote wallet api
type IWallet interface {
	WalletNew(ctx context.Context, kt core.KeyType) (core.Address, error)
	WalletHas(ctx context.Context, address core.Address) (bool, error)
	WalletList(ctx context.Context) ([]core.Address, error)
	WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error)
	WalletExport(ctx context.Context, addr core.Address) (*core.KeyInfo, error)
	WalletImport(ctx context.Context, ki *core.KeyInfo) (core.Address, error)
	WalletDelete(ctx context.Context, addr core.Address) error
}

type GetPwdFunc func() string

var _ IWallet = &wallet{}

// wallet implementation
type wallet struct {
	keyCache map[string]crypto.PrivateKey // simple key cache
	ws       storage.KeyStore             // key storage
	mw       storage.KeyMiddleware        //
	verify   strategy.IStrategyVerify     // check wallet strategy with token
	bus      EventBus.Bus
	m        sync.RWMutex
}

func NewWallet(ks storage.KeyStore, mw storage.KeyMiddleware, bus EventBus.Bus, verify strategy.ILocalStrategy, getPwd GetPwdFunc) ILocalWallet {
	w := &wallet{
		ws:       ks,
		mw:       mw,
		verify:   verify,
		bus:      bus,
		keyCache: make(map[string]crypto.PrivateKey),
	}
	if getPwd != nil {
		if pwd := getPwd(); len(pwd) != 0 {
			if err := w.SetPassword(context.Background(), pwd); err != nil {
				log.Fatalf("set password(%s) failed %v", pwd, err)
			}
		}
	}

	return w
}
func (w *wallet) SetPassword(ctx context.Context, password string) error {
	if err := w.checkPassword(ctx, password); err != nil {
		return err
	}
	return w.mw.SetPassword(ctx, password)
}
func (w *wallet) checkPassword(ctx context.Context, password string) error {
	hashPasswd := aes.Keccak256([]byte(password))
	addrs, err := w.WalletList(ctx)
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		key, err := w.ws.Get(addr)
		if err != nil {
			return err
		}
		_, err = w.mw.Decrypt(hashPasswd, key)
		if err != nil {
			return err
		}
	}

	return nil
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
	ckey, err := w.mw.Encrypt(storage.EmptyPassword, prv)
	if err != nil {
		return core.NilAddress, err
	}
	err = w.ws.Put(ckey)
	if err != nil {
		return core.NilAddress, err
	}
	//notify
	w.bus.Publish("wallet:add_address", addr)
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
	// Do not validate strategy
	if meta.Type == core.MTVerifyAddress {
		_, toSign, err := core.GetSignBytes(toSign, meta)
		if err != nil {
			return nil, xerrors.Errorf("get sign bytes failed: %v", err)
		}
		owner = signer
		data = toSign
	} else if meta.Type == core.MTChainMsg {
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
			return nil, xerrors.Errorf("get sign bytes failed: %w", err)
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
		prvKey, err = w.mw.Decrypt(storage.EmptyPassword, key)
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
	pkey, err := w.mw.Decrypt(storage.EmptyPassword, key)
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
	key, err := w.mw.Encrypt(storage.EmptyPassword, pk)
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
	err := w.mw.CheckToken(ctx)
	if err != nil {
		return err
	}
	err = w.ws.Delete(addr)
	if err != nil {
		return err
	}
	w.pullCache(addr)
	w.bus.Publish("wallet:remove_address", addr)
	return nil
}

func (w *wallet) VerifyPassword(ctx context.Context, password string) error {
	if err := w.mw.Next(); err != nil {
		return err
	}
	return w.mw.VerifyPassword(ctx, password)
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
