package handler

import (
	"encoding/json"
	"net/http"
)

// HealthCheck 測試服務是否正常運作
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "success",
		"message": "GopherMind is thinking...",
		"version": "1.26.0",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
