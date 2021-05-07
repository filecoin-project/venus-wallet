package storage

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/filecoin-project/venus-wallet/core"
)

type StrategyStore interface {
	PutMsgTypeTemplate(mtt *MsgTypeTemplate) error
	GetMsgTypeTemplate(mttId uint) (*MsgTypeTemplate, error)
	GetMsgTypeTemplateByName(name string) (*MsgTypeTemplate, error)
	ListMsgTypeTemplates(fromIndex, toIndex int) ([]*MsgTypeTemplate, error)
	DeleteMsgTypeTemplate(mttId uint) error

	PutMethodTemplate(mt *MethodTemplate) error
	GetMethodTemplate(mtId uint) (*MethodTemplate, error)
	GetMethodTemplateByName(name string) (*MethodTemplate, error)
	ListMethodTemplates(fromIndex, toIndex int) ([]*MethodTemplate, error)
	DeleteMethodTemplate(mtId uint) error

	UpdateKeyBindMetaTypes(kb *KeyBind) error
	PutKeyBind(kb *KeyBind) error
	GetKeyBinds(address string) ([]*KeyBind, error)
	GetKeyBindByName(name string) (*KeyBind, error)
	GetKeyBindByNames(names []string) ([]*KeyBind, error)
	GetKeyBindById(kbId uint) (*KeyBind, error)
	ListKeyBinds(fromIndex, toIndex int) ([]*KeyBind, error)
	DeleteKeyBind(kbId uint) error
	DeleteKeyBindsByAddress(address string) (int64, error)

	PutGroup(name string, keyBindIds []uint) error
	GetGroupByName(name string) (*Group, error)
	GetGroup(gId uint) (*Group, error)
	ListGroups(fromIndex, toIndex int) ([]*Group, error)
	DeleteGroup(gId uint) error

	PutGroupAuth(token string, groupId uint) error
	GetTokensByGroupId(groupId uint) ([]string, error)
	GetGroupAuth(token string) (*GroupAuth, error)
	GetGroupKeyBind(token string, address string) (*KeyBind, error)
	DeleteGroupAuth(token string) error
}

type KeyStrategy struct {
	Address   core.Address
	MetaTypes core.MsgEnum
	Methods   []core.MethodName
}

// GroupAuth relation with Group and generate a token for external invocation
type GroupAuth struct {
	Token    string
	GroupId  uint
	Name     string
	KeyBinds []*KeyBind
}

// KeyBind  bind wallet usage strategy
// allow designated rule to pass
type KeyBind struct {
	BindId  uint
	Name    string
	Address string
	// source from MsgTypeTemplate or temporary create
	MetaTypes core.MsgEnum
	// source from MethodTemplate
	Methods []core.MethodName
}

func (kb *KeyBind) ContainMsgType(mt core.MsgType) bool {
	return core.ContainMsgType(kb.MetaTypes, mt)
}
func (kb *KeyBind) ContainMethod(m string) bool {
	return linq.From(kb.Methods).Contains(m)
}

// Group multi KeyBind
type Group struct {
	GroupId uint
	Name    string
	// NOTE: not fill data when query groups
	KeyBinds []*KeyBind
}

// MethodTemplate to quickly create a private key usage strategy
// msg actor and methodNum agg to method name
// NOTE: routeType 4
type MethodTemplate struct {
	MTId uint
	Name string
	// method name join with ','
	Methods []core.MethodName
}

// MsgTypeTemplate to quickly create a private key usage strategy
type MsgTypeTemplate struct {
	MTTId     uint
	Name      string
	MetaTypes core.MsgEnum
}
