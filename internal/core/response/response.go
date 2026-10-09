package response

import (
	"encoding/json"
	"net/http"

	"github.com/N1ktarchik/habbit-graph/internal/core/errors"
)

func responseWithJSON(w http.ResponseWriter, statusCode int, payload any) error {
	resp, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(resp)

	return err
}

func ResponseWithMap(w http.ResponseWriter, m map[string]string) error {
	w.Header().Set("Content-Type", "application/json")

	return responseWithJSON(w, http.StatusOK, m)
}

func ResponseWithError(w http.ResponseWriter, err error) error {
	errorApp, ok := errors.IsErrorApp(err)
	if !ok {
		newErr := errors.UnknownErr()
		return responseWithJSON(w, newErr.Code, newErr.Message)
	}

	return responseWithJSON(w, errorApp.Code, errorApp.Message)
}
