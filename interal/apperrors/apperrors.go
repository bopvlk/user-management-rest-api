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
		Code:     "DATA_NOT_FOUND_ERR",
		HTTPCode: http.StatusBadRequest,
	}

	CanNotCreateTokenErr = AppError{
		Message:  "can't create token",
		Code:     "TOKEN_CREATE_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	CanNotParseTokenErr = AppError{
		Message:  "can't parse token",
		Code:     "TOKEN_PARSE_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	CanNotCreateUserErr = AppError{
		Message:  "can't create user",
		Code:     "SING_UP_ERR",
		HTTPCode: http.StatusBadRequest,
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
		Message:  "there is pagination problem",
		Code:     "PAGINATION_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	CanNotInitializeDBSessionErr = AppError{
		Message:  "can't initialize db session",
		Code:     "DB_SESSION_INIT_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	ValidatorErr = AppError{
		Message:  "validation cannot be passed",
		Code:     "VALIDATOR_ERR",
		HTTPCode: http.StatusBadRequest,
	}

	ValidatorInitializeErr = AppError{
		Message:  "can't initialize validation",
		Code:     "VALIDATOR_INIT_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	CanNotBindErr = AppError{
		Message:  "couldn't bind some data",
		Code:     "BINDING_ERR",
		HTTPCode: http.StatusBadRequest,
	}
	SomeCookieErr = AppError{
		Message:  "couldn't through out cookie",
		Code:     "COOKIE_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	WrongRoleErr = AppError{
		Message:  "you couldn't do this request. you shoud change the role",
		Code:     "ROLE_ERR",
		HTTPCode: http.StatusForbidden,
	}

	CanNotUpdateErr = AppError{
		Message:  "couldn't update the user",
		Code:     "UPDATE_ERR",
		HTTPCode: http.StatusInternalServerError,
	}

	CanNotRateYorself = AppError{
		Message:  "U can't rate yourself",
		Code:     "RATE_YORSELF_ERR",
		HTTPCode: http.StatusForbidden,
	}

	ProblemWithGivingRating = AppError{
		Message:  "Some problem with rate, probably you rate too often. Please try again later",
		Code:     "RATE_ERR",
		HTTPCode: http.StatusForbidden,
	}

	CanNotRateAgain = AppError{
		Message:  "You can not rate this user again",
		Code:     "RATE_USER_AGAIN_ERR",
		HTTPCode: http.StatusForbidden,
	}

	WrongTextInRateRequest = AppError{
		Message:  "Wrong the field in RateRequest. You must fill the field like \"up\", \"rm\" or \"down\"",
		Code:     "WRONG_TEXT_IN_RATE_REQUEST",
		HTTPCode: http.StatusBadRequest,
	}

	CanNotCreateWhoRatedErr = AppError{
		Message:  "can't create WhoRated",
		Code:     "CAN_NOT_CREATE_TABLE_WHORATED_ERR",
		HTTPCode: http.StatusInternalServerError,
	}
)

func (appError *AppError) Error() string {
	return appError.Code + ": " + appError.Message
}

func (appError *AppError) AppendMessage(anyErrs ...interface{}) *AppError {
	return &AppError{
		Message:  fmt.Sprintf("%v : %v", appError.Message, anyErrs),
		Code:     appError.Code,
		HTTPCode: appError.HTTPCode,
	}
}

func Is(err1 error, err2 *AppError) bool {
	err, ok := err1.(*AppError)
	if !ok {
		return false
	}

	return err.Code == err2.Code
}
