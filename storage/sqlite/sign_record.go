package sqlite

import (
	"fmt"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/storage"
	"github.com/filecoin-project/venus/venus-shared/types"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/gorm"
)

const MTUndefined types.MsgType = ""

var log = logging.Logger("recorder")

type sqliteSignRecord struct {
	ID        string    `gorm:"primaryKey;type:varchar(256);not null"`
	CreatedAt time.Time `gorm:"index"`
	Type      types.MsgType
	Signer    string            `gorm:"type:varchar(256);index;not null"`
	Err       string            `gorm:"type:varchar(256);default:null"`
	RawMsg    []byte            `gorm:"type:blob;default:null"`
	Signature *crypto.Signature `gorm:"embedded;embeddedPrefix:signature_"`
}

func (s *sqliteSignRecord) TableName() string {
	return "sign_record"
}

func newFromSignRecord(record *storage.SignRecord) *sqliteSignRecord {
	ret := &sqliteSignRecord{
		ID:        record.ID,
		CreatedAt: record.CreateAt,
		Type:      record.Type,
		Signer:    record.Signer.String(),
		RawMsg:    record.RawMsg,
		Signature: record.Signature,
	}
	if record.Err != nil {
		ret.Err = record.Err.Error()
	}
	return ret
}

func (s *sqliteSignRecord) toSignRecord() *storage.SignRecord {
	ret := &storage.SignRecord{
		ID:        s.ID,
		CreateAt:  s.CreatedAt,
		Type:      s.Type,
		Signer:    MustParseAddress(s.Signer),
		Err:       fmt.Errorf(s.Err),
		RawMsg:    s.RawMsg,
		Signature: s.Signature,
	}
	if s.Err == "" {
		ret.Err = nil
	}
	return ret
}

type SqliteRecorder struct {
	db *gorm.DB
}

func NewSqliteRecorder(db *gorm.DB, cfg *config.SignRecorderConfig) (storage.IRecorder, error) {
	enable, keepDuration := true, time.Hour*7*24
	if cfg != nil {
		enable = cfg.Enable
		if cfg.KeepDuration != "" {
			d, err := time.ParseDuration(cfg.KeepDuration)
			if err != nil {
				return nil, fmt.Errorf("init sqlite_recorder: %w", err)
			}
			keepDuration = d
		}
	}

	if !enable {
		return &RecorderStub{}, nil
	}

	err := db.AutoMigrate(&sqliteSignRecord{})
	if err != nil {
		return nil, fmt.Errorf("init sqlite_recorder: %w", err)
	}

	go func() {
		ticker := time.NewTicker(time.Hour)
		for {
			<-ticker.C
			err := db.Where("created_at < ?", time.Now().Add(-keepDuration)).Delete(&sqliteSignRecord{}).Error
			if err != nil {
				log.Errorf("clean sqlite recorder: %s", err)
			}
		}
	}()

	return &SqliteRecorder{db: db}, nil
}

func (s *SqliteRecorder) Record(record *storage.SignRecord) error {
	return s.db.Create(newFromSignRecord(record)).Error
}

func (s *SqliteRecorder) QueryRecord(params *storage.QueryParams) ([]storage.SignRecord, error) {
	var records []*sqliteSignRecord
	query := s.db

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	} else {
		if params.Signer != address.Undef {
			query = query.Where("signer = ?", params.Signer.String())
		}
		if !params.After.IsZero() {
			query = query.Where("created_at >= ?", params.After)
		}
		if !params.Before.IsZero() {
			query = query.Where("created_at <= ?", params.Before)
		}
		if params.IsError {
			query = query.Where("err is not null")
		}
		if params.Type != MTUndefined {
			query = query.Where("type = ?", params.Type)
		}
		if params.Skip > 0 {
			query = query.Offset(params.Skip)
		}
		if params.Limit > 0 {
			query = query.Limit(params.Limit)
		}
	}

	err := query.Order("created_at desc").Find(&records).Error
	if err != nil {
		return nil, err
	}

	ret := make([]storage.SignRecord, 0, len(records))
	for _, r := range records {
		ret = append(ret, *r.toSignRecord())
	}
	return ret, nil
}

func MustParseAddress(addr string) address.Address {
	a, err := address.NewFromString(addr)
	if err != nil {
		panic(err)
	}
	return a
}

type RecorderStub struct {
}

func (r *RecorderStub) Record(record *storage.SignRecord) error {
	return nil
}

func (r *RecorderStub) QueryRecord(params *storage.QueryParams) ([]storage.SignRecord, error) {
	return nil, nil
}
