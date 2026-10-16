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

func BadJSON() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "invalid JSON",
	}
}

func BadDate() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "invalid date",
	}
}

func BadTitle() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "invalid title",
	}
}

func BadID() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "invalid ID",
	}
}

func BadYear() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "invalid year",
	}
}

func LongDescription() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "too long description",
	}
}

func NothingToUpdate() *ErrorApp {
	return &ErrorApp{
		Code:    http.StatusBadRequest,
		Message: "no data to update",
	}
}
