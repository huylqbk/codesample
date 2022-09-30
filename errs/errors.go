package errs

import (
	"net/http"

	"github.com/pkg/errors"
)

type ErrorMessage error

var (
	ErrorForbidden        ErrorMessage = errors.New("you are not allowed to access this resource")
	ErrorInvalidRequest   ErrorMessage = errors.New("invalid request")
	ErrorIncorrectData    ErrorMessage = errors.New("incorrect data request")
	ErrorResourceNotFound ErrorMessage = errors.New("resource not found")
	ErrorSomethingWrong   ErrorMessage = errors.New("something wrong happened")
	ErrorServerFailure    ErrorMessage = errors.New("server failure")
	ErrorUnauthorized     ErrorMessage = errors.New("unauthorized request")
	ErrorRedisConnection  ErrorMessage = errors.New("redis connection error")
	ErrorTransaction      ErrorMessage = errors.New("transaction error")
	ErrorCreateResource   ErrorMessage = errors.New("create resource error")
	ErrorUpdateResource   ErrorMessage = errors.New("update resource error")
	ErrorAccessResource   ErrorMessage = errors.New("access resource error")
	ErrorDeleteResource   ErrorMessage = errors.New("delete resource error")
)

func ToCode(err error) int {
	code := http.StatusInternalServerError
	cause := errors.Cause(err)
	switch cause {
	case ErrorForbidden:
		code = http.StatusForbidden
	case ErrorInvalidRequest, ErrorIncorrectData, ErrorResourceNotFound:
		code = http.StatusBadRequest
	case ErrorServerFailure, ErrorSomethingWrong, ErrorTransaction, ErrorRedisConnection:
		code = http.StatusInternalServerError
	}
	return code
}
