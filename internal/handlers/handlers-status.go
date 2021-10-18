package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (m *Repository) StatusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := make(map[string]string)
	currentStatus["Status"] = "Available"

	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
