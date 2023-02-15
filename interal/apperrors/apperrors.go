package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Message  string
	Code     string
	HTTPCode int
}

var (
	ConfigUnmarshallErr = AppError{
		Message:  "couldn't unmarshal a response",
		Code:     "UNMARSHAL_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	ConfigReadErr = AppError{
		Message:  "couldn't read config",
		Code:     "CONFIG_READ_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	UserNotFoundErr = AppError{
		Message:  "can't find user",
		Code:     "SIGN_IN_ERR",
		HTTPCode: http.StatusBadRequest,
	}

	CanNotCreateTokenErr = AppError{
		Message:  "can't create token",
		Code:     "TOKEN_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	CanNotCreateUserErr = AppError{
		Message:  "can't create user",
		Code:     "SING_UP_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	HashingPasswordErr = AppError{
		Message:  "it is'nt hashing pasword",
		Code:     "PASSWORD_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	CanNotDeleteUserErr = AppError{
		Message:  "can not delete the user",
		Code:     "DELETION_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	PaginationErr = AppError{
		Message:  "it is pagination problem",
		Code:     "PAGINATION_ERR",
		HTTPCode: http.StatusInternalServerError,
	}
)

func (appError *AppError) Error() string {
	return appError.Code + ": " + appError.Message
}

func (appError *AppError) AppendMessage(anyErrs ...interface{}) *AppError {
	return &AppError{
		Message: fmt.Sprintf("%v : %v", appError.Message, anyErrs),
		Code:    appError.Code,
	}
}
