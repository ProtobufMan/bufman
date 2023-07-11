package e

import (
	"errors"
	"github.com/bufbuild/connect-go"
)

type ResponseError interface {
	Error() string
	Err() error
	Code() connect.Code
}

type BaseResponseError struct {
	msg  string
	code connect.Code
}

func NewBaseResponseError(msg string, code connect.Code) *BaseResponseError {
	return &BaseResponseError{
		msg:  msg,
		code: code,
	}
}

func (e *BaseResponseError) Error() string {
	if e == nil {
		return ""
	}
	return e.msg
}

func (e *BaseResponseError) Err() error {
	if e == nil {
		return errors.New("")
	}
	return errors.New(e.msg)
}

func (e *BaseResponseError) Code() connect.Code {
	if e == nil {
		return connect.CodeUnknown
	}

	return e.code
}
