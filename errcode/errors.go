package errcode

import "errors"

var (
	ErrDBOperationFailed = errors.New("database operation exception")
	ErrDataAlreadyExists = errors.New("the data already exists")
	ErrDataNotExists     = errors.New("data does not exists")

	ErrNilReference      = errors.New("nil reference")
	ErrWithoutPermission = errors.New("without permission")
	ErrParameterMismatch = errors.New("parameter mismatch")
)
