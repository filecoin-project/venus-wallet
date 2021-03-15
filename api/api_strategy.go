package api

import (
	"github.com/ipfs-force-community/venus-wallet/msgrouter"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"github.com/ipfs-force-community/venus-wallet/storage/strategy"
)

var _ strategy.IStrategy = &StrategyAuth{}

type StrategyAuth struct {
	Internal struct {
		NewMsgTypeTemplate     func(name string, codes []int) error                                          `perm:"admin" local:"required"`
		NewMethodTemplate      func(name string, methods []string) error                                     `perm:"admin" local:"required"`
		NewKeyBindCustom       func(name, address string, codes []int, methods []msgrouter.MethodName) error `perm:"admin" local:"required"`
		NewKeyBindFromTemplate func(name, address, mttName, mtName string) error                             `perm:"admin" local:"required"`
		NewGroup               func(name string, keyBindNames []string) error                                `perm:"admin" local:"required"`

		GetMsgTypeTemplate      func(name string) (*storage.MsgTypeTemplate, error) `perm:"admin" local:"required"`
		GetMethodTemplateByName func(name string) (*storage.MethodTemplate, error)  `perm:"admin" local:"required"`
		GetKeyBindByName        func(name string) (*storage.KeyBind, error)         `perm:"admin" local:"required"`
		GetKeyBinds             func(address string) ([]*storage.KeyBind, error)    `perm:"admin" local:"required"`
		GetGroupByName          func(name string) (*storage.Group, error)           `perm:"admin" local:"required"`

		ListGroups           func(fromIndex, toIndex int) ([]*storage.Group, error)           `perm:"admin" local:"required"`
		ListKeyBinds         func(fromIndex, toIndex int) ([]*storage.KeyBind, error)         `perm:"admin" local:"required"`
		ListMethodTemplates  func(fromIndex, toIndex int) ([]*storage.MethodTemplate, error)  `perm:"admin" local:"required"`
		ListMsgTypeTemplates func(fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) `perm:"admin" local:"required"`

		RemoveKeyBind          func(name string) error             `perm:"admin" local:"required"`
		RemoveKeyBindByAddress func(address string) (int64, error) `perm:"admin" local:"required"`
		RemoveGroup            func(name string) error             `perm:"admin" local:"required"`
		RemoveMethodTemplate   func(name string) error             `perm:"admin" local:"required"`
		RemoveMsgTypeTemplate  func(name string) error             `perm:"admin" local:"required"`
	}
}

func (o *StrategyAuth) NewMsgTypeTemplate(name string, codes []int) error {
	return o.Internal.NewMsgTypeTemplate(name, codes)
}
func (o *StrategyAuth) NewMethodTemplate(name string, methods []string) error {
	return o.Internal.NewMethodTemplate(name, methods)
}
func (o *StrategyAuth) NewKeyBindCustom(name, address string, codes []int, methods []msgrouter.MethodName) error {
	return o.Internal.NewKeyBindCustom(name, address, codes, methods)
}
func (o *StrategyAuth) NewKeyBindFromTemplate(name, address, mttName, mtName string) error {
	return o.Internal.NewKeyBindFromTemplate(name, address, mttName, mtName)
}
func (o *StrategyAuth) NewGroup(name string, keyBindNames []string) error {
	return o.Internal.NewGroup(name, keyBindNames)
}

func (o *StrategyAuth) GetMsgTypeTemplate(name string) (*storage.MsgTypeTemplate, error) {
	return o.Internal.GetMsgTypeTemplate(name)
}
func (o *StrategyAuth) GetMethodTemplateByName(name string) (*storage.MethodTemplate, error) {
	return o.Internal.GetMethodTemplateByName(name)
}
func (o *StrategyAuth) GetKeyBindByName(name string) (*storage.KeyBind, error) {
	return o.Internal.GetKeyBindByName(name)
}
func (o *StrategyAuth) GetKeyBinds(address string) ([]*storage.KeyBind, error) {
	return o.Internal.GetKeyBinds(address)
}
func (o *StrategyAuth) GetGroupByName(name string) (*storage.Group, error) {
	return o.Internal.GetGroupByName(name)
}

func (o *StrategyAuth) ListGroups(fromIndex, toIndex int) ([]*storage.Group, error) {
	return o.Internal.ListGroups(fromIndex, toIndex)
}
func (o *StrategyAuth) ListKeyBinds(fromIndex, toIndex int) ([]*storage.KeyBind, error) {
	return o.Internal.ListKeyBinds(fromIndex, toIndex)
}
func (o *StrategyAuth) ListMethodTemplates(fromIndex, toIndex int) ([]*storage.MethodTemplate, error) {
	return o.Internal.ListMethodTemplates(fromIndex, toIndex)
}
func (o *StrategyAuth) ListMsgTypeTemplates(fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) {
	return o.Internal.ListMsgTypeTemplates(fromIndex, toIndex)
}

func (o *StrategyAuth) RemoveKeyBind(name string) error {
	return o.Internal.RemoveKeyBind(name)
}
func (o *StrategyAuth) RemoveKeyBindByAddress(address string) (int64, error) {
	return o.Internal.RemoveKeyBindByAddress(address)
}
func (o *StrategyAuth) RemoveGroup(name string) error {
	return o.Internal.RemoveGroup(name)
}
func (o *StrategyAuth) RemoveMethodTemplate(name string) error {
	return o.Internal.RemoveMethodTemplate(name)
}
func (o *StrategyAuth) RemoveMsgTypeTemplate(name string) error {
	return o.Internal.RemoveMsgTypeTemplate(name)
}
