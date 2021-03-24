package api

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"github.com/ipfs-force-community/venus-wallet/storage/strategy"
)

var _ strategy.ILocalStrategy = &StrategyAuth{}

type StrategyAuth struct {
	Internal struct {
		NewMsgTypeTemplate     func(ctx context.Context, name string, codes []int) error                                     `perm:"admin" local:"required"`
		NewMethodTemplate      func(ctx context.Context, name string, methods []string) error                                `perm:"admin" local:"required"`
		NewKeyBindCustom       func(ctx context.Context, name, address string, codes []int, methods []core.MethodName) error `perm:"admin" local:"required"`
		NewKeyBindFromTemplate func(ctx context.Context, name, address, mttName, mtName string) error                        `perm:"admin" local:"required"`
		NewGroup               func(ctx context.Context, name string, keyBindNames []string) error                           `perm:"admin" local:"required"`
		NewWalletToken         func(ctx context.Context, groupName string) (token string, err error)                         `perm:"admin" local:"required"`

		GetMsgTypeTemplate      func(ctx context.Context, name string) (*storage.MsgTypeTemplate, error) `perm:"admin" local:"required"`
		GetMethodTemplateByName func(ctx context.Context, name string) (*storage.MethodTemplate, error)  `perm:"admin" local:"required"`
		GetKeyBindByName        func(ctx context.Context, name string) (*storage.KeyBind, error)         `perm:"admin" local:"required"`
		GetKeyBinds             func(ctx context.Context, address string) ([]*storage.KeyBind, error)    `perm:"admin" local:"required"`
		GetGroupByName          func(ctx context.Context, name string) (*storage.Group, error)           `perm:"admin" local:"required"`
		GetWalletTokensByGroup  func(ctx context.Context, groupName string) ([]string, error)            `perm:"admin" local:"required"`
		GetWalletTokenInfo      func(ctx context.Context, token string) (*storage.GroupAuth, error)      `perm:"admin" local:"required"`

		ListGroups           func(ctx context.Context, fromIndex, toIndex int) ([]*storage.Group, error)           `perm:"admin" local:"required"`
		ListKeyBinds         func(ctx context.Context, fromIndex, toIndex int) ([]*storage.KeyBind, error)         `perm:"admin" local:"required"`
		ListMethodTemplates  func(ctx context.Context, fromIndex, toIndex int) ([]*storage.MethodTemplate, error)  `perm:"admin" local:"required"`
		ListMsgTypeTemplates func(ctx context.Context, fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) `perm:"admin" local:"required"`

		PushMsgTypeIntoKeyBind func(ctx context.Context, name string, codes []int) (*storage.KeyBind, error)      `perm:"admin" local:"required"`
		PushMethodIntoKeyBind  func(ctx context.Context, name string, methods []string) (*storage.KeyBind, error) `perm:"admin" local:"required"`
		PullMsgTypeFromKeyBind func(ctx context.Context, name string, codes []int) (*storage.KeyBind, error)      `perm:"admin" local:"required"`
		PullMethodFromKeyBind  func(ctx context.Context, name string, methods []string) (*storage.KeyBind, error) `perm:"admin" local:"required"`

		RemoveToken            func(ctx context.Context, token string) error                                                  `perm:"admin" local:"required"`
		RemoveKeyBind          func(ctx context.Context, name string) error                                                   `perm:"admin" local:"required"`
		RemoveKeyBindByAddress func(ctx context.Context, address string) (int64, error)                                       `perm:"admin" local:"required"`
		RemoveGroup            func(ctx context.Context, name string) error                                                   `perm:"admin" local:"required"`
		RemoveMethodTemplate   func(ctx context.Context, name string) error                                                   `perm:"admin" local:"required"`
		RemoveMsgTypeTemplate  func(ctx context.Context, name string) error                                                   `perm:"admin" local:"required"`
		ScopeWallet            func(ctx context.Context) ([]core.Address, error)                                              `perm:"admin" local:"required"`
		Verify                 func(ctx context.Context, address core.Address, msgType core.MsgType, msg *core.Message) error `perm:"admin" local:"required"`
	}
}

func (o *StrategyAuth) NewWalletToken(ctx context.Context, groupName string) (token string, err error) {
	return o.Internal.NewWalletToken(ctx, groupName)
}

func (o *StrategyAuth) GetWalletTokensByGroup(ctx context.Context, groupName string) ([]string, error) {
	return o.Internal.GetWalletTokensByGroup(ctx, groupName)
}

func (o *StrategyAuth) NewMsgTypeTemplate(ctx context.Context, name string, codes []int) error {
	return o.Internal.NewMsgTypeTemplate(ctx, name, codes)
}
func (o *StrategyAuth) NewMethodTemplate(ctx context.Context, name string, methods []string) error {
	return o.Internal.NewMethodTemplate(ctx, name, methods)
}
func (o *StrategyAuth) NewKeyBindCustom(ctx context.Context, name, address string, codes []int, methods []core.MethodName) error {
	return o.Internal.NewKeyBindCustom(ctx, name, address, codes, methods)
}
func (o *StrategyAuth) NewKeyBindFromTemplate(ctx context.Context, name, address, mttName, mtName string) error {
	return o.Internal.NewKeyBindFromTemplate(ctx, name, address, mttName, mtName)
}
func (o *StrategyAuth) NewGroup(ctx context.Context, name string, keyBindNames []string) error {
	return o.Internal.NewGroup(ctx, name, keyBindNames)
}

func (o *StrategyAuth) GetMsgTypeTemplate(ctx context.Context, name string) (*storage.MsgTypeTemplate, error) {
	return o.Internal.GetMsgTypeTemplate(ctx, name)
}
func (o *StrategyAuth) GetMethodTemplateByName(ctx context.Context, name string) (*storage.MethodTemplate, error) {
	return o.Internal.GetMethodTemplateByName(ctx, name)
}
func (o *StrategyAuth) GetKeyBindByName(ctx context.Context, name string) (*storage.KeyBind, error) {
	return o.Internal.GetKeyBindByName(ctx, name)
}
func (o *StrategyAuth) GetKeyBinds(ctx context.Context, address string) ([]*storage.KeyBind, error) {
	return o.Internal.GetKeyBinds(ctx, address)
}
func (o *StrategyAuth) GetGroupByName(ctx context.Context, name string) (*storage.Group, error) {
	return o.Internal.GetGroupByName(ctx, name)
}

func (o *StrategyAuth) GetWalletTokenInfo(ctx context.Context, token string) (*storage.GroupAuth, error) {
	return o.Internal.GetWalletTokenInfo(ctx, token)
}

func (o *StrategyAuth) ListGroups(ctx context.Context, fromIndex, toIndex int) ([]*storage.Group, error) {
	return o.Internal.ListGroups(ctx, fromIndex, toIndex)
}
func (o *StrategyAuth) ListKeyBinds(ctx context.Context, fromIndex, toIndex int) ([]*storage.KeyBind, error) {
	return o.Internal.ListKeyBinds(ctx, fromIndex, toIndex)
}
func (o *StrategyAuth) ListMethodTemplates(ctx context.Context, fromIndex, toIndex int) ([]*storage.MethodTemplate, error) {
	return o.Internal.ListMethodTemplates(ctx, fromIndex, toIndex)
}
func (o *StrategyAuth) ListMsgTypeTemplates(ctx context.Context, fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) {
	return o.Internal.ListMsgTypeTemplates(ctx, fromIndex, toIndex)
}

func (o *StrategyAuth) PushMsgTypeIntoKeyBind(ctx context.Context, name string, codes []int) (*storage.KeyBind, error) {
	return o.Internal.PushMsgTypeIntoKeyBind(ctx, name, codes)
}
func (o *StrategyAuth) PushMethodIntoKeyBind(ctx context.Context, name string, methods []string) (*storage.KeyBind, error) {
	return o.Internal.PushMethodIntoKeyBind(ctx, name, methods)
}
func (o *StrategyAuth) PullMsgTypeFromKeyBind(ctx context.Context, name string, codes []int) (*storage.KeyBind, error) {
	return o.Internal.PullMsgTypeFromKeyBind(ctx, name, codes)
}
func (o *StrategyAuth) PullMethodFromKeyBind(ctx context.Context, name string, methods []string) (*storage.KeyBind, error) {
	return o.Internal.PullMethodFromKeyBind(ctx, name, methods)
}

func (o *StrategyAuth) RemoveKeyBind(ctx context.Context, name string) error {
	return o.Internal.RemoveKeyBind(ctx, name)
}
func (o *StrategyAuth) RemoveKeyBindByAddress(ctx context.Context, address string) (int64, error) {
	return o.Internal.RemoveKeyBindByAddress(ctx, address)
}
func (o *StrategyAuth) RemoveGroup(ctx context.Context, name string) error {
	return o.Internal.RemoveGroup(ctx, name)
}
func (o *StrategyAuth) RemoveMethodTemplate(ctx context.Context, name string) error {
	return o.Internal.RemoveMethodTemplate(ctx, name)
}
func (o *StrategyAuth) RemoveMsgTypeTemplate(ctx context.Context, name string) error {
	return o.Internal.RemoveMsgTypeTemplate(ctx, name)
}
func (o *StrategyAuth) RemoveToken(ctx context.Context, token string) error {
	return o.Internal.RemoveToken(ctx, token)
}
func (o *StrategyAuth) ScopeWallet(ctx context.Context) ([]core.Address, error) {
	return o.Internal.ScopeWallet(ctx)
}
func (o *StrategyAuth) Verify(ctx context.Context, address core.Address, msgType core.MsgType, msg *core.Message) error {
	return o.Internal.Verify(ctx, address, msgType, msg)
}
