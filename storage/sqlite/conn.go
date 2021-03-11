package sqlite

import (
	"database/sql"
	"github.com/ipfs-force-community/venus-wallet/config"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log = logging.Logger("sqlite")

// for sqlite use a single file socket
type Conn struct {
	DB *gorm.DB
}

type TableName = string

const (
	TBWallet          TableName = "wallets"
	TBGroupAuth       TableName = "groupAuth"
	TBKeyBind         TableName = "keyBind"
	TBGroup           TableName = "group"
	TBMethodTemplate  TableName = "methodTmp"
	TBMsgTypeTemplate TableName = "typeTmp"
)

func NewSQLiteConn(cfg *config.DBConfig) (*Conn, error) {
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
	// key_types 1
	if !db.Migrator().HasTable(TBWallet) {
		if err = db.Table(TBWallet).AutoMigrate(&Wallet{}); err != nil {
			return nil, xerrors.Errorf("migrate failed:%w", err)
		}
	}

	// NOTE: routeType 1
	if !db.Migrator().HasTable(TBGroupAuth) {
		if err = db.Table(TBGroupAuth).AutoMigrate(&GroupAuth{}); err != nil {
			return nil, xerrors.Errorf("migrate failed:%w", err)
		}
	}
	// NOTE: routeType 2
	if !db.Migrator().HasTable(TBKeyBind) {
		if err = db.Table(TBKeyBind).AutoMigrate(&KeyBind{}); err != nil {
			return nil, xerrors.Errorf("migrate failed:%w", err)
		}
	}
	// NOTE: routeType 3
	if !db.Migrator().HasTable(TBGroup) {
		if err = db.Table(TBGroup).AutoMigrate(&Group{}); err != nil {
			return nil, xerrors.Errorf("migrate failed:%w", err)
		}
	}
	// NOTE: routeType 4
	if !db.Migrator().HasTable(TBMethodTemplate) {
		if err = db.Table(TBMethodTemplate).AutoMigrate(&MethodTemplate{}); err != nil {
			return nil, xerrors.Errorf("migrate failed:%w", err)
		}
	}
	// NOTE: routeType 5
	if !db.Migrator().HasTable(TBMsgTypeTemplate) {
		if err = db.Table(TBMsgTypeTemplate).AutoMigrate(&MsgTypeTemplate{}); err != nil {
			return nil, xerrors.Errorf("migrate failed:%w", err)
		}
	}
	return &Conn{DB: db}, nil
}
