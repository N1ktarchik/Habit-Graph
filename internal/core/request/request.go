package request

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeAndValidateJSON(r *http.Request, payload any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, payload); err != nil {
		return err
	}

	return nil
}
