package sqlite

import (
	"encoding/json"
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
	"github.com/filecoin-project/venus-wallet/storage"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/gorm"
)

var ksLog = logging.Logger("key_store")

// keystore sqlite implementation
type sqliteStorage struct {
	db       *gorm.DB
	walletTB string
}

func NewKeyStore(db *gorm.DB) storage.KeyStore {
	store := &sqliteStorage{db: db, walletTB: TBWallet}
	_ = store.migrateCompatibleAddress()
	return store
}

func (s *sqliteStorage) Put(key *aes.EncryptedKey) error {
	keyBytes, err := json.Marshal(key.Crypto)
	if err != nil {
		return err
	}
	ki := SqlKeyInfo{
		Type:       key.KeyType,
		PrivateKey: keyBytes,
	}
	wallet := &Wallet{
		Address: key.Address,
		KeyInfo: &ki,
	}
	var sqlAddr shortAddress
	if sqlAddr, err = shortAddressFromString(key.Address); err != nil {
		return fmt.Errorf("%s is not an address:%w", key.Address, err)
	}

	if err = s.db.Table(s.walletTB).First(wallet, "address=?",
		sqlAddr).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else if err == gorm.ErrRecordNotFound {
		wallet.Address = sqlAddr.String()
		return s.db.Table(s.walletTB).Create(wallet).Error
	}
	return err
}

func (s *sqliteStorage) Has(addr address.Address) (bool, error) {
	var counts int64 = 0
	err := s.db.Table(s.walletTB).Where("address=?", shortAddress(addr)).Count(&counts).Error
	if err != nil {
		return false, err
	}
	return counts > 0, err
}

func (s *sqliteStorage) List() ([]address.Address, error) {
	var ws []Wallet
	err := s.db.Table(s.walletTB).Scan(&ws).Error
	if err != nil {
		return nil, err
	}
	addresses := make([]address.Address, 0, len(ws))
	for _, val := range ws {
		addr, err := shortAddressFromString(val.Address)
		if err != nil {
			ksLog.Error("can't decode:%s to address:%s", val.Address, err.Error())
			continue
		}
		addresses = append(addresses, addr.Address())
	}
	return addresses, err
}

func (s *sqliteStorage) Get(addr address.Address) (*aes.EncryptedKey, error) {
	res := &Wallet{}
	if err := s.db.Table(s.walletTB).Where("address=?", shortAddress(addr)).First(res).Error; err != nil {
		return nil, err
	}
	cj := new(aes.CryptoJSON)
	err := json.Unmarshal(res.KeyInfo.PrivateKey, cj)
	if err != nil {
		return nil, err
	}
	return &aes.EncryptedKey{
		Address: addr.String(),
		KeyType: res.KeyInfo.Type,
		Crypto:  cj,
	}, nil
}

func (s *sqliteStorage) Delete(addr address.Address) error {
	var err error
	tmpDb := s.db.Table(s.walletTB).Delete(nil, "address = ?", shortAddress(addr))
	if err = tmpDb.Error; err != nil {
		// may be it isn't explicit, but acceptable
		return fmt.Errorf("delete wallet(%s) failed:%w",
			addr.String(), err)
	}
	return nil
}

func (s *sqliteStorage) migrateCompatibleAddress() error {
	var ws []Wallet
	err := s.db.Table(s.walletTB).Scan(&ws).Error
	if err != nil {
		return err
	}

	for _, w := range ws {
		if addr, err := address.NewFromString(w.Address); err == nil {
			if err = s.db.Table(s.walletTB).Where("address = ?", w.Address).
				Update("address", shortAddress(addr)).Error; err != nil {
				return fmt.Errorf("migrate address:%s failed:%w", addr.String(), err)
			}
		}
	}
	return nil
}
