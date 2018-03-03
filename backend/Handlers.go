package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Person is user on MeetOver
type Person struct {
	ID          string         `json:"uid,omitempty"`
	Firstname   string         `json:"firstName,omitempty"`
	Lastname    string         `json:"lastName,omitempty"`
	Address     *Address       `json:"address,omitempty"`
	AccessToken ATokenResponse `json:"accessToken"`
	LiProfile   LiProfile      `json:"profile"`
}

var people []Person

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
	json.NewEncoder(w).Encode(people)
}

// VerifyUser will get a code object to obtain an access token
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

	users, err := fbClient.Ref("/users")
	if err != nil {
		respondWithError(w, FailedDBCall, err.Error())
		fmt.Println("Failed to save user profile")
		return
	}
	person := Person{
		ID:          lip.ID,
		Firstname:   lip.FirstName,
		Lastname:    lip.LastName,
		AccessToken: aTokenResp,
		LiProfile:   lip,
	}
	addUser := make(map[string]interface{})
	addUser[lip.ID] = person
	defer users.Update(addUser)

	customToken, err := CreateCustomToken(lip.ID)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		fmt.Println("Sending failed token exchange error")
		return
	}

	var resp AuthResponse
	resp.AccessToken = aTokenResp
	resp.LiProfile = lip
	resp.FirebaseCustomToken = customToken

	json.NewEncoder(w).Encode(resp)
}

// Match will set a flag to notify the system the suer is matched
func Match(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
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
