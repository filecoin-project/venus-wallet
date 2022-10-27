package storage

import (
	types "github.com/filecoin-project/venus/venus-shared/types/wallet"
)

type StrategyStore interface {
	// PutMsgTypeTemplate saves a msgType template
	PutMsgTypeTemplate(mtt *types.MsgTypeTemplate) error
	// GetMsgTypeTemplate gets a msgType template by Id
	GetMsgTypeTemplate(mttId uint) (*types.MsgTypeTemplate, error)
	// GetMsgTypeTemplateByName gets a mesType template by name
	GetMsgTypeTemplateByName(name string) (*types.MsgTypeTemplate, error)
	// ListMsgTypeTemplates gets a list of msgType templates sort by Id
	ListMsgTypeTemplates(fromIndex, toIndex int) ([]*types.MsgTypeTemplate, error)
	// DeleteMsgTypeTemplate  removes a msgType template from strategy store
	DeleteMsgTypeTemplate(mttId uint) error

	// PutMethodTemplate saves a method template
	PutMethodTemplate(mt *types.MethodTemplate) error
	// GetMethodTemplate gets a method template by Id
	GetMethodTemplate(mtId uint) (*types.MethodTemplate, error)
	// GetMethodTemplateByName gets a method template by name
	GetMethodTemplateByName(name string) (*types.MethodTemplate, error)
	// ListMethodTemplates gets a list of method templates sort by Id
	ListMethodTemplates(fromIndex, toIndex int) ([]*types.MethodTemplate, error)
	// DeleteMethodTemplate  removes a method template from strategy store
	DeleteMethodTemplate(mtId uint) error

	// UpdateKeyBindMetaTypes update the elements of MetaTypes and Methods
	UpdateKeyBindMetaTypes(kb *types.KeyBind) error
	// PutKeyBind saves a keyBind
	PutKeyBind(kb *types.KeyBind) error
	// GetKeyBinds gets a list of keyBinds which wallet is address
	GetKeyBinds(address string) ([]*types.KeyBind, error)
	// GetKeyBindByName gets a keyBind by name
	GetKeyBindByName(name string) (*types.KeyBind, error)
	// GetKeyBindByNames gets a list of keyBinds that contain names
	GetKeyBindByNames(names []string) ([]*types.KeyBind, error)
	// GetKeyBindById gets a keyBind by Id
	GetKeyBindById(kbId uint) (*types.KeyBind, error)
	// ListKeyBinds gets a list of keyBinds sort by Id
	ListKeyBinds(fromIndex, toIndex int) ([]*types.KeyBind, error)
	// DeleteKeyBind removes a keyBinds from strategy store
	DeleteKeyBind(kbId uint) error
	// DeleteKeyBindsByAddress  removes a keyBind from strategy store
	DeleteKeyBindsByAddress(address string) (int64, error)

	// PutGroup saves a group
	PutGroup(name string, keyBindIDs []uint) error
	// GetGroupByName gets a group by name
	GetGroupByName(name string) (*types.Group, error)
	// GetGroup gets a group by Id
	GetGroup(gId uint) (*types.Group, error)
	// ListGroups gets a list of groups sort by Id
	ListGroups(fromIndex, toIndex int) ([]*types.Group, error)
	// DeleteGroup  removes a group from strategy store
	DeleteGroup(gId uint) error

	// PutGroup saves a strategy token
	PutGroupAuth(token string, groupId uint) error
	// GetTokensByGroupId gets a strategy token by groupId
	GetTokensByGroupId(groupId uint) ([]string, error)
	// GetGroupAuth gets a group by strategy token
	GetGroupAuth(token string) (*types.GroupAuth, error)
	// GetGroupKeyBind gets an element of keyBinds in group by token and wallet address
	GetGroupKeyBind(token string, address string) (*types.KeyBind, error)
	// DeleteGroupAuth  removes a group auth token from strategy store
	DeleteGroupAuth(token string) error
}
