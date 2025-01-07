package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIError struct {
	Message string `json:"message"`
	status  int
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(payload); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, http.StatusText(500), 500)
	}
}
