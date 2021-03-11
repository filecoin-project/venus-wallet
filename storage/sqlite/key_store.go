package sqlite

import (
	"encoding/json"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

// keystore sqlite implementation
type sqliteStorage struct {
	db       *gorm.DB
	walletTB string
}

func NewKeyStore(conn *Conn) storage.KeyStore {
	return &sqliteStorage{db: conn.DB, walletTB: TBWallet}
}

func (s *sqliteStorage) Put(key *storage.EncryptedKey) error {
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
	if err = s.db.Table(s.walletTB).First(wallet, "address=?", wallet.Address).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else if err == gorm.ErrRecordNotFound {
		return s.db.Table(s.walletTB).Create(wallet).Error
	}
	return err
}

func (s *sqliteStorage) Has(addr core.Address) (bool, error) {
	var counts int64 = 0
	err := s.db.Table(s.walletTB).Where("address=?", addr.String()).Count(&counts).Error
	if err != nil {
		return false, err
	}
	return counts > 0, err
}

func (s *sqliteStorage) List() ([]core.Address, error) {
	var ws []Wallet
	err := s.db.Table(s.walletTB).Scan(&ws).Error
	if err != nil {
		return nil, err
	}
	addresses := make([]core.Address, len(ws))
	for idx, val := range ws {
		addresses[idx], _ = address.NewFromString(val.Address)
	}
	return addresses, err
}

func (s *sqliteStorage) Get(addr core.Address) (*storage.EncryptedKey, error) {
	res := &Wallet{}
	if err := s.db.Table(s.walletTB).Where("address=?", addr.String()).First(res).Error; err != nil {
		return nil, err
	}
	cj := new(storage.CryptoJSON)
	err := json.Unmarshal(res.KeyInfo.PrivateKey, cj)
	if err != nil {
		return nil, err
	}
	return &storage.EncryptedKey{
		Address: addr.String(),
		KeyType: res.KeyInfo.Type,
		Crypto:  cj,
	}, nil
}

func (s *sqliteStorage) Delete(addr core.Address) error {
	var err error
	tmpDb := s.db.Table(s.walletTB).Delete(nil, "address = ?", addr.String())
	if err = tmpDb.Error; err != nil {
		// may be it isn't explicit, but acceptable
		return xerrors.Errorf("delete wallet(%s) failed:%w",
			addr.String(), err)
	}
	return nil
}
