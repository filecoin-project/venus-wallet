package storage

import (
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/msgrouter"
)

type RouterStore interface {
	PutMsgTypeTemplate(mtt *MsgTypeTemplate) error
	GetMsgTypeTemplate(mttId uint) (*MsgTypeTemplate, error)
	GetMsgTypeTemplateByName(name string) ([]*MsgTypeTemplate, error)
	ListMsgTypeTemplates(fromIndex, toIndex int) ([]*MsgTypeTemplate, error)
	DeleteMsgTypeTemplate(mttId uint) error

	PutMethodTemplate(mt *MethodTemplate) error
	GetMethodTemplate(mtId uint) (*MethodTemplate, error)
	ListMethodTemplates(fromIndex, toIndex int) ([]*MethodTemplate, error)
	DeleteMethodTemplate(mtId uint) error

	PutKeyBind(kb *KeyBind) error
	GetKeyBind(address string) ([]*KeyBind, error)
	GetKeyBindById(kbId uint) (*KeyBind, error)
	ListKeyBinds(fromIndex, toIndex int) ([]*KeyBind, error)

	PutGroup(g *Group) error
	GetGroup(gId uint) (*Group, error)
	ListGroups(fromIndex, toIndex int) ([]*Group, error)
	DeleteGroup(gId uint) error

	PutGroupAuth(ga *GroupAuth) error
	GetGroupAuth(token string) (*GroupAuth, error)
	DeleteGroupAuth(token string) error

	GetKeyStrategy(token, address string) (*GroupAuth, error)
}

type KeyStrategy struct {
	Address   core.Address
	MetaTypes core.MsgEnum
	Methods   []msgrouter.MethodName
}

// GroupAuth relation with Group and generate a token for external invocation
type GroupAuth struct {
	Token   string
	BindIds []string
}

// KeyBind  bind wallet usage strategy
// allow designated rule to pass
type KeyBind struct {
	BindId  uint
	Address core.Address
	// source from MsgTypeTemplate or temporary create
	MetaTypes core.MsgEnum
	// source from MethodTemplate
	Methods []msgrouter.MethodName
}

// Group multi KeyBind
type Group struct {
	GroupId uint
	Name    string
	BindIds []string
}

// MethodTemplate to quickly create a private key usage strategy
// msg actor and methodNum agg to method name
// NOTE: routeType 4
type MethodTemplate struct {
	MTId uint
	// method name join with ','
	Methods []msgrouter.MethodName
}

// MsgTypeTemplate to quickly create a private key usage strategy
type MsgTypeTemplate struct {
	MTTId     uint
	Name      string
	MetaTypes core.MsgEnum
}
