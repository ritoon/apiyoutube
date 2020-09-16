package db

import (
	"fmt"
	"net/http"
)

type ErrorDB struct {
	Code   int
	Params string
	Err    error
}

func NewError(param string, err error) *ErrorDB {
	return &ErrorDB{
		Code:   http.StatusInternalServerError,
		Params: param,
		Err:    err,
	}
}

func (e ErrorDB) Error() string {
	return fmt.Sprintf("db: params:%v err:%v", e.Params, e.Err)
}

func NewErrNotFound(params string, err error) error {
	return &ErrorDB{
		Code:   http.StatusNotFound,
		Params: params,
		Err:    err,
	}
}
