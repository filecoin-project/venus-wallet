package api

import (
	"github.com/ipfs-force-community/venus-wallet/storage"
)

var _ storage.StrategyStore = &StrategyAuth{}

type StrategyAuth struct {
	Internal struct {
		PutMsgTypeTemplate       func(mtt *storage.MsgTypeTemplate) error                         `perm:"admin" local:"required"`
		GetMsgTypeTemplate       func(mttId uint) (*storage.MsgTypeTemplate, error)               `perm:"admin" local:"required"`
		GetMsgTypeTemplateByName func(name string) (*storage.MsgTypeTemplate, error)              `perm:"admin" local:"required"`
		ListMsgTypeTemplates     func(fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) `perm:"admin" local:"required"`
		DeleteMsgTypeTemplate    func(mttId uint) error                                           `perm:"admin" local:"required"`

		PutMethodTemplate       func(mt *storage.MethodTemplate) error                          `perm:"admin" local:"required"`
		GetMethodTemplate       func(mtId uint) (*storage.MethodTemplate, error)                `perm:"admin" local:"required"`
		GetMethodTemplateByName func(name string) (*storage.MethodTemplate, error)              `perm:"admin" local:"required"`
		ListMethodTemplates     func(fromIndex, toIndex int) ([]*storage.MethodTemplate, error) `perm:"admin" local:"required"`
		DeleteMethodTemplate    func(mtId uint) error                                           `perm:"admin" local:"required"`

		PutKeyBind              func(kb *storage.KeyBind) error                          `perm:"admin" local:"required"`
		GetKeyBinds             func(address string) ([]*storage.KeyBind, error)         `perm:"admin" local:"required"`
		GetKeyBindByName        func(name string) (*storage.KeyBind, error)              `perm:"admin" local:"required"`
		GetKeyBindById          func(kbId uint) (*storage.KeyBind, error)                `perm:"admin" local:"required"`
		ListKeyBinds            func(fromIndex, toIndex int) ([]*storage.KeyBind, error) `perm:"admin" local:"required"`
		DeleteKeyBind           func(kbId uint) error                                    `perm:"admin" local:"required"`
		DeleteKeyBindsByAddress func(address string) (int64, error)                      `perm:"admin" local:"required"`

		PutGroup       func(name string, keyBindIds []uint) error             `perm:"admin" local:"required"`
		GetGroupByName func(name string) (*storage.Group, error)              `perm:"admin" local:"required"`
		GetGroup       func(gId uint) (*storage.Group, error)                 `perm:"admin" local:"required"`
		ListGroups     func(fromIndex, toIndex int) ([]*storage.Group, error) `perm:"admin" local:"required"`
		DeleteGroup    func(gId uint) error                                   `perm:"admin" local:"required"`

		PutGroupAuth    func(token string, groupId uint) error         `perm:"admin" local:"required"`
		GetGroupAuth    func(token string) (*storage.GroupAuth, error) `perm:"admin" local:"required"`
		DeleteGroupAuth func(token string) error                       `perm:"admin" local:"required"`
	}
}

func (s *StrategyAuth) PutMsgTypeTemplate(mtt *storage.MsgTypeTemplate) error {
	return s.Internal.PutMsgTypeTemplate(mtt)
}
func (s *StrategyAuth) GetMsgTypeTemplateByName(name string) (*storage.MsgTypeTemplate, error) {
	return s.Internal.GetMsgTypeTemplateByName(name)
}
func (s *StrategyAuth) GetMsgTypeTemplate(mttId uint) (*storage.MsgTypeTemplate, error) {
	return s.Internal.GetMsgTypeTemplate(mttId)
}
func (s *StrategyAuth) ListMsgTypeTemplates(fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) {
	return s.Internal.ListMsgTypeTemplates(fromIndex, toIndex)
}
func (s *StrategyAuth) DeleteMsgTypeTemplate(mttId uint) error {
	return s.Internal.DeleteMsgTypeTemplate(mttId)
}

func (s *StrategyAuth) PutMethodTemplate(mt *storage.MethodTemplate) error {
	return s.Internal.PutMethodTemplate(mt)
}
func (s *StrategyAuth) GetMethodTemplate(mtId uint) (*storage.MethodTemplate, error) {
	return s.Internal.GetMethodTemplate(mtId)
}
func (s *StrategyAuth) GetMethodTemplateByName(name string) (*storage.MethodTemplate, error) {
	return s.Internal.GetMethodTemplateByName(name)
}

func (s *StrategyAuth) ListMethodTemplates(fromIndex, toIndex int) ([]*storage.MethodTemplate, error) {
	return s.Internal.ListMethodTemplates(fromIndex, toIndex)
}
func (s *StrategyAuth) DeleteMethodTemplate(mtId uint) error {
	return s.Internal.DeleteMethodTemplate(mtId)
}

func (s *StrategyAuth) PutKeyBind(kb *storage.KeyBind) error {
	return s.Internal.PutKeyBind(kb)
}
func (s *StrategyAuth) GetKeyBinds(address string) ([]*storage.KeyBind, error) {
	return s.Internal.GetKeyBinds(address)
}
func (s *StrategyAuth) GetKeyBindByName(name string) (*storage.KeyBind, error) {
	return s.Internal.GetKeyBindByName(name)
}
func (s *StrategyAuth) GetKeyBindById(kbId uint) (*storage.KeyBind, error) {
	return s.Internal.GetKeyBindById(kbId)
}
func (s *StrategyAuth) ListKeyBinds(fromIndex, toIndex int) ([]*storage.KeyBind, error) {
	return s.Internal.ListKeyBinds(fromIndex, toIndex)
}

func (s *StrategyAuth) DeleteKeyBind(kbId uint) error {
	return s.Internal.DeleteKeyBind(kbId)
}
func (s *StrategyAuth) DeleteKeyBindsByAddress(address string) (int64, error) {
	return s.Internal.DeleteKeyBindsByAddress(address)
}

func (s *StrategyAuth) PutGroup(name string, keyBindIds []uint) error {
	return s.Internal.PutGroup(name, keyBindIds)
}

func (s *StrategyAuth) GetGroupByName(name string) (*storage.Group, error) {
	return s.Internal.GetGroupByName(name)
}

func (s *StrategyAuth) GetGroup(gId uint) (*storage.Group, error) {
	return s.Internal.GetGroup(gId)
}

func (s *StrategyAuth) ListGroups(fromIndex, toIndex int) ([]*storage.Group, error) {
	return s.Internal.ListGroups(fromIndex, toIndex)
}
func (s *StrategyAuth) DeleteGroup(gId uint) error {
	return s.Internal.DeleteGroup(gId)
}

func (s *StrategyAuth) PutGroupAuth(token string, groupId uint) error {
	return s.Internal.PutGroupAuth(token, groupId)
}
func (s *StrategyAuth) GetGroupAuth(token string) (*storage.GroupAuth, error) {
	return s.Internal.GetGroupAuth(token)
}
func (s *StrategyAuth) DeleteGroupAuth(token string) error {
	return s.Internal.DeleteGroupAuth(token)
}
