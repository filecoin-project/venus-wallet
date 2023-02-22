package sqlite

import (
	"fmt"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/storage"
	"github.com/filecoin-project/venus/venus-shared/types"
	"gorm.io/gorm"
)

const MTUndefined types.MsgType = ""

type sqliteSignRecord struct {
	CreatedAt time.Time     `gorm:"index"`
	Type      types.MsgType `gorm:"index"`
	Signer    string        `gorm:"type:varchar(256);not null"`
	Err       string        `gorm:"type:varchar(256);default:null"`
	Msg       []byte        `gorm:"type:blob;default:null"`
}

func (s *sqliteSignRecord) TableName() string {
	return "sign_record"
}

func newFromSignRecord(record *storage.SignRecord) *sqliteSignRecord {
	ret := &sqliteSignRecord{
		CreatedAt: record.CreateAt,
		Type:      record.Type,
		Signer:    record.Signer.String(),
		Msg:       record.Msg,
	}
	if record.Err != nil {
		ret.Err = record.Err.Error()
	}
	return ret
}

func (s *sqliteSignRecord) toSignRecord() *storage.SignRecord {
	ret := &storage.SignRecord{
		ID:       s.getId(),
		CreateAt: s.CreatedAt,
		Type:     s.Type,
		Signer:   MustParseAddress(s.Signer),
		Err:      fmt.Errorf(s.Err),
		Msg:      s.Msg,
	}
	if s.Err == "" {
		ret.Err = nil
	}
	return ret
}

func (s *sqliteSignRecord) getId() string {
	s1 := s.Signer
	s2 := s.CreatedAt.UnixNano()
	return fmt.Sprintf("%s-%d", s1, s2)
}

func parseId(id string) (signer string, createAt time.Time, err error) {
	_, err = fmt.Sscanf(id, "%s-%d", &signer, &createAt)
	return signer, createAt, err
}

type SqliteRecorder struct {
	db *gorm.DB
}

func NewSqliteRecorder(db *gorm.DB) (storage.IRecorder, error) {
	err := db.AutoMigrate(&sqliteSignRecord{})
	if err != nil {
		return nil, fmt.Errorf("init sqlite_recorder: %w", err)
	}
	return &SqliteRecorder{db: db}, nil
}

func (s *SqliteRecorder) Record(record *storage.SignRecord) error {
	return s.db.Create(newFromSignRecord(record)).Error
}

func (s *SqliteRecorder) QueryRecord(params *storage.QueryParams) ([]storage.SignRecord, error) {
	var records []*sqliteSignRecord
	query := s.db

	if params.ID != "" {
		signer, createAt, err := parseId(params.ID)
		if err != nil {
			return nil, fmt.Errorf("parse id: %w", err)
		}
		query = query.Where("signer = ?", signer).Where("created_at = ?", createAt)
	} else {
		if params.Signer != address.Undef {
			query = query.Where("signer = ?", params.Signer)
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
		err := query.Order("created_at desc").Find(&records).Error
		if err != nil {
			return nil, err
		}
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
