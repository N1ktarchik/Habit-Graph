package response

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, m map[string]string) error {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)

	if err != nil {
		return err
	}
	return nil
}
