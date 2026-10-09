package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorApp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorApp) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func IsErrorApp(err error) (*ErrorApp, bool) {
	var appErr *ErrorApp
	if errors.As(err, &appErr) {
		return appErr, true
	}

	return nil, false
}

func UnknownErr() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusInternalServerError,
		Message: "unknown error",
	}
}

func BadValue() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "invalid value",
	}
}

func BadDate() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "invalide date",
	}
}
