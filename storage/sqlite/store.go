package sqlite

import (
	"database/sql"
	"encoding/json"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/venus-wallet/config"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"golang.org/x/xerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// keystore sqlite implementation
type sqliteStorage struct {
	db *gorm.DB
}

func NewSQLiteStorage(cfg *config.DBConfig) (storage.KeyStore, error) {
	var db, err = gorm.Open(sqlite.Open(cfg.Conn), &gorm.Config{})
	var sqldb *sql.DB
	if err != nil {
		return nil, xerrors.Errorf("open database(%s) failed:%w", cfg.Conn, err)
	}

	if sqldb, err = db.DB(); err != nil {
		return nil, xerrors.Errorf("sqlDb failed, %w", err)
	}

	sqldb.SetConnMaxIdleTime(300)
	sqldb.SetMaxIdleConns(8)
	sqldb.SetMaxOpenConns(64)
	//db = db.Debug()
	if !db.Migrator().HasTable(&Wallet{}) {
		if err = db.AutoMigrate(&Wallet{}); err != nil {
			return nil, xerrors.Errorf("migrate failed:%w", err)
		}
	}
	return &sqliteStorage{db: db}, err
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
	wallet := Wallet{
		Address: key.Address,
		KeyInfo: &ki,
	}
	if err = s.db.First(&wallet, "address=?", wallet.Address).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else if err == gorm.ErrRecordNotFound {
		return s.db.Create(&wallet).Error
	}
	return err
}

func (s *sqliteStorage) Has(addr core.Address) (bool, error) {
	var counts int64 = 0
	err := s.db.Table("wallets").Where("address=?", addr.String()).Count(&counts).Error
	if err != nil {
		return false, err
	}
	return counts > 0, err
}

func (s *sqliteStorage) List() ([]core.Address, error) {
	var ws []Wallet
	err := s.db.Table("wallets"). /*.Select("address", "nonce").*/ Scan(&ws).Error
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
	if err := s.db.Where("address=?", addr.String()).First(res).Error; err != nil {
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
	tmpDb := s.db.Table("wallets").Delete(nil, "address = ?", addr.String())
	if err = tmpDb.Error; err != nil {
		// may be it isn't explicit, but acceptable
		return xerrors.Errorf("delete wallet(%s) failed:%w",
			addr.String(), err)
	}
	return nil
}
