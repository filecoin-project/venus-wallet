package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/filecoin-project/venus-wallet/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// for sqlite use a single file socket
type Conn struct {
	DB *gorm.DB
}

type TableName = string

const (
	TBWallet TableName = "wallets"
)

func NewSQLiteConn(cfg *config.DBConfig) (*Conn, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Conn), &gorm.Config{})
	var sqldb *sql.DB
	if err != nil {
		return nil, fmt.Errorf("open database(%s) failed:%w", cfg.Conn, err)
	}

	if sqldb, err = db.DB(); err != nil {
		return nil, fmt.Errorf("sqlDb failed, %w", err)
	}

	sqldb.SetConnMaxIdleTime(300)
	sqldb.SetMaxIdleConns(8)
	sqldb.SetMaxOpenConns(64)
	// key_types 1
	if !db.Migrator().HasTable(TBWallet) {
		if err = db.Table(TBWallet).AutoMigrate(&Wallet{}); err != nil {
			return nil, fmt.Errorf("migrate failed:%w", err)
		}
	}

	return &Conn{DB: db}, nil
}
