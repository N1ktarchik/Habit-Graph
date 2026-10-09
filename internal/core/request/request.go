package request

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(r *http.Request, payload any) error {
	return json.NewDecoder(r.Body).Decode(payload)
}

func GetID(r *http.Request) string {
	return r.Header.Get("id")
}

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
