package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ClientID - temp
var ClientID string

// Person is user on MeetOver
type Person struct {
	ID          string   `json:"id,omitempty"`
	Firstname   string   `json:"firstname,omitempty"`
	Lastname    string   `json:"lastname,omitempty"`
	Address     *Address `json:"address,omitempty"`
	AccessToken string   `json:"accesstoken"`
}

// Address is a our location metric
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	Area  string `json:"area,omitempty"`
}

// ServerResponse is what we give to our clients
type ServerResponse struct {
	Code    ResponseCode `json:"id"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
}

var people []Person

// ResponseCode Global codes for client - backend connections
type ResponseCode int

const (
	FailedTokenExchange ResponseCode = 506
	FailedDBCall        ResponseCode = 507
	FailedProfileFetch  ResponseCode = 508
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
	// params := mux.Vars(r)
	//  location = getFromDataSet( params["coords"].split(",") )
	location := Address{City: "Chicago", State: "IL", Area: "ORD"}
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
	fmt.Println(tempUserCode)
	aTokenResp, err := ExchangeToken(tempUserCode)
	if err != nil {
		respondWithError(w, FailedTokenExchange, err.Error())
		fmt.Println("Sending failed token exchange error")
	} else {
		json.NewEncoder(w).Encode(aTokenResp)
		fmt.Println(aTokenResp.AToken)
	}
}

// Match will set a flag to notify the system the suer is matched
func Match(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
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
