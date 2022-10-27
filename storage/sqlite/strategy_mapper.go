package sqlite

import (
	"strings"

	types "github.com/filecoin-project/venus/venus-shared/types/wallet"
)

type iStrategyMapper interface {
	toInnerMsgTypeTemplate(outer *types.MsgTypeTemplate) *MsgTypeTemplate
	toOuterMsgTypeTemplate(inner *MsgTypeTemplate) *types.MsgTypeTemplate
	toOuterMsgTypeTemplates(inners []*MsgTypeTemplate) []*types.MsgTypeTemplate

	toInnerMethodTemplate(outer *types.MethodTemplate) *MethodTemplate
	toInnerMethodTemplates(outers []*types.MethodTemplate) []*MethodTemplate
	toOuterMethodTemplate(inner *MethodTemplate) *types.MethodTemplate
	toOuterMethodTemplates(inner []*MethodTemplate) []*types.MethodTemplate

	toInnerKeyBind(outer *types.KeyBind) *KeyBind
	toOuterKeyBind(inner *KeyBind) *types.KeyBind
	toOuterKeyBinds(inner []*KeyBind) []*types.KeyBind

	toOuterGroup(inner *Group, kbs []*KeyBind) *types.Group
	toOuterGroups(inner []*Group) []*types.Group
}

const splitChar = ","

type strategyMapper struct{}

func newRouterMapper() iStrategyMapper {
	return &strategyMapper{}
}

func (m *strategyMapper) toInnerMsgTypeTemplate(outer *types.MsgTypeTemplate) *MsgTypeTemplate {
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

func (m *strategyMapper) toOuterMsgTypeTemplate(inner *MsgTypeTemplate) *types.MsgTypeTemplate {
	if inner == nil {
		return nil
	}
	res := &types.MsgTypeTemplate{
		MTTId:     inner.MTTId,
		Name:      inner.Name,
		MetaTypes: inner.MetaTypes,
	}
	return res
}

func (m *strategyMapper) toOuterMsgTypeTemplates(inners []*MsgTypeTemplate) []*types.MsgTypeTemplate {
	res := make([]*types.MsgTypeTemplate, 0, len(inners))
	for _, v := range inners {
		res = append(res, m.toOuterMsgTypeTemplate(v))
	}
	return res
}

func (m *strategyMapper) toInnerMethodTemplate(outer *types.MethodTemplate) *MethodTemplate {
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

func (m *strategyMapper) toInnerMethodTemplates(outers []*types.MethodTemplate) []*MethodTemplate {
	res := make([]*MethodTemplate, 0, len(outers))
	for _, v := range outers {
		res = append(res, m.toInnerMethodTemplate(v))
	}
	return res
}

func (m *strategyMapper) toOuterMethodTemplate(inner *MethodTemplate) *types.MethodTemplate {
	if inner == nil {
		return nil
	}
	return &types.MethodTemplate{
		MTId:    inner.MTId,
		Name:    inner.Name,
		Methods: strings.Split(inner.Methods, splitChar),
	}
}

func (m *strategyMapper) toOuterMethodTemplates(inner []*MethodTemplate) []*types.MethodTemplate {
	res := make([]*types.MethodTemplate, 0, len(inner))
	for _, v := range inner {
		res = append(res, m.toOuterMethodTemplate(v))
	}
	return res
}

func (m *strategyMapper) toInnerKeyBind(outer *types.KeyBind) *KeyBind {
	if outer == nil {
		return nil
	}

	return &KeyBind{
		BindId:    outer.BindID,
		Name:      outer.Name,
		Address:   outer.Address,
		MetaTypes: outer.MetaTypes,
		// source from MethodTemplate
		MethodNames: strings.Join(outer.Methods, splitChar),
	}
}

func (m *strategyMapper) toOuterKeyBind(inner *KeyBind) *types.KeyBind {
	if inner == nil {
		return nil
	}
	return &types.KeyBind{
		BindID:    inner.BindId,
		Name:      inner.Name,
		Address:   inner.Address,
		MetaTypes: inner.MetaTypes,
		// source from MethodTemplate
		Methods: strings.Split(inner.MethodNames, splitChar),
	}
}

func (m *strategyMapper) toOuterKeyBinds(inner []*KeyBind) []*types.KeyBind {
	res := make([]*types.KeyBind, 0, len(inner))
	for _, v := range inner {
		res = append(res, m.toOuterKeyBind(v))
	}
	return res
}

func (m *strategyMapper) toOuterGroup(inner *Group, kbs []*KeyBind) *types.Group {
	if inner == nil {
		return nil
	}
	return &types.Group{
		GroupID:  inner.GroupId,
		Name:     inner.Name,
		KeyBinds: m.toOuterKeyBinds(kbs),
	}
}

func (m *strategyMapper) toOuterGroups(inner []*Group) []*types.Group {
	res := make([]*types.Group, 0, len(inner))
	for _, v := range inner {
		res = append(res, m.toOuterGroup(v, nil))
	}
	return res
}
