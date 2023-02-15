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
	// ConfigUnmarshallErr = AppError{
	// 	Message: "can't ",
	// }
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
