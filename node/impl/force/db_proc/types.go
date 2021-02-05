package db_proc

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"gorm.io/gorm"
)

type SqlScannerValuer interface {
	sql.Scanner
	driver.Valuer
}

type Key struct {
	types.KeyInfo

	PublicKey []byte
	Address   address.Address
}

type Wallet struct {
	gorm.Model
	Address string           `gorm:"type:varchar(255);uniqueIndex"`
	KeyInfo SqlKeyInfo       `gorm:"type:blob;column:private_key;not null"`
	Nonce   uint64           `gorm:"type:unsigned integer;default:0"`
	Meta    SqlScannerValuer `gorm:"type:blob;column:meta;default:null"`
}

type SignedMsg struct {
	gorm.Model
	Address   string            `gorm:"type:varchar(255);uniqueIndex:uidx_address_nonce"`
	Cid       string            `gorm:"type:varchar(128);uniqueIndex:uidx_cid"`
	SignedMsg SqlChainSignedMsg `gorm:"type:blob;column:signed_msg"`
	Nonce     uint64            `gorm:"type:unsigned integer;uniqueIndex:uidx_address_nonce;default:0"`
	Epoch     uint64            `gorm:"type:unsigned integer;column:epoch;not null"`
}

type SignedData struct {
	gorm.Model
	Signer string     `gorm:"type:varchar(255);index"`
	Data   SqlBytes   `grom:"type:blob"`
	Meta   SqlMsgMeta `gorm:"type:blob"`
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
		return fmt.Errorf("Failed to unmarshal JSONB value:", value)
	}
	err := json.Unmarshal(data, mki)
	return err
}

// Value return json value, implement driver.Valuer interface
func (mki SqlKeyInfo) Value() (driver.Value, error) {
	if !mki.IsValid() {
		return nil, fmt.Errorf("invalid key-info: (type:%d, value:0x%x)", mki.Type, mki.PrivateKey)
	}
	return json.Marshal(mki)
}

func (w Wallet) GetKey() (*Key, error) {
	return NewKey((types.KeyInfo)(w.KeyInfo))
}

type SqlBytes []byte

func (sb *SqlBytes) Scan(value interface{}) error {
	tmp, isok := value.([]byte)
	if !isok {
		return fmt.Errorf("value must be []byte")
	}
	*sb = tmp
	return nil
}

func (sb SqlBytes) Value() (driver.Value, error) {
	return []byte(sb), nil
}

// conflict on field 'Value' and to-implement interface function 'Value'
type SqlChainSignedMsg struct{ types.SignedMessage }

func (sc *SqlChainSignedMsg) Scan(value interface{}) error {
	data, isok := value.([]byte)
	if !isok {
		return fmt.Errorf("value must be []byte")
	}
	return (&sc.Message).UnmarshalCBOR(bytes.NewReader(data))
}

func (sc SqlChainSignedMsg) Value() (driver.Value, error) {
	msg, err := sc.Message.Serialize()
	return driver.Value(msg), err
}

type SqlMsgMeta api.MsgMeta

func (sc *SqlMsgMeta) Scan(value interface{}) error {
	data, isok := value.([]byte)
	if !isok {
		return fmt.Errorf("value must be []byte")
	}
	return json.Unmarshal(data, sc)
}

func (sc SqlMsgMeta) Value() (driver.Value, error) {
	data, err := json.Marshal(&sc)
	return driver.Value(data), err
}

var _ = SqlScannerValuer((*SqlMsgMeta)(nil))
var _ = SqlScannerValuer((*SqlKeyInfo)(nil))
var _ = SqlScannerValuer((*SqlKeyInfo)(nil))
var _ = SqlScannerValuer((*SqlBytes)(nil))
var _ = SqlScannerValuer((*SqlChainSignedMsg)(nil))

