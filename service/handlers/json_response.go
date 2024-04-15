package handlers

import (
	"encoding/json"
	"net/http"
)

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}
