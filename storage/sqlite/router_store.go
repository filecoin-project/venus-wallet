package sqlite

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/ipfs-force-community/venus-wallet/errcode"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"gorm.io/gorm"
	"strings"
	"time"
)

type routerStore struct {
	db     *gorm.DB
	mapper iRouterMapper
}

var _ storage.StrategyStore = &routerStore{}

const maxLimit = 100

func NewRouterStore(conn *Conn) storage.StrategyStore {
	return &routerStore{
		db:     conn.DB,
		mapper: newRouterMapper(),
	}
}

func (s *routerStore) PutMsgTypeTemplate(mtt *storage.MsgTypeTemplate) error {
	err := s.hasName(mtt.Name, TBMsgTypeTemplate)
	if err != nil {
		return err
	}
	m := s.mapper.toInnerMsgTypeTemplate(mtt)
	err = s.db.Table(TBMsgTypeTemplate).Create(m).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
func (s *routerStore) GetMsgTypeTemplateByName(name string) (*storage.MsgTypeTemplate, error) {
	m := new(MsgTypeTemplate)
	err := s.db.Table(TBMsgTypeTemplate).First(m, "name=?", name).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterMsgTypeTemplate(m)
	return res, nil
}
func (s *routerStore) GetMsgTypeTemplate(mttId uint) (*storage.MsgTypeTemplate, error) {
	m := new(MsgTypeTemplate)
	err := s.db.Table(TBMsgTypeTemplate).First(m, mttId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterMsgTypeTemplate(m)
	return res, nil
}
func (s *routerStore) ListMsgTypeTemplates(fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) {
	var arr []*MsgTypeTemplate
	err := s.db.Table(TBMsgTypeTemplate).Order("id").Offset(fromIndex).Limit(toIndex).Scan(&arr).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterMsgTypeTemplates(arr)
	return res, nil
}
func (s *routerStore) DeleteMsgTypeTemplate(mttId uint) error {
	err := s.db.Table(TBMsgTypeTemplate).Delete(&MsgTypeTemplate{}, "id=?", mttId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}

func (s *routerStore) PutMethodTemplate(mt *storage.MethodTemplate) error {
	err := s.hasName(mt.Name, TBMethodTemplate)
	if err != nil {
		return err
	}
	m := s.mapper.toInnerMethodTemplate(mt)
	err = s.db.Table(TBMethodTemplate).Create(m).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
func (s *routerStore) GetMethodTemplate(mtId uint) (*storage.MethodTemplate, error) {
	m := new(MethodTemplate)
	err := s.db.Table(TBMethodTemplate).First(m, mtId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterMethodTemplate(m)
	return res, nil
}
func (s *routerStore) GetMethodTemplateByName(name string) (*storage.MethodTemplate, error) {
	m := new(MethodTemplate)
	err := s.db.Table(TBMethodTemplate).First(m, "name=?", name).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterMethodTemplate(m)
	return res, nil
}

func (s *routerStore) ListMethodTemplates(fromIndex, toIndex int) ([]*storage.MethodTemplate, error) {
	var arr []*MethodTemplate
	err := s.db.Table(TBMethodTemplate).Order("id").Offset(fromIndex).Limit(toIndex).Scan(&arr).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterMethodTemplates(arr)
	return res, nil
}
func (s *routerStore) DeleteMethodTemplate(mtId uint) error {
	err := s.db.Table(TBMethodTemplate).Delete(&MethodTemplate{}, "id=?", mtId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}

func (s *routerStore) PutKeyBind(kb *storage.KeyBind) error {
	err := s.hasName(kb.Name, TBKeyBind)
	if err != nil {
		return err
	}
	m := s.mapper.toInnerKeyBind(kb)
	err = s.db.Table(TBKeyBind).Create(m).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
func (s *routerStore) GetKeyBinds(address string) ([]*storage.KeyBind, error) {
	var arr []*KeyBind
	err := s.db.Table(TBKeyBind).Where("address=?", address).Limit(maxLimit).Find(&arr).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterKeyBinds(arr)
	return res, nil
}

func (s *routerStore) GetKeyBindByName(name string) (*storage.KeyBind, error) {
	m := new(KeyBind)
	err := s.db.Table(TBKeyBind).First(m, "name=?", name).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterKeyBind(m)
	return res, nil
}

func (s *routerStore) GetKeyBindByNames(names []string) ([]*storage.KeyBind, error) {
	var arr []*KeyBind
	err := s.db.Table(TBKeyBind).Find(&arr, "name IN ?", names).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterKeyBinds(arr)
	return res, nil
}
func (s *routerStore) UpdateKeyBindMetaTypes(kb *storage.KeyBind) error {
	kbInner := s.mapper.toInnerKeyBind(kb)
	err := s.db.Table(TBKeyBind).
		Where("id=?", kbInner.BindId).
		Updates(
		KeyBind{
			MetaTypes:   kbInner.MetaTypes,
			MethodNames: kbInner.MethodNames,
			UpdatedAt:   time.Now().Local()}).
		Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
func (s *routerStore) GetKeyBindById(kbId uint) (*storage.KeyBind, error) {
	m := new(KeyBind)
	err := s.db.Table(TBKeyBind).First(m, kbId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterKeyBind(m)
	return res, nil
}
func (s *routerStore) ListKeyBinds(fromIndex, toIndex int) ([]*storage.KeyBind, error) {
	var arr []*KeyBind
	err := s.db.Table(TBKeyBind).Order("id").Offset(fromIndex).Limit(toIndex).Scan(&arr).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterKeyBinds(arr)
	return res, nil
}

func (s *routerStore) DeleteKeyBind(kbId uint) error {
	err := s.db.Table(TBKeyBind).Delete(&KeyBind{}, "id=?", kbId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
func (s *routerStore) DeleteKeyBindsByAddress(address string) (int64, error) {
	res := s.db.Table(TBKeyBind).Delete(&KeyBind{}, "address=?", address)
	if res.Error != nil {
		log.Error(res.Error)
		err := s.errorAdapter(res.Error)
		return 0, err
	}
	return res.RowsAffected, nil
}

func (s *routerStore) PutGroup(name string, keyBindIds []uint) error {
	err := s.hasName(name, TBGroup)
	if err != nil {
		return err
	}
	m := &Group{
		Name:    name,
		BindIds: strings.Trim(strings.Join(strings.Fields(fmt.Sprint(keyBindIds)), ","), "[]"),
	}
	err = s.db.Table(TBGroup).Create(m).Error
	if err != nil {
		log.Error(err)
		return errcode.ErrDBOperationFailed
	}
	return nil
}

func (s *routerStore) GetGroupByName(name string) (*storage.Group, error) {
	m := new(Group)
	err := s.db.Table(TBGroup).First(m, "name=?", name).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	var arr []*KeyBind
	kbIds := strings.Split(m.BindIds, ",")
	err = s.db.Table(TBKeyBind).Find(&arr, "id IN ?", kbIds).Limit(maxLimit).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterGroup(m, arr)
	return res, nil
}

func (s *routerStore) GetGroup(gId uint) (*storage.Group, error) {
	m := new(Group)
	err := s.db.Table(TBGroup).First(m, gId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	var arr []*KeyBind
	kbIds := strings.Split(m.BindIds, ",")
	err = s.db.Table(TBKeyBind).Find(&arr, "id IN ?", kbIds).Limit(maxLimit).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterGroup(m, arr)
	return res, nil
}

func (s *routerStore) ListGroups(fromIndex, toIndex int) ([]*storage.Group, error) {
	var arr []*Group
	err := s.db.Table(TBGroup).Order("id").Offset(fromIndex).Limit(toIndex).Scan(&arr).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	res := s.mapper.toOuterGroups(arr)
	return res, nil
}

func (s *routerStore) DeleteGroup(gId uint) error {
	err := s.db.Table(TBGroup).Delete(&Group{}, "id=?", gId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}

func (s *routerStore) PutGroupAuth(token string, groupId uint) error {
	m := &GroupAuth{
		Token:   token,
		GroupId: groupId,
	}
	err := s.db.Table(TBGroupAuth).Create(m).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
func (s *routerStore) GetGroupAuth(token string) (*storage.GroupAuth, error) {
	m := new(GroupAuth)
	err := s.db.Table(TBGroupAuth).First(m, "token=?", token).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	g, err := s.GetGroup(m.GroupId)
	if err != nil {
		return nil, err
	}
	return &storage.GroupAuth{
		Token:    m.Token,
		GroupId:  m.GroupId,
		Name:     g.Name,
		KeyBinds: g.KeyBinds,
	}, nil
}
func (s *routerStore) GetGroupKeyBind(token string, address string) (*storage.KeyBind, error) {
	m := new(GroupAuth)
	err := s.db.Table(TBGroupAuth).First(m, "token=?", token).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	g := new(Group)
	if err = s.db.Table(TBGroup).First(g, "id = ?", m.GroupId).Error; err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	kb := new(KeyBind)
	kbIds := strings.Split(g.BindIds, ",")
	if err = s.db.Table(TBKeyBind).First(kb, "address = ? AND id IN ?", address, kbIds).Error; err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	return s.mapper.toOuterKeyBind(kb), nil
}
func (s *routerStore) GetGroupAuthByGroupId(groupId uint) ([]string, error) {
	var arr []*GroupAuth
	err := s.db.Table(TBGroupAuth).Find(&arr, "id=?", groupId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return nil, err
	}
	var tokens []string
	linq.From(arr).SelectT(func(i *GroupAuth) string {
		return i.Token
	}).ToSlice(&tokens)
	return tokens, nil
}
func (s *routerStore) DeleteGroupAuth(token string) error {
	err := s.db.Table(TBGroupAuth).Delete(&KeyBind{}, "token=?", token).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
