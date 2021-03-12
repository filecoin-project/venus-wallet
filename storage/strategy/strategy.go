package strategy

import (
	"bytes"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/errcode"
	"github.com/ipfs-force-community/venus-wallet/msgrouter"
	"github.com/ipfs-force-community/venus-wallet/storage"
)

type strategy struct {
	store storage.StrategyStore
}

func (s *strategy) NewMsgTypeTemplate(name string, enum core.MsgEnum) error {
	return s.store.PutMsgTypeTemplate(&storage.MsgTypeTemplate{
		Name:      name,
		MetaTypes: enum,
	})
}

func (s *strategy) GetMsgTypeTemplate(name string) (*storage.MsgTypeTemplate, error) {
	return s.store.GetMsgTypeTemplateByName(name)
}

func (s *strategy) ListMsgTypeTemplates(fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) {
	return s.store.ListMsgTypeTemplates(fromIndex, toIndex)
}

func (s *strategy) DeleteMsgTypeTemplate(name string) error {
	m, err := s.store.GetMsgTypeTemplateByName(name)
	if err != nil {
		return err
	}
	return s.store.DeleteMsgTypeTemplate(m.MTTId)
}

func (s *strategy) NewMethodTemplate(name string, methods []string) error {
	return s.store.PutMethodTemplate(&storage.MethodTemplate{
		Name:    name,
		Methods: methods,
	})
}

func (s *strategy) GetMethodTemplateByName(name string) (*storage.MethodTemplate, error) {
	return s.store.GetMethodTemplateByName(name)
}

func (s *strategy) ListMethodTemplates(fromIndex, toIndex int) ([]*storage.MethodTemplate, error) {
	return s.store.ListMethodTemplates(fromIndex, toIndex)
}

func (s *strategy) RemoveMethodTemplate(name string) error {
	m, err := s.store.GetMethodTemplateByName(name)
	if err != nil {
		return err
	}
	return s.store.DeleteMethodTemplate(m.MTId)
}

func (s *strategy) NewKeyBindCustom(name, address string, codes []int, methods []msgrouter.MethodName) error {
	em, err := core.AggregateMsgEnumCode(codes)
	if err != nil {
		return err
	}
	ms, err := msgrouter.AggregateMethodNames(methods)
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
func (s *strategy) NewKeyBindFromTemplate(name, address, mttName, mtName string) error {
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
func (s *strategy) GetKeyBindByName(name string) (*storage.KeyBind, error) {
	return s.store.GetKeyBindByName(name)
}

func (s *strategy) GetKeyBinds(address string) ([]*storage.KeyBind, error) {
	return s.store.GetKeyBinds(address)
}

func (s *strategy) ListKeyBinds(fromIndex, toIndex int) ([]*storage.KeyBind, error) {
	return s.store.ListKeyBinds(fromIndex, toIndex)
}

func (s *strategy) RemoveKeyBind(name string) error {
	kb, err := s.store.GetKeyBindByName(name)
	if err != nil {
		return err
	}
	return s.store.DeleteKeyBind(kb.BindId)
}
func (s *strategy) RemoveKeyBindByAddress(address string) (int64, error) {
	return s.store.DeleteKeyBindsByAddress(address)
}

func (s *strategy) NewGroup(name string, keyBindNames []string) error {
	if len(keyBindNames) == 0 {
		return errcode.ErrNilReference
	}
	linq.From(keyBindNames).Distinct().ToSlice(keyBindNames)
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

func (s *strategy) GetGroupByName(name string) (*storage.Group, error) {
	return s.store.GetGroupByName(name)
}
func (s *strategy) ListGroups(fromIndex, toIndex int) ([]*storage.Group, error) {
	return s.store.ListGroups(fromIndex, toIndex)
}

func (s *strategy) RemoveGroup(name string) error {
	g, err := s.store.GetGroupByName(name)
	if err != nil {
		return err
	}
	return s.store.DeleteGroup(g.GroupId)
}

func (s *strategy) Verify(token, address string, enum core.MsgEnum, method msgrouter.MethodName) error {
	return nil
}
