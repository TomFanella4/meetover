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
	LiProfile           LiProfile      `json:"profile"`
	AccessToken         ATokenResponse `json:"token"`
	FirebaseCustomToken string         `json:"firebaseCustomToken"`
}

// ResponseCode Global codes for client - backend connections
type ResponseCode int

const (
	FailedTokenExchange ResponseCode = 506
	FailedDBCall        ResponseCode = 507
	FailedProfileFetch  ResponseCode = 508
	FailedLocationQuery ResponseCode = 509
	FailedUserInit      ResponseCode = 510
)

// GetUserProfile will give back a json object of user's LinkedIn Profile
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accessToken := params["accessToken"]
	profile, err := GetLiProfile(accessToken)
	if err != nil {
		respondWithError(w, FailedProfileFetch, err.Error())
	}
	json.NewEncoder(w).Encode(profile)
}

// Test returns a sample LinkedIn Profile JSON object
func Test(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tt := params["testType"]
	if tt == "liprofile" {
		json.NewEncoder(w).Encode(strings.Replace(sampleProfile, "\n", "", -1))
	}
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

// GetPeople returns all users
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cachedUsers)
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

	lip, err := GetLiProfile(aTokenResp.AToken)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		fmt.Println("Sending failed token exchange error")
		return
	}
	// Updates access token if user exists or adds a new User
	err = InitUser(lip, aTokenResp)
	if err != nil {
		respondWithError(w, FailedUserInit, err.Error())
		return
	}

	// gets firebase access token for user's IM chat
	customToken, err := CreateCustomToken(lip.ID)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		return
	}
	var resp AuthResponse
	resp.AccessToken = aTokenResp
	resp.LiProfile = lip
	resp.FirebaseCustomToken = customToken

	json.NewEncoder(w).Encode(resp)
}

// Match will set a flag to notify the system the user is matched
func Match(w http.ResponseWriter, r *http.Request) {
	// requesters uid
	params := mux.Vars(r)
	userID := params["uid"]
	// get requester's coords
	var userLocation Geolocation
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err := json.Unmarshal(bodyBytes, &userLocation); err != nil {
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return // no match was returned
	}
	radius := 1     //1 km
	lastUpdate := 2 // 2hrs
	PMatchList, err := GetProspectiveUsers(userLocation, radius, lastUpdate)
	if err != nil {
		return // unable to get anyone from db
	}
	MatchList, err := GetMatches(userID, PMatchList)
	if err != nil {
		return // unable to get anyone to match
	}
	json.NewEncoder(w).Encode(MatchList)
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
