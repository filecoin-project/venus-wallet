package sqlite

import (
	"database/sql"
	"github.com/ipfs-force-community/venus-wallet/core"
	"time"
)

// GroupAuth relation with Group and generate a token for external invocation
// NOTE: routeType 1
type GroupAuth struct {
	// uuid
	Token   string `gorm:"primaryKey;type:varchar(36);column:token;not null"`
	GroupId uint   `gorm:"index;type:int;column:id;not null"`
	//Public    []byte //pubkey for verify data source
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}

// KeyBind  bind wallet usage strategy
// allow designated rule to pass
// NOTE: routeType 2
type KeyBind struct {
	BindId  uint   `gorm:"primaryKey;type:int;column:id;not null"`
	Name    string `gorm:"type:varchar(32);column:name;not null"`
	Address string `gorm:"index;type:varchar(255);column:address;not null"`
	// source from MsgTypeTemplate or temporary create
	MetaTypes core.MsgEnum `gorm:"type:int;column:meta_types;not null"`
	// source from MethodTemplate
	MethodNames string `gorm:"type:text;column:methods;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime `gorm:"index"`
}

// Group multi KeyBind
// NOTE: routeType 3
type Group struct {
	GroupId   uint   `gorm:"primaryKey;type:int;column:id;not null"`
	Name      string `gorm:"type:varchar(32);column:name;not null"`
	BindIds   string `gorm:"type:text;column:bind_ids;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}

// MethodTemplate to quickly create a private key usage strategy
// msg actor and methodNum agg to method name
// NOTE: routeType 4
type MethodTemplate struct {
	MTId uint   `gorm:"primaryKey;type:int;column:id;not null"`
	Name string `gorm:"type:varchar(32);column:name;not null"`
	// method name join with ','
	Methods   string `gorm:"type:text;column:methods;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}

// MsgTypeTemplate to quickly create a private key usage strategy
// NOTE: routeType 5
type MsgTypeTemplate struct {
	MTTId     uint         `gorm:"primaryKey;type:int;column:id;not null"`
	Name      string       `gorm:"type:varchar(36);column:name;not null"`
	MetaTypes core.MsgEnum `gorm:"type:int;column:meta_types;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}
