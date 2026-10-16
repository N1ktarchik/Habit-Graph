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
		Message: "internal service error",
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

func BadTitle() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "invalide title",
	}
}

func BadID() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "ivalid ID",
	}
}

func BadYear() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "ivalid year",
	}
}

func LongDescription() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "to long description",
	}
}

func NothingToUpdate() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "nothing data for update",
	}
}
