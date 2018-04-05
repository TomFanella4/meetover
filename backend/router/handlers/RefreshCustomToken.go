package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"meetover/backend/firebase"
)

// RefreshCustomToken will refresh an authorized users Firebase custom token
func RefreshCustomToken(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		id := r.Header.Get("Identity")

		customToken, err := firebase.CreateCustomToken(id)
		if err != nil {
			respondWithError(w, FailedTokenExchange, err.Error())
			fmt.Println("Failed to create Firebase custom token")
			return
		}

		var resp RefreshResponse
		resp.FirebaseCustomToken = customToken

		json.NewEncoder(w).Encode(resp)
	}
}
