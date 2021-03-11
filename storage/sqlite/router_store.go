package sqlite

import (
	"github.com/ipfs-force-community/venus-wallet/errcode"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"gorm.io/gorm"
)

type routerStore struct {
	db     *gorm.DB
	mapper iRouterMapper
}

var _ storage.RouterStore = &routerStore{}

func NewRouterStore(conn *Conn) storage.RouterStore {
	return &routerStore{
		db:     conn.DB,
		mapper: newRouterMapper(),
	}
}

func (s *routerStore) PutMsgTypeTemplate(mtt *storage.MsgTypeTemplate) error {
	m := s.mapper.toInnerMsgTypeTemplate(mtt)
	err := s.db.Table(TBMsgTypeTemplate).Create(m).Error
	if err != nil {
		log.Error(err)
		return errcode.ErrDBOperationFailed
	}
	return nil
}
func (s *routerStore) GetMsgTypeTemplateByName(name string) ([]*storage.MsgTypeTemplate, error) {
	var arr []*MsgTypeTemplate
	err := s.db.Table(TBMsgTypeTemplate).Where("name=?", name).Find(&arr).Error
	if err != nil {
		log.Error(err)
		return nil, errcode.ErrDBOperationFailed
	}
	res := s.mapper.toOuterMsgTypeTemplates(arr)
	return res, nil
}
func (s *routerStore) GetMsgTypeTemplate(mttId uint) (*storage.MsgTypeTemplate, error) {
	m := new(MsgTypeTemplate)
	err := s.db.Table(TBMsgTypeTemplate).First(m, mttId).Error
	if err != nil {
		log.Error(err)
		return nil, errcode.ErrDBOperationFailed
	}
	res := s.mapper.toOuterMsgTypeTemplate(m)
	return res, nil
}
func (s *routerStore) ListMsgTypeTemplates(fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) {
	var arr []*MsgTypeTemplate
	err := s.db.Table(TBMsgTypeTemplate).Order("id").Offset(fromIndex).Limit(toIndex).Scan(&arr).Error
	if err != nil {
		log.Error(err)
		return nil, errcode.ErrDBOperationFailed
	}
	res := s.mapper.toOuterMsgTypeTemplates(arr)
	return res, nil
}
func (s *routerStore) DeleteMsgTypeTemplate(mttId uint) error {
	err := s.db.Table(TBMsgTypeTemplate).Delete(&MsgTypeTemplate{}, "id=?", mttId).Error
	if err != nil {
		log.Error(err)
		return errcode.ErrDBOperationFailed
	}
	return nil
}

func (s *routerStore) PutMethodTemplate(mt *storage.MethodTemplate) error {
	return nil
}
func (s *routerStore) GetMethodTemplate(mtId uint) (*storage.MethodTemplate, error) {
	return nil, nil
}
func (s *routerStore) ListMethodTemplates(fromIndex, toIndex int) ([]*storage.MethodTemplate, error) {
	return nil, nil
}
func (s *routerStore) DeleteMethodTemplate(mtId uint) error {
	return nil
}

func (s *routerStore) PutKeyBind(kb *storage.KeyBind) error {
	return nil
}
func (s *routerStore) GetKeyBind(address string) ([]*storage.KeyBind, error) {
	return nil, nil
}
func (s *routerStore) GetKeyBindById(kbId uint) (*storage.KeyBind, error) {
	return nil, nil
}
func (s *routerStore) ListKeyBinds(fromIndex, toIndex int) ([]*storage.KeyBind, error) {
	return nil, nil
}

func (s *routerStore) PutGroup(g *storage.Group) error {
	return nil
}
func (s *routerStore) GetGroup(gId uint) (*storage.Group, error) {
	return nil, nil
}

func (s *routerStore) ListGroups(fromIndex, toIndex int) ([]*storage.Group, error) {
	return nil, nil
}
func (s *routerStore) DeleteGroup(gId uint) error {
	return nil
}

func (s *routerStore) PutGroupAuth(ga *storage.GroupAuth) error {
	return nil
}
func (s *routerStore) GetGroupAuth(token string) (*storage.GroupAuth, error) {
	return nil, nil
}
func (s *routerStore) DeleteGroupAuth(token string) error {
	return nil
}

func (s *routerStore) GetKeyStrategy(token, address string) (*storage.GroupAuth, error) {
	return nil, nil
}
