package sqlite

import "github.com/filecoin-project/venus-wallet/errcode"

func (s *strategyStore) hasName(name, tb string) error {
	var count int64
	err := s.db.Table(tb).Where("name=?", name).Count(&count).Error
	if err != nil {
		log.Error(err)
		return errcode.ErrDBOperationFailed
	}
	if count > 0 {
		return errcode.ErrDataAlreadyExists
	}
	return nil
}

func (s *strategyStore) errorAdapter(err error) error {
	errMsg := err.Error()
	if errMsg == "record not found" {
		return errcode.ErrDataNotExists
	}
	return errcode.ErrDBOperationFailed
}
