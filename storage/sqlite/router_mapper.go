package sqlite

import "github.com/ipfs-force-community/venus-wallet/storage"

type iRouterMapper interface {
	toInnerMsgTypeTemplate(outer *storage.MsgTypeTemplate) *MsgTypeTemplate
	toOuterMsgTypeTemplate(inner *MsgTypeTemplate) *storage.MsgTypeTemplate
	toOuterMsgTypeTemplates(inners []*MsgTypeTemplate) []*storage.MsgTypeTemplate
}

type routerMapper struct {
}

func newRouterMapper() iRouterMapper {
	return &routerMapper{}
}

func (m *routerMapper) toInnerMsgTypeTemplate(outer *storage.MsgTypeTemplate) *MsgTypeTemplate {
	if outer == nil {
		return nil
	}
	res := &MsgTypeTemplate{
		Name:      outer.Name,
		MetaTypes: outer.MetaTypes,
	}
	return res
}

func (m *routerMapper) toOuterMsgTypeTemplate(inner *MsgTypeTemplate) *storage.MsgTypeTemplate {
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

func (m *routerMapper) toOuterMsgTypeTemplates(inners []*MsgTypeTemplate) []*storage.MsgTypeTemplate {
	res := make([]*storage.MsgTypeTemplate, 0, len(inners))
	for _, v := range inners {
		res = append(res, m.toOuterMsgTypeTemplate(v))
	}
	return res
}
