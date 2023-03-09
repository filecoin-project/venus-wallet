package wallet

import (
	"context"
	"fmt"
	"sync"

	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/google/uuid"

	"github.com/asaskevich/EventBus"
	wallet_api "github.com/filecoin-project/venus/venus-shared/api/wallet"
	"github.com/filecoin-project/venus/venus-shared/types"
	logging "github.com/ipfs/go-log/v2"

	c "github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/venus-wallet/crypto"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
	"github.com/filecoin-project/venus-wallet/storage"
)

var log = logging.Logger("wallet")

type GetPwdFunc func() string

var _ wallet_api.IWallet = &wallet{}

// wallet implementation
type wallet struct {
	keyCache map[string]crypto.PrivateKey // simple key cache
	ws       storage.KeyStore             // key storage
	mw       storage.KeyMiddleware        //
	bus      EventBus.Bus
	filter   ISignMsgFilter
	m        sync.RWMutex
	recorder storage.IRecorder
}

func NewWallet(ks storage.KeyStore, rd storage.IRecorder, mw storage.KeyMiddleware, filter ISignMsgFilter, bus EventBus.Bus, getPwd GetPwdFunc) wallet_api.ILocalWallet {
	w := &wallet{
		ws:       ks,
		recorder: rd,
		mw:       mw,
		bus:      bus,
		filter:   filter,
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
	if err := w.checkPassword(ctx, password); err != nil {
		return err
	}
	return w.mw.Unlock(ctx, password)
}

func (w *wallet) Lock(ctx context.Context, password string) error {
	return w.mw.Lock(ctx, password)
}

func (w *wallet) LockState(ctx context.Context) bool {
	return w.mw.LockState(ctx)
}

func (w *wallet) WalletNew(ctx context.Context, kt types.KeyType) (address.Address, error) {
	if err := w.mw.Next(); err != nil {
		return address.Undef, err
	}
	err := w.mw.CheckToken(ctx)
	if err != nil {
		return address.Undef, err
	}
	prv, err := crypto.GeneratePrivateKey(types.KeyType2Sign(kt))
	if err != nil {
		return address.Undef, err
	}
	addr, err := prv.Address()
	if err != nil {
		return address.Undef, err
	}
	ckey, err := w.mw.Encrypt(storage.EmptyPassword, prv)
	if err != nil {
		return address.Undef, err
	}
	err = w.ws.Put(ckey)
	if err != nil {
		return address.Undef, err
	}
	// notify
	w.bus.Publish("wallet:add_address", addr)
	return addr, nil
}

func (w *wallet) WalletHas(ctx context.Context, address address.Address) (bool, error) {
	return w.ws.Has(address)
}

func (w *wallet) WalletList(ctx context.Context) ([]address.Address, error) {
	addrs, err := w.ws.List()
	if err != nil {
		return nil, err
	}
	return addrs, nil
}

func (w *wallet) WalletSign(ctx context.Context, signer address.Address, data []byte, meta types.MsgMeta) (*c.Signature, error) {
	if err := w.mw.Next(); err != nil {
		return nil, err
	}

	// parse msg
	signObj, toSign, err := GetSignBytesAndObj(data, meta)
	if err != nil {
		return nil, fmt.Errorf("get sign bytes: %w", err)
	}

	// check owner
	if meta.Type == types.MTChainMsg {
		if signer != signObj.(*types.Message).From {
			return nil, fmt.Errorf("signer(%s) is not msg sender(%s)", signer, signObj.(*types.Message).From)
		}

		// Use the data passed directly, because the message of f4 address is not signed for cid.
		// https://github.com/filecoin-project/venus/blob/master/venus-shared/actors/types/message.go#L228
		toSign = data
	}

	// check filter
	if meta.Type != types.MTVerifyAddress {
		err = w.filter.CheckSignMsg(ctx, SignMsg{
			SignType: meta.Type,
			Data:     signObj,
		})
		if err != nil {
			return nil, err
		}
	}

	// sign
	prvKey := w.cacheKey(signer)
	if prvKey == nil {
		key, err := w.ws.Get(signer)
		if err != nil {
			return nil, err
		}
		prvKey, err = w.mw.Decrypt(storage.EmptyPassword, key)
		if err != nil {
			return nil, err
		}
		w.pushCache(signer, prvKey)
	}
	signature, signErr := prvKey.Sign(toSign)

	// record
	go func() {
		msg, err := cborutil.Dump(signObj)
		if err != nil {
			log.Errorf("dump signObj failed %v", err)
		}

		err = w.recorder.Record(&storage.SignRecord{
			ID:     uuid.New().String(),
			Type:   meta.Type,
			Signer: signer,
			RawMsg: msg,
			Err:    signErr,
		})
		if err != nil {
			log.Errorf("record sign failed: %v", err)
		}
	}()

	return signature, signErr
}

func (w *wallet) WalletExport(ctx context.Context, addr address.Address) (*types.KeyInfo, error) {
	if err := w.mw.Next(); err != nil {
		return nil, err
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

func (w *wallet) WalletImport(ctx context.Context, ki *types.KeyInfo) (address.Address, error) {
	if err := w.mw.Next(); err != nil {
		return address.Undef, err
	}
	err := w.mw.CheckToken(ctx)
	if err != nil {
		return address.Undef, err
	}
	pk, err := crypto.NewKeyFromKeyInfo(ki)
	if err != nil {
		return address.Undef, err
	}
	addr, err := pk.Address()
	if err != nil {
		return address.Undef, err
	}
	exist, err := w.ws.Has(addr)
	if err != nil {
		return address.Undef, err
	}
	if exist {
		return addr, nil
	}
	key, err := w.mw.Encrypt(storage.EmptyPassword, pk)
	if err != nil {
		return address.Undef, err
	}
	err = w.ws.Put(key)
	if err != nil {
		return address.Undef, err
	}
	// notify
	w.bus.Publish("wallet:add_address", addr)
	return addr, nil
}

func (w *wallet) WalletDelete(ctx context.Context, addr address.Address) error {
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

func (w *wallet) pushCache(address address.Address, prv crypto.PrivateKey) {
	w.m.Lock()
	defer w.m.Unlock()
	w.keyCache[address.String()] = prv
}

func (w *wallet) pullCache(address address.Address) {
	w.m.Lock()
	defer w.m.Unlock()
	delete(w.keyCache, address.String())
}

func (w *wallet) cacheKey(address address.Address) crypto.PrivateKey {
	w.m.RLock()
	defer w.m.RUnlock()
	return w.keyCache[address.String()]
}
