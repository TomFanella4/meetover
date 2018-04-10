package handlers

import (
	"encoding/json"
	"meetover/backend/firebase"
	"net/http"

	"github.com/gorilla/mux"
)

const allowImitate = true

// ImitateUser allows someone to sign in to an existing user for demo purposes
func ImitateUser(w http.ResponseWriter, r *http.Request) {
	if allowImitate {
		params := mux.Vars(r)
		uid := params["uid"]

		imitate, err := firebase.GetUser(uid)
		if err != nil {
			respondWithError(w, FailedDBCall, "Error fetching user "+uid+" from Firebase for imitation.")
			return
		}

		customToken, err := firebase.CreateCustomToken(uid)
		if err != nil {
			respondWithError(w, FailedTokenExchange, err.Error())
			return
		}

		resp := AuthResponse{
			AccessToken:         imitate.AccessToken,
			Profile:             imitate.Profile,
			FirebaseCustomToken: customToken,
			UserExists:          true,
		}

		json.NewEncoder(w).Encode(resp)
		return
	}

	respondWithError(w, Unauthorized, "User imitation disabled")
}
