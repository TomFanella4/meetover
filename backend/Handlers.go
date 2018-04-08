package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ServerResponse - Error message JSON structure
type ServerResponse struct {
	Code    ResponseCode `json:"id"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
}

// AuthResponse is the JSON returned to client during login to backend
type AuthResponse struct {
	Profile             Profile        `json:"profile"`
	AccessToken         ATokenResponse `json:"token"`
	FirebaseCustomToken string         `json:"firebaseCustomToken"`
	UserExists          bool           `json:"userExists"`
}

// RefreshResponse - returned to the client during firebase custom token refresh
type RefreshResponse struct {
	FirebaseCustomToken string `json:"firebaseCustomToken"`
}

// MatchResponse returned to the UI when /match is hit
type MatchResponse struct {
	Matches []MatchValue `json:"matches"`
}

// UserProfileResponse returned when requesting a user profile
type UserProfileResponse struct {
	Profile  Profile     `json:"profile"`
	Location Geolocation `json:"location"`
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

// GetUserProfiles will give back a json object of user's Profile
func GetUserProfiles(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		var ids []string
		profiles := make(map[string]UserProfileResponse)
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &ids); err != nil {
			fmt.Println(err.Error())
			respondWithError(w, FailedProfileFetch, "Unable to get user profiles")
			return
		}

		for _, id := range ids {
			user, err := GetUser(id)
			if err != nil {
				fmt.Println(err.Error())
				respondWithError(w, FailedProfileFetch, "Failed to get User")
				return
			}
			var profile UserProfileResponse
			profile.Profile = user.Profile
			profile.Location = user.Location
			profiles[id] = profile
		}
		json.NewEncoder(w).Encode(profiles)
	}
}

// Test returns a sample LinkedIn Profile JSON object
func Test(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tt := params["testType"]
	if tt == "profile" {
		json.NewEncoder(w).Encode(strings.Replace(sampleProfile, "\n", "", -1))
	}
	resp := ServerResponse{Success, "Success", true}
	json.NewEncoder(w).Encode(resp)
}

// GetAddress will give back a json object based on coordinates
func GetAddress(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	location, err := QueryLocation(params["coords"])
	if err != nil {
		respondWithError(w, FailedLocationQuery, err.Error())
	}
	json.NewEncoder(w).Encode(location)
}

// Index general success place holder
func Index(w http.ResponseWriter, r *http.Request) {
	response := ServerResponse{Message: "Welcome to MeetOverAPI", Success: true}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// VerifyUser - token exchange and authentication at user login
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tempUserCode := params["code"]
	redirectURI := r.URL.Query().Get("redirect_uri")
	fmt.Println("[+] Recieved code: " + tempUserCode)
	fmt.Println("[+] Recieved redirect_uri: " + redirectURI)

	aTokenResp, err := ExchangeToken(tempUserCode, redirectURI)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		fmt.Println("Sending failed token exchange error")
		return
	}
	fmt.Println("[+] After ExchangeToken: " + aTokenResp.AToken)

	p, err := GetProfile(aTokenResp.AToken)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		fmt.Println("Sending failed token exchange error")
		return
	}
	// Updates access token if user exists or adds a new User
	userExists, err := InitUser(p, aTokenResp)
	if err != nil {
		respondWithError(w, FailedUserInit, err.Error())
		return
	}

	// gets firebase access token for user's IM chat
	customToken, err := CreateCustomToken(p.ID)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		return
	}
	var resp AuthResponse
	resp.AccessToken = aTokenResp
	resp.Profile = p
	resp.FirebaseCustomToken = customToken
	resp.UserExists = userExists

	json.NewEncoder(w).Encode(resp)
}

// Match will set a flag to notify the system the user is matched
func Match(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		// requesters uid
		params := mux.Vars(r)
		userID := params["uid"]
		// get requester's coords
		user, err := GetUser(userID)
		if err != nil {
			fmt.Println(err)
			return // Unable to get location
		}
		radius := 1     //1 km
		lastUpdate := 2 // 2hrs
		PMatchList, err := GetProspectiveUsers(&user.Location, radius, lastUpdate)
		if err != nil {
			fmt.Println(err)
			return // unable to get anyone from db
		}
		MatchList, err := GetMatches(userID, PMatchList)
		if err != nil {
			fmt.Println(err)
			return // unable to get anyone to match
		}
		json.NewEncoder(w).Encode(MatchList)
	}
}

// RefreshCustomToken will refresh an authorized users Firebase custom token
func RefreshCustomToken(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		id := r.Header.Get("Identity")

		customToken, err := CreateCustomToken(id)
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

// InitiateMeetover called to begin the meetover appointment betwen two users
func InitiateMeetover(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		initiatorID := r.Header.Get("Identity")
		params := mux.Vars(r)
		requestedID := params["otherID"]

		if err := AddThread(initiatorID, requestedID); err != nil {
			respondWithError(w, FailedDBCall, "Could not create chat thread")
			fmt.Println("Failed to create the thread " + initiatorID + ", " + requestedID)
			return
		}

		// Send a push notification to the requested user
		formattedName, err := fbClient.Ref("/users/" + initiatorID + "/profile/formattedName")
		if err != nil {
			respondWithError(w, FailedDBCall, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}
		var name string
		if err = formattedName.Value(&name); err != nil {
			respondWithError(w, FailedDBCall, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}

		title := "New MeetOver Request"
		body := name + " would like to MeetOver"
		pushNotification := PushNotification{
			ID:    requestedID,
			Title: title,
			Body:  body,
		}
		err = SendPushNotification(&pushNotification)
		if err != nil {
			respondWithError(w, FailedSendPush, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}
		resp := ServerResponse{Success, "Success", true}
		json.NewEncoder(w).Encode(resp)
	}
}

// ProcessDecision Updates the status of the specified user pair
func ProcessDecision(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		requestedID := r.Header.Get("Identity")
		params := mux.Vars(r)
		initiatorID := params["otherID"]

		var decision MeetOverDecisionBody
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &decision); err != nil {
			fmt.Println(err.Error())
			respondWithError(w, FailedMeetOver, "Failed to set MeetOver decision")
			return
		}

		if err := SetThreadStatus(requestedID, decision.ThreadID, decision.Status); err != nil {
			fmt.Println(err.Error())
			respondWithError(w, FailedMeetOver, "Failed to set MeetOver decision")
			return
		}
		if err := SetThreadStatus(initiatorID, decision.ThreadID, decision.Status); err != nil {
			fmt.Println(err.Error())
			respondWithError(w, FailedMeetOver, "Failed to set MeetOver decision")
			return
		}

		// Send a push notification to the initiator
		formattedName, err := fbClient.Ref("/users/" + requestedID + "/profile/formattedName")
		if err != nil {
			respondWithError(w, FailedDBCall, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}
		var name string
		if err = formattedName.Value(&name); err != nil {
			respondWithError(w, FailedDBCall, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}

		pushNotification := PushNotification{
			ID:    initiatorID,
			Title: name + " " + decision.Status + " request",
			Body:  "",
		}
		if err := SendPushNotification(&pushNotification); err != nil {
			fmt.Println("Unable to send push notification")
			respondWithError(w, FailedSendPush, "Could not send push notification")
			return
		}
		resp := ServerResponse{Success, "Success", true}
		json.NewEncoder(w).Encode(resp)
	}
}

// SendPush sends sends a push notification for a verified user
func SendPush(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		var pushNotification PushNotification
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &pushNotification); err != nil {
			fmt.Println("Unable to send push notification")
			respondWithError(w, FailedSendPush, "Could not send push notification")
			return
		}
		err := SendPushNotification(&pushNotification)
		if err != nil {
			respondWithError(w, FailedSendPush, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}
		resp := ServerResponse{Success, "Success", true}
		json.NewEncoder(w).Encode(resp)
	}
}

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
