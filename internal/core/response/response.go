package response

import (
	"encoding/json"
	"net/http"

	"github.com/N1ktarchik/habbit-graph/internal/core/errors"
)

func ResponseOK(w http.ResponseWriter, statusCode int, payload any) error {

	resp, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	_, err = w.Write(resp)

	return err
}

func ResponseWithError(w http.ResponseWriter, err error) error {
	errorApp, ok := errors.IsErrorApp(err)
	if !ok {
		newErr := errors.UnknownErr()
		return ResponseOK(w, newErr.Code, newErr.Message)
	}

	return ResponseOK(w, errorApp.Code, errorApp.Message)
}
