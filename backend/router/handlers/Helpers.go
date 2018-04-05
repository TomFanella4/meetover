package handlers

import (
	// "meetover/backend/matching"
	"encoding/json"
	"fmt"
	"meetover/backend/firebase"
	"meetover/backend/location"
	"meetover/backend/matching"
	"meetover/backend/user"
	"net/http"
)

// ServerResponse - Error message JSON structure
type ServerResponse struct {
	Code    ResponseCode `json:"id"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
}

// AuthResponse is the JSON returned to client during login to backend
type AuthResponse struct {
	Profile             user.Profile        `json:"profile"`
	AccessToken         user.ATokenResponse `json:"token"`
	FirebaseCustomToken string              `json:"firebaseCustomToken"`
	UserExists          bool                `json:"userExists"`
}

// RefreshResponse - returned to the client during firebase custom token refresh
type RefreshResponse struct {
	FirebaseCustomToken string `json:"firebaseCustomToken"`
}

// MatchResponse returned to the UI when /match is hit
type MatchResponse struct {
	Matches []matching.MatchValue `json:"matches"`
}

// UserProfileResponse returned when requesting a user profile
type UserProfileResponse struct {
	Profile  user.Profile         `json:"profile"`
	Location location.Geolocation `json:"location"`
}

// MeetOverDecisionBody processed for a meetover decision
type MeetOverDecisionBody struct {
	Status   string `json:"status"`
	ThreadID string `json:"_id"`
}

// ResponseCode Global codes for client - backend connections
type ResponseCode int

// ResponseCodes
const (
	Success             ResponseCode = 200
	Unauthorized        ResponseCode = 401
	FailedTokenExchange ResponseCode = 506
	FailedDBCall        ResponseCode = 507
	FailedProfileFetch  ResponseCode = 508
	FailedLocationQuery ResponseCode = 509
	FailedUserInit      ResponseCode = 510
	FailedSendPush      ResponseCode = 511
	FailedMeetOver      ResponseCode = 512
)

func respondWithError(w http.ResponseWriter, code ResponseCode, message string) {
	content := ServerResponse{Success: false, Message: message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(code))
	if err := json.NewEncoder(w).Encode(content); err != nil {
		fmt.Println("Unable to respond with JSON")
	}
}

func respondWithJSON(w http.ResponseWriter, code ResponseCode, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(code))
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		fmt.Println("Unable to respond with JSON")
	}
}

// CheckAuthorized checks if a user is authorized to make a request
func CheckAuthorized(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Token")
	id := r.Header.Get("Identity")

	if token == "" || id == "" {
		respondWithError(w, Unauthorized, "You are not authorized to make this request")
		return false
	}

	u, err := firebase.GetUser(id)
	if err != nil {
		fmt.Println("Error fetching value of user " + id + " from Firebase for authentication")
		respondWithError(w, FailedDBCall, err.Error())
		return false
	}

	if u.AccessToken.AToken == token {
		return true
	}

	respondWithError(w, Unauthorized, "You are not authorized to make this request")
	return false
}
