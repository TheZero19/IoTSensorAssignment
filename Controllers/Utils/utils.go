package Utils

import (
	"io"
	"net/http"
)

func CheckPOSTRequestValidity(w http.ResponseWriter, r *http.Request) (bool, []byte) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return false, nil
	}
	return true, body
}
