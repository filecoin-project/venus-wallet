package strategy

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/filecoin-project/go-address"
	"github.com/google/uuid"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/errcode"
	"github.com/filecoin-project/venus-wallet/node"
	"github.com/filecoin-project/venus-wallet/storage"
	"sync"
)

var (
	ErrGenToken            = errors.New("token generation failed")
	ErrIllegalMetaType     = errors.New("this metaType or method is not allowed to pass")
	ErrGenerateTokenFailed = errors.New("generate token failed")
)

var _ ILocalStrategy = &strategy{}

type IStrategy interface {
	NewMsgTypeTemplate(ctx context.Context, name string, codes []int) error
	NewMethodTemplate(ctx context.Context, name string, methods []string) error
	NewKeyBindCustom(ctx context.Context, name, address string, codes []int, methods []core.MethodName) error
	NewKeyBindFromTemplate(ctx context.Context, name, address, mttName, mtName string) error
	NewGroup(ctx context.Context, name string, keyBindNames []string) error
	NewWalletToken(ctx context.Context, groupName string) (token string, err error)

	GetMsgTypeTemplate(ctx context.Context, name string) (*storage.MsgTypeTemplate, error)
	GetMethodTemplateByName(ctx context.Context, name string) (*storage.MethodTemplate, error)
	GetKeyBindByName(ctx context.Context, name string) (*storage.KeyBind, error)
	GetKeyBinds(ctx context.Context, address string) ([]*storage.KeyBind, error)
	GetGroupByName(ctx context.Context, name string) (*storage.Group, error)
	GetWalletTokensByGroup(ctx context.Context, groupName string) ([]string, error)
	GetWalletTokenInfo(ctx context.Context, token string) (*storage.GroupAuth, error)

	ListGroups(ctx context.Context, fromIndex, toIndex int) ([]*storage.Group, error)
	ListKeyBinds(ctx context.Context, fromIndex, toIndex int) ([]*storage.KeyBind, error)
	ListMethodTemplates(ctx context.Context, fromIndex, toIndex int) ([]*storage.MethodTemplate, error)
	ListMsgTypeTemplates(ctx context.Context, fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error)

	PushMsgTypeIntoKeyBind(ctx context.Context, name string, codes []int) (*storage.KeyBind, error)
	PushMethodIntoKeyBind(ctx context.Context, name string, methods []string) (*storage.KeyBind, error)
	PullMsgTypeFromKeyBind(ctx context.Context, name string, codes []int) (*storage.KeyBind, error)
	PullMethodFromKeyBind(ctx context.Context, name string, methods []string) (*storage.KeyBind, error)

	RemoveMsgTypeTemplate(nctx context.Context, ame string) error
	RemoveGroup(ctx context.Context, name string) error
	RemoveMethodTemplate(ctx context.Context, name string) error
	RemoveKeyBind(ctx context.Context, name string) error
	RemoveKeyBindByAddress(ctx context.Context, address string) (int64, error)
	RemoveToken(ctx context.Context, token string) error
}
type ILocalStrategy interface {
	IStrategyVerify
	IStrategy
}
type VerifyFunc func(token, address string, enum core.MsgEnum, method core.MethodName) error

// NOTE: for wallet
type IStrategyVerify interface {
	Verify(ctx context.Context, address core.Address, msgType core.MsgType, msg *core.Message) error
	//@bool: root can do anything
	ScopeWallet(ctx context.Context) (*core.AddressScope, error)
	ContainWallet(ctx context.Context, address core.Address) bool
}

// TODO: add Cache
type strategy struct {
	scache  StrategyCache
	store   storage.StrategyStore
	nodeCli *node.NodeClient
	mw      storage.KeyMiddleware
	sync.RWMutex
}

func NewStrategy(store storage.StrategyStore, nodeCli *node.NodeClient, mw storage.KeyMiddleware) ILocalStrategy {
	return &strategy{
		store:   store,
		nodeCli: nodeCli,
		scache:  newStrategyCache(),
		mw:      mw,
	}
}
func (s *strategy) NewMsgTypeTemplate(ctx context.Context, name string, codes []int) error {
	enum, err := core.AggregateMsgEnumCode(codes)
	if err != nil {
		return err
	}
	return s.store.PutMsgTypeTemplate(&storage.MsgTypeTemplate{
		Name:      name,
		MetaTypes: enum,
	})
}

func (s *strategy) GetMsgTypeTemplate(ctx context.Context, name string) (*storage.MsgTypeTemplate, error) {
	return s.store.GetMsgTypeTemplateByName(name)
}

func (s *strategy) ListMsgTypeTemplates(ctx context.Context, fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) {
	return s.store.ListMsgTypeTemplates(fromIndex, toIndex)
}

func (s *strategy) RemoveMsgTypeTemplate(ctx context.Context, name string) error {
	m, err := s.store.GetMsgTypeTemplateByName(name)
	if err != nil {
		return err
	}
	return s.store.DeleteMsgTypeTemplate(m.MTTId)
}

func (s *strategy) NewMethodTemplate(ctx context.Context, name string, methods []string) error {
	return s.store.PutMethodTemplate(&storage.MethodTemplate{
		Name:    name,
		Methods: methods,
	})
}

func (s *strategy) GetMethodTemplateByName(ctx context.Context, name string) (*storage.MethodTemplate, error) {
	return s.store.GetMethodTemplateByName(name)
}

func (s *strategy) ListMethodTemplates(ctx context.Context, fromIndex, toIndex int) ([]*storage.MethodTemplate, error) {
	return s.store.ListMethodTemplates(fromIndex, toIndex)
}

func (s *strategy) RemoveMethodTemplate(ctx context.Context, name string) error {
	m, err := s.store.GetMethodTemplateByName(name)
	if err != nil {
		return err
	}
	return s.store.DeleteMethodTemplate(m.MTId)
}

func (s *strategy) NewKeyBindCustom(ctx context.Context, name, address string, codes []int, methods []core.MethodName) error {
	em, err := core.AggregateMsgEnumCode(codes)
	if err != nil {
		return err
	}
	ms, err := core.AggregateMethodNames(methods)
	if err != nil {
		return err
	}
	kb := &storage.KeyBind{
		Name:      name,
		Address:   address,
		MetaTypes: em,
		Methods:   ms,
	}
	err = s.store.PutKeyBind(kb)
	if err != nil {
		return fmt.Errorf("failed to store keybind:%s ", err)
	}
	return nil
}

func (s *strategy) NewKeyBindFromTemplate(ctx context.Context, name, address, mttName, mtName string) error {
	mtt, err := s.store.GetMsgTypeTemplateByName(mttName)
	if err != nil {
		return fmt.Errorf("find msgType template failed:%s", err)
	}
	mt, err := s.store.GetMethodTemplateByName(mtName)
	if err != nil {
		return fmt.Errorf("find msg method template failed:%s", err)
	}
	kb := &storage.KeyBind{
		Name:      name,
		Address:   address,
		MetaTypes: mtt.MetaTypes,
		Methods:   mt.Methods,
	}
	err = s.store.PutKeyBind(kb)
	if err != nil {
		return fmt.Errorf("failed to store keybind:%s ", err)
	}
	return nil
}
func (s *strategy) GetKeyBindByName(ctx context.Context, name string) (*storage.KeyBind, error) {
	return s.store.GetKeyBindByName(name)
}

func (s *strategy) GetKeyBinds(ctx context.Context, address string) ([]*storage.KeyBind, error) {
	return s.store.GetKeyBinds(address)
}

func (s *strategy) ListKeyBinds(ctx context.Context, fromIndex, toIndex int) ([]*storage.KeyBind, error) {
	return s.store.ListKeyBinds(fromIndex, toIndex)
}

func (s *strategy) RemoveKeyBind(ctx context.Context, name string) error {
	s.Lock()
	defer s.Unlock()
	err := s.mw.CheckToken(ctx)
	if err != nil {
		return err
	}
	kb, err := s.store.GetKeyBindByName(name)
	if err != nil {
		return err
	}
	err = s.store.DeleteKeyBind(kb.BindId)
	if err != nil {
		s.scache.removeKeyBind(kb)
	}
	return nil
}

func (s *strategy) RemoveKeyBindByAddress(ctx context.Context, address string) (int64, error) {
	s.Lock()
	defer s.Unlock()
	err := s.mw.CheckToken(ctx)
	if err != nil {
		return 0, err
	}
	num, err := s.store.DeleteKeyBindsByAddress(address)
	if err != nil {
		s.scache.removeAddress(address)
	}
	return num, nil
}

func (s *strategy) PushMsgTypeIntoKeyBind(ctx context.Context, name string, codes []int) (*storage.KeyBind, error) {
	s.Lock()
	defer s.Unlock()
	err := s.mw.CheckToken(ctx)
	if err != nil {
		return nil, err
	}
	em, err := core.AggregateMsgEnumCode(codes)
	if err != nil {
		return nil, err
	}
	kb, err := s.store.GetKeyBindByName(name)
	if err != nil {
		return nil, err
	}
	inc := em - kb.MetaTypes&em
	if inc == 0 {
		return kb, nil
	}
	kb.MetaTypes += inc
	err = s.store.UpdateKeyBindMetaTypes(kb)
	if err != nil {
		return nil, err
	}
	s.scache.removeKeyBind(kb)
	return kb, nil
}

func (s *strategy) PushMethodIntoKeyBind(ctx context.Context, name string, methods []string) (*storage.KeyBind, error) {
	s.Lock()
	defer s.Unlock()
	err := s.mw.CheckToken(ctx)
	if err != nil {
		return nil, err
	}
	em, err := core.AggregateMethodNames(methods)
	if err != nil {
		return nil, err
	}
	kb, err := s.store.GetKeyBindByName(name)
	if err != nil {
		return nil, err
	}
	linq.From(em).Except(linq.From(kb.Methods)).ToSlice(&em)
	if len(em) == 0 {
		return kb, nil
	}
	kb.Methods = append(kb.Methods, em...)
	err = s.store.UpdateKeyBindMetaTypes(kb)
	if err != nil {
		return nil, err
	}
	s.scache.removeKeyBind(kb)
	return kb, nil
}

func (s *strategy) PullMsgTypeFromKeyBind(ctx context.Context, name string, codes []int) (*storage.KeyBind, error) {
	s.Lock()
	defer s.Unlock()
	err := s.mw.CheckToken(ctx)
	if err != nil {
		return nil, err
	}
	em, err := core.AggregateMsgEnumCode(codes)
	if err != nil {
		return nil, err
	}
	kb, err := s.store.GetKeyBindByName(name)
	if err != nil {
		return nil, err
	}
	dec := kb.MetaTypes & em
	if dec == 0 {
		return kb, nil
	}
	kb.MetaTypes -= dec
	err = s.store.UpdateKeyBindMetaTypes(kb)
	if err != nil {
		return nil, err
	}
	s.scache.removeKeyBind(kb)
	return kb, nil
}

func (s *strategy) PullMethodFromKeyBind(ctx context.Context, name string, methods []string) (*storage.KeyBind, error) {
	s.Lock()
	defer s.Unlock()
	err := s.mw.CheckToken(ctx)
	if err != nil {
		return nil, err
	}
	em, err := core.AggregateMethodNames(methods)
	if err != nil {
		return nil, err
	}
	kb, err := s.store.GetKeyBindByName(name)
	if err != nil {
		return nil, err
	}
	linq.From(em).Intersect(linq.From(kb.Methods)).ToSlice(&em)
	if len(em) == 0 {
		return kb, nil
	}
	linq.From(kb.Methods).Except(linq.From(em)).ToSlice(&(kb.Methods))
	err = s.store.UpdateKeyBindMetaTypes(kb)
	if err != nil {
		return nil, err
	}
	s.scache.removeKeyBind(kb)
	return kb, nil
}

func (s *strategy) NewGroup(ctx context.Context, name string, keyBindNames []string) error {
	if len(keyBindNames) == 0 {
		return errcode.ErrNilReference
	}
	linq.From(keyBindNames).Distinct().ToSlice(&keyBindNames)
	kbs, err := s.store.GetKeyBindByNames(keyBindNames)
	if err != nil {
		return err
	}
	if len(kbs) < len(keyBindNames) {
		var names []string
		linq.From(keyBindNames).Except(
			linq.From(kbs).SelectT(func(o *storage.KeyBind) string {
				return o.Name
			})).ToSlice(&names)
		buf := new(bytes.Buffer)
		for _, v := range names {
			buf.WriteString(v)
			buf.WriteString(" ")
		}
		return fmt.Errorf("keybind name not exists:%s", buf.String())
	}
	var kbIds []uint
	linq.From(kbs).SelectT(func(o *storage.KeyBind) uint {
		return o.BindId
	}).ToSlice(&kbIds)

	return s.store.PutGroup(name, kbIds)
}

func (s *strategy) GetGroupByName(ctx context.Context, name string) (*storage.Group, error) {
	return s.store.GetGroupByName(name)
}

func (s *strategy) ListGroups(ctx context.Context, fromIndex, toIndex int) ([]*storage.Group, error) {
	return s.store.ListGroups(fromIndex, toIndex)
}

func (s *strategy) RemoveGroup(ctx context.Context, name string) error {
	s.Lock()
	defer s.Unlock()
	err := s.mw.CheckToken(ctx)
	if err != nil {
		return err
	}
	g, err := s.store.GetGroupByName(name)
	if err != nil {
		return err
	}
	err = s.store.DeleteGroup(g.GroupId)
	if err != nil {
		tokens, err := s.store.GetTokensByGroupId(g.GroupId)
		if err == nil {
			s.scache.removeTokens(tokens)
		}
		if err != nil {
			if err == errcode.ErrDataNotExists {
				return nil
			} else {
				s.scache.refresh()
			}
		}
	}
	return nil
}

func (s *strategy) RemoveToken(ctx context.Context, token string) error {
	s.Lock()
	defer s.Unlock()
	err := s.mw.CheckToken(ctx)
	if err != nil {
		return err
	}
	err = s.store.DeleteGroupAuth(token)
	if err != nil {
		s.scache.removeToken(token)
	}
	return err
}

func (s *strategy) NewWalletToken(ctx context.Context, groupName string) (token string, err error) {
	err = s.mw.CheckToken(ctx)
	if err != nil {
		return core.StringEmpty, err
	}
	g, err := s.store.GetGroupByName(groupName)
	if err != nil {
		return core.StringEmpty, err
	}
	tk, err := uuid.NewRandom()
	if err != nil {
		return core.StringEmpty, ErrGenToken
	}
	token = tk.String()
	err = s.store.PutGroupAuth(token, g.GroupId)
	if err != nil {
		return core.StringEmpty, err
	}
	for _, v := range g.KeyBinds {
		s.scache.removeBlank(token, v.Address)
	}
	return token, nil
}

func (s *strategy) GetWalletTokensByGroup(ctx context.Context, groupName string) ([]string, error) {
	g, err := s.store.GetGroupByName(groupName)
	if err != nil {
		return nil, err
	}
	tokens, err := s.store.GetTokensByGroupId(g.GroupId)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *strategy) GetWalletTokenInfo(ctx context.Context, token string) (*storage.GroupAuth, error) {
	return s.store.GetGroupAuth(token)
}

// NOTE: for wallet
func (s *strategy) Verify(ctx context.Context, address core.Address, msgType core.MsgType, msg *core.Message) error {
	s.RLock()
	defer s.RUnlock()
	if core.WalletStrategyLevel == core.SLDisable {
		return nil
	}
	token := core.ContextStrategyToken(ctx)
	var (
		err error
		kb  *storage.KeyBind
	)
	addrStr := address.String()
	// cache
	kb = s.scache.get(token, addrStr)
	if kb != nil {
		goto Verify
	}
	// is data through
	if s.scache.through(token, addrStr) {
		return errcode.ErrDataNotExists
	}
	kb, err = s.store.GetGroupKeyBind(token, address.String())
	if err == nil {
		s.scache.set(token, kb)
		goto Verify
	}
	if err == errcode.ErrDataNotExists {
		s.scache.setBlank(token, addrStr)
	}
	return err
Verify:
	if !kb.ContainMsgType(msgType) {
		return fmt.Errorf("%s: msgType %s", ErrIllegalMetaType, msgType)
	}
	if core.WalletStrategyLevel == core.SLMethod && msgType == core.MTChainMsg {
		if msg.Method == 0 {
			if !linq.From(kb.Methods).Contains("Send") {
				return fmt.Errorf("%s: method %s", ErrIllegalMetaType, "Send")
			}
		} else {
			actor, err := s.nodeCli.StateGetActor(ctx, msg.To, node.TipSetKey{})
			if err != nil {
				return err
			}
			fn, err := core.GetMethodName(actor.Code, msg.Method)
			if err != nil {
				return err
			}
			if !kb.ContainMethod(fn) {
				return fmt.Errorf("%s: method %s", ErrIllegalMetaType, fn)
			}
		}
	}
	return nil
}

// NOTE: for wallet
// list wallets binding with group token
// level: 0  all pass
// level: 1  check token
// level: 2  check token
func (s *strategy) ScopeWallet(ctx context.Context) (*core.AddressScope, error) {
	// strategy disable, root view
	if core.WalletStrategyLevel == core.SLDisable {
		return &core.AddressScope{Root: true}, nil
	}
	token := core.ContextStrategyToken(ctx)
	err := s.mw.EqualRootToken(token)
	if err == nil {
		return &core.AddressScope{Root: true}, nil
	}
	//TODO: Rich Domain Mode, need replace
	kb, err := s.store.GetGroupAuth(token)
	if err != nil {
		return &core.AddressScope{Root: false}, err
	}
	var addresses []core.Address
	linq.From(kb.KeyBinds).SelectT(func(i *storage.KeyBind) core.Address {
		addr, _ := address.NewFromString(i.Address)
		return addr
	}).ToSlice(&addresses)
	return &core.AddressScope{Root: false, Addresses: addresses}, nil
}

// level: 0  all pass
// level: 1  check token
// level: 2  check token
func (s *strategy) ContainWallet(ctx context.Context, address core.Address) bool {
	// strategy disable, root view
	if core.WalletStrategyLevel == core.SLDisable {
		return true
	}
	token := core.ContextStrategyToken(ctx)
	err := s.mw.EqualRootToken(token)
	if err == nil {
		return true
	}
	_, err = s.store.GetGroupKeyBind(token, address.String())
	return err == nil
}
