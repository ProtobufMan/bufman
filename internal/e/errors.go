package e

import (
	"fmt"
	"github.com/bufbuild/connect-go"
)

type UnknownError struct {
	*BaseResponseError
}

func NewUnknownError(procedure string) *UnknownError {
	msg := fmt.Sprintf("unknown error occurs while %s", procedure)
	return &UnknownError{
		NewBaseResponseError(msg, connect.CodeUnknown),
	}
}

type NotFoundError struct {
	*BaseResponseError
}

func NewNotFoundError(object string) *NotFoundError {
	msg := fmt.Sprintf("%s not found", object)
	return &NotFoundError{
		NewBaseResponseError(msg, connect.CodeNotFound),
	}
}

type AlreadyExistsError struct {
	*BaseResponseError
}

func NewAlreadyExistsError(object string) *AlreadyExistsError {
	msg := fmt.Sprintf("%s is already exists", object)
	return &AlreadyExistsError{
		NewBaseResponseError(msg, connect.CodeAlreadyExists),
	}
}

type PermissionDeniedError struct {
	*BaseResponseError
}

func NewPermissionDeniedError(procedure string) *PermissionDeniedError {
	msg := fmt.Sprintf("%s is permission denied", procedure)
	return &PermissionDeniedError{
		NewBaseResponseError(msg, connect.CodePermissionDenied),
	}
}

type InternalError struct {
	*BaseResponseError
}

func NewInternalError(procedure string) *InternalError {
	msg := fmt.Sprintf("internal error occurs while %s", procedure)
	return &InternalError{
		NewBaseResponseError(msg, connect.CodeInternal),
	}
}

type CodeUnauthenticatedError struct {
	*BaseResponseError
}

func NewUnauthenticatedError(procedure string) *CodeUnauthenticatedError {
	msg := fmt.Sprintf("unauthenticated error occurs while %s", procedure)
	return &CodeUnauthenticatedError{
		NewBaseResponseError(msg, connect.CodeUnauthenticated),
	}
}

type InvalidArgumentError struct {
	*BaseResponseError
}

func NewInvalidArgumentError(object string) *InvalidArgumentError {
	msg := fmt.Sprintf("%s is invalid", object)
	return &InvalidArgumentError{
		NewBaseResponseError(msg, connect.CodeAlreadyExists),
	}
}
