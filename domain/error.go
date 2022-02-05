package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrNotFound            = errors.New("your requested Record not found")
	ErrBadParamInput       = errors.New("given Param is not valid")
)
