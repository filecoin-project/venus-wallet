package sqlite

import (
	"github.com/filecoin-project/venus-wallet/storage"
	"strings"
)

type iStrategyMapper interface {
	toInnerMsgTypeTemplate(outer *storage.MsgTypeTemplate) *MsgTypeTemplate
	toOuterMsgTypeTemplate(inner *MsgTypeTemplate) *storage.MsgTypeTemplate
	toOuterMsgTypeTemplates(inners []*MsgTypeTemplate) []*storage.MsgTypeTemplate

	toInnerMethodTemplate(outer *storage.MethodTemplate) *MethodTemplate
	toInnerMethodTemplates(outers []*storage.MethodTemplate) []*MethodTemplate
	toOuterMethodTemplate(inner *MethodTemplate) *storage.MethodTemplate
	toOuterMethodTemplates(inner []*MethodTemplate) []*storage.MethodTemplate

	toInnerKeyBind(outer *storage.KeyBind) *KeyBind
	toOuterKeyBind(inner *KeyBind) *storage.KeyBind
	toOuterKeyBinds(inner []*KeyBind) []*storage.KeyBind

	toOuterGroup(inner *Group, kbs []*KeyBind) *storage.Group
	toOuterGroups(inner []*Group) []*storage.Group
}

const splitChar = ","

type strategyMapper struct {
}

func newRouterMapper() iStrategyMapper {
	return &strategyMapper{}
}

func (m *strategyMapper) toInnerMsgTypeTemplate(outer *storage.MsgTypeTemplate) *MsgTypeTemplate {
	if outer == nil {
		return nil
	}
	res := &MsgTypeTemplate{
		MTTId:     outer.MTTId,
		Name:      outer.Name,
		MetaTypes: outer.MetaTypes,
	}
	return res
}

func (m *strategyMapper) toOuterMsgTypeTemplate(inner *MsgTypeTemplate) *storage.MsgTypeTemplate {
	if inner == nil {
		return nil
	}
	res := &storage.MsgTypeTemplate{
		MTTId:     inner.MTTId,
		Name:      inner.Name,
		MetaTypes: inner.MetaTypes,
	}
	return res
}

func (m *strategyMapper) toOuterMsgTypeTemplates(inners []*MsgTypeTemplate) []*storage.MsgTypeTemplate {
	res := make([]*storage.MsgTypeTemplate, 0, len(inners))
	for _, v := range inners {
		res = append(res, m.toOuterMsgTypeTemplate(v))
	}
	return res
}

func (m *strategyMapper) toInnerMethodTemplate(outer *storage.MethodTemplate) *MethodTemplate {
	if outer == nil {
		return nil
	}
	return &MethodTemplate{
		MTId: outer.MTId,
		Name: outer.Name,
		// method name join with ','
		Methods: strings.Join(outer.Methods, splitChar),
	}
}

func (m *strategyMapper) toInnerMethodTemplates(outers []*storage.MethodTemplate) []*MethodTemplate {
	res := make([]*MethodTemplate, 0, len(outers))
	for _, v := range outers {
		res = append(res, m.toInnerMethodTemplate(v))
	}
	return res
}

func (m *strategyMapper) toOuterMethodTemplate(inner *MethodTemplate) *storage.MethodTemplate {
	if inner == nil {
		return nil
	}
	return &storage.MethodTemplate{
		MTId:    inner.MTId,
		Name:    inner.Name,
		Methods: strings.Split(inner.Methods, splitChar),
	}
}
func (m *strategyMapper) toOuterMethodTemplates(inner []*MethodTemplate) []*storage.MethodTemplate {
	res := make([]*storage.MethodTemplate, 0, len(inner))
	for _, v := range inner {
		res = append(res, m.toOuterMethodTemplate(v))
	}
	return res
}

func (m *strategyMapper) toInnerKeyBind(outer *storage.KeyBind) *KeyBind {
	if outer == nil {
		return nil
	}
	return &KeyBind{
		BindId:    outer.BindId,
		Name:      outer.Name,
		Address:   outer.Address,
		MetaTypes: outer.MetaTypes,
		// source from MethodTemplate
		MethodNames: strings.Join(outer.Methods, splitChar),
	}
}

func (m *strategyMapper) toOuterKeyBind(inner *KeyBind) *storage.KeyBind {
	if inner == nil {
		return nil
	}
	return &storage.KeyBind{
		BindId:    inner.BindId,
		Name:      inner.Name,
		Address:   inner.Address,
		MetaTypes: inner.MetaTypes,
		// source from MethodTemplate
		Methods: strings.Split(inner.MethodNames, splitChar),
	}
}

func (m *strategyMapper) toOuterKeyBinds(inner []*KeyBind) []*storage.KeyBind {
	res := make([]*storage.KeyBind, 0, len(inner))
	for _, v := range inner {
		res = append(res, m.toOuterKeyBind(v))
	}
	return res
}

func (m *strategyMapper) toOuterGroup(inner *Group, kbs []*KeyBind) *storage.Group {
	if inner == nil {
		return nil
	}
	return &storage.Group{
		GroupId:  inner.GroupId,
		Name:     inner.Name,
		KeyBinds: m.toOuterKeyBinds(kbs),
	}
}

func (m *strategyMapper) toOuterGroups(inner []*Group) []*storage.Group {
	res := make([]*storage.Group, 0, len(inner))
	for _, v := range inner {
		res = append(res, m.toOuterGroup(v, nil))
	}
	return res
}
