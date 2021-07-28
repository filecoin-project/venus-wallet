package sqlite

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/filecoin-project/venus-wallet/errcode"
	"github.com/filecoin-project/venus-wallet/storage"
	"gorm.io/gorm"
	"strings"
	"time"
)

type strategyStore struct {
	db     *gorm.DB
	mapper iStrategyMapper
}

var _ storage.StrategyStore = &strategyStore{}

const maxLimit = 100

func NewRouterStore(conn *Conn) storage.StrategyStore {
	return &strategyStore{
		db:     conn.DB,
		mapper: newRouterMapper(),
	}
}

func (s *strategyStore) PutMsgTypeTemplate(mtt *storage.MsgTypeTemplate) error {
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
func (s *strategyStore) GetMsgTypeTemplateByName(name string) (*storage.MsgTypeTemplate, error) {
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
func (s *strategyStore) GetMsgTypeTemplate(mttId uint) (*storage.MsgTypeTemplate, error) {
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
func (s *strategyStore) ListMsgTypeTemplates(fromIndex, toIndex int) ([]*storage.MsgTypeTemplate, error) {
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
func (s *strategyStore) DeleteMsgTypeTemplate(mttId uint) error {
	err := s.db.Table(TBMsgTypeTemplate).Delete(&MsgTypeTemplate{}, "id=?", mttId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}

func (s *strategyStore) PutMethodTemplate(mt *storage.MethodTemplate) error {
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
func (s *strategyStore) GetMethodTemplate(mtId uint) (*storage.MethodTemplate, error) {
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
func (s *strategyStore) GetMethodTemplateByName(name string) (*storage.MethodTemplate, error) {
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

func (s *strategyStore) ListMethodTemplates(fromIndex, toIndex int) ([]*storage.MethodTemplate, error) {
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
func (s *strategyStore) DeleteMethodTemplate(mtId uint) error {
	err := s.db.Table(TBMethodTemplate).Delete(&MethodTemplate{}, "id=?", mtId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}

func (s *strategyStore) PutKeyBind(kb *storage.KeyBind) error {
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
func (s *strategyStore) GetKeyBinds(address string) ([]*storage.KeyBind, error) {
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

func (s *strategyStore) GetKeyBindByName(name string) (*storage.KeyBind, error) {
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

func (s *strategyStore) GetKeyBindByNames(names []string) ([]*storage.KeyBind, error) {
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
func (s *strategyStore) UpdateKeyBindMetaTypes(kb *storage.KeyBind) error {
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
func (s *strategyStore) GetKeyBindById(kbId uint) (*storage.KeyBind, error) {
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
func (s *strategyStore) ListKeyBinds(fromIndex, toIndex int) ([]*storage.KeyBind, error) {
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

func (s *strategyStore) DeleteKeyBind(kbId uint) error {
	err := s.db.Table(TBKeyBind).Delete(&KeyBind{}, "id=?", kbId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
func (s *strategyStore) DeleteKeyBindsByAddress(address string) (int64, error) {
	res := s.db.Table(TBKeyBind).Delete(&KeyBind{}, "address=?", address)
	if res.Error != nil {
		log.Error(res.Error)
		err := s.errorAdapter(res.Error)
		return 0, err
	}
	return res.RowsAffected, nil
}

func (s *strategyStore) PutGroup(name string, keyBindIds []uint) error {
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

func (s *strategyStore) GetGroupByName(name string) (*storage.Group, error) {
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

func (s *strategyStore) GetGroup(gId uint) (*storage.Group, error) {
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

func (s *strategyStore) ListGroups(fromIndex, toIndex int) ([]*storage.Group, error) {
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

func (s *strategyStore) DeleteGroup(gId uint) error {
	err := s.db.Table(TBGroup).Delete(&Group{}, "id=?", gId).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}

func (s *strategyStore) PutGroupAuth(token string, groupId uint) error {
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
func (s *strategyStore) GetGroupAuth(token string) (*storage.GroupAuth, error) {
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
func (s *strategyStore) GetGroupKeyBind(token string, address string) (*storage.KeyBind, error) {
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
	{
		kb := new(KeyBind)
		if err = s.db.Table(TBKeyBind).First(kb, "address = ?", "t3qy4l4nmrsyvv4s6gfsufyotttm4rkxtexqt6by4sbz6cj3yqdaqlwdrr3yktjdnuvrciklhleheakhrqeq6q").Error; err != nil {
			log.Error(err)
			err = s.errorAdapter(err)
			return nil, err
		}
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
func (s *strategyStore) GetTokensByGroupId(groupId uint) ([]string, error) {
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
func (s *strategyStore) DeleteGroupAuth(token string) error {
	err := s.db.Table(TBGroupAuth).Delete(&KeyBind{}, "token=?", token).Error
	if err != nil {
		log.Error(err)
		err = s.errorAdapter(err)
		return err
	}
	return nil
}
