package request

import (
	"encoding/json"
	"net/http"

	"github.com/N1ktarchik/habbit-graph/internal/core/errors"
	"github.com/google/uuid"
)

func DecodeJSON(r *http.Request, payload any) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return errors.BadJSON()
	}
	return nil
}

func GetAndValidateID(r *http.Request) (string, error) {
	id := r.PathValue("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return "", errors.BadID()
	}

	return id, nil

}

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
