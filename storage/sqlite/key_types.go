package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus/venus-shared/types"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type SqlScannerValuer interface {
	sql.Scanner
	driver.Valuer
}

type Key struct {
	types.KeyInfo
	PublicKey []byte
	Address   core.Address
}

type Wallet struct {
	gorm.Model
	Address string            `gorm:"type:varchar(255);uniqueIndex"`
	KeyInfo *SqlKeyInfo       `gorm:"type:blob;column:private_key;not null"`
	Meta    *SqlScannerValuer `gorm:"type:blob;column:meta;default:null"`
}

func (w *Wallet) TableName() string {
	return TBWallet
}

type SqlKeyInfo types.KeyInfo

func (mki *SqlKeyInfo) IsValid() bool {
	return mki != nil &&
		(mki.Type == types.KTBLS ||
			mki.Type == types.KTSecp256k1) &&
		len(mki.PrivateKey) != 0
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (mki *SqlKeyInfo) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return xerrors.Errorf("Failed to unmarshal JSONB value:%v", value)
	}
	err := json.Unmarshal(data, mki)
	return err
}

// Value return json value, implement driver.Valuer interface
func (mki SqlKeyInfo) Value() (driver.Value, error) {
	if !mki.IsValid() {
		return nil, xerrors.Errorf("invalid key-info: (type:%s, value:0x%x)", mki.Type, mki.PrivateKey)
	}
	return json.Marshal(mki)
}
