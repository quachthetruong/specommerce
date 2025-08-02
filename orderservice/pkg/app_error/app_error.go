package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Err     error
	Code    int
	Message string
}

type AppErrorOption func(appError *AppError)

func New(err error, opts ...AppErrorOption) AppError {
	appErr := AppError{
		Err:  err,
		Code: http.StatusInternalServerError,
	}
	for _, opt := range opts {
		opt(&appErr)
	}
	return appErr
}

func (e AppError) Error() string {
	if e.Err != nil {
		return e.Message + e.Err.Error()
	}
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Err
}

func WithCode(code int) AppErrorOption {
	return func(appError *AppError) {
		appError.Code = code
	}
}

func WithMessage(message string) AppErrorOption {
	return func(appError *AppError) {
		appError.Message = message
	}
}

func ErrParamInvalid(param string) AppError {
	return New(nil, WithCode(400_0001), WithMessage(fmt.Sprintf("invalid param: %s", param)))
}

var NotFoundIdWhenUpdate = New(nil, WithCode(404_0000), WithMessage("not found id when update"))
var NotFoundPrimaryKey = New(nil, WithCode(404_0001), WithMessage("not found primary key"))
