package storage

import (
	linq "github.com/ahmetb/go-linq/v3"

	"github.com/filecoin-project/venus-wallet/core"
)

type StrategyStore interface {
	//PutMsgTypeTemplate saves a msgType template
	PutMsgTypeTemplate(mtt *MsgTypeTemplate) error
	// GetMsgTypeTemplate gets a msgType template by Id
	GetMsgTypeTemplate(mttId uint) (*MsgTypeTemplate, error)
	// GetMsgTypeTemplateByName gets a mesType template by name
	GetMsgTypeTemplateByName(name string) (*MsgTypeTemplate, error)
	// ListMsgTypeTemplates gets a list of msgType templates sort by Id
	ListMsgTypeTemplates(fromIndex, toIndex int) ([]*MsgTypeTemplate, error)
	// DeleteMsgTypeTemplate  removes a msgType template from strategy store
	DeleteMsgTypeTemplate(mttId uint) error

	// PutMethodTemplate saves a method template
	PutMethodTemplate(mt *MethodTemplate) error
	// GetMethodTemplate gets a method template by Id
	GetMethodTemplate(mtId uint) (*MethodTemplate, error)
	// GetMethodTemplateByName gets a method template by name
	GetMethodTemplateByName(name string) (*MethodTemplate, error)
	// ListMethodTemplates gets a list of method templates sort by Id
	ListMethodTemplates(fromIndex, toIndex int) ([]*MethodTemplate, error)
	// DeleteMethodTemplate  removes a method template from strategy store
	DeleteMethodTemplate(mtId uint) error

	// UpdateKeyBindMetaTypes update the elements of MetaTypes and Methods
	UpdateKeyBindMetaTypes(kb *KeyBind) error
	// PutKeyBind saves a keyBind
	PutKeyBind(kb *KeyBind) error
	// GetKeyBinds gets a list of keyBinds which wallet is address
	GetKeyBinds(address string) ([]*KeyBind, error)
	// GetKeyBindByName gets a keyBind by name
	GetKeyBindByName(name string) (*KeyBind, error)
	// GetKeyBindByNames gets a list of keyBinds that contain names
	GetKeyBindByNames(names []string) ([]*KeyBind, error)
	// GetKeyBindById gets a keyBind by Id
	GetKeyBindById(kbId uint) (*KeyBind, error)
	// ListKeyBinds gets a list of keyBinds sort by Id
	ListKeyBinds(fromIndex, toIndex int) ([]*KeyBind, error)
	// DeleteKeyBind removes a keyBinds from strategy store
	DeleteKeyBind(kbId uint) error
	// DeleteKeyBindsByAddress  removes a keyBind from strategy store
	DeleteKeyBindsByAddress(address string) (int64, error)

	// PutGroup saves a group
	PutGroup(name string, keyBindIds []uint) error
	// GetGroupByName gets a group by name
	GetGroupByName(name string) (*Group, error)
	// GetGroup gets a group by Id
	GetGroup(gId uint) (*Group, error)
	// ListGroups gets a list of groups sort by Id
	ListGroups(fromIndex, toIndex int) ([]*Group, error)
	// DeleteGroup  removes a group from strategy store
	DeleteGroup(gId uint) error

	// PutGroup saves a strategy token
	PutGroupAuth(token string, groupId uint) error
	// GetTokensByGroupId gets a strategy token by groupId
	GetTokensByGroupId(groupId uint) ([]string, error)
	// GetGroupAuth gets a group by strategy token
	GetGroupAuth(token string) (*GroupAuth, error)
	// GetGroupKeyBind gets an element of keyBinds in group by token and wallet address
	GetGroupKeyBind(token string, address string) (*KeyBind, error)
	// DeleteGroupAuth  removes a group auth token from strategy store
	DeleteGroupAuth(token string) error
}

// KeyStrategy a uint of wallet strategy
type KeyStrategy struct {
	Address   core.Address      // wallet address
	MetaTypes core.MsgEnum      // sum MsgEnum
	Methods   []core.MethodName // msg method array
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
