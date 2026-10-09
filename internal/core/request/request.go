package request

import (
	"encoding/json"
	"net/http"
)

func DecodeAndValidateJSON(r *http.Request, payload any) error {
	return json.NewDecoder(r.Body).Decode(payload)
}
