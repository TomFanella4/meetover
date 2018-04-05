package handlers

import (
	"encoding/json"
	"net/http"
)

// Index general success place holder
func Index(w http.ResponseWriter, r *http.Request) {
	response := ServerResponse{Message: "Welcome to MeetOverAPI", Success: true}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
