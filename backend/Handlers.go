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

// ResponseCode Global codes for client - backend connections
type ResponseCode int

// ResponseCodes
const (
	Unauthorized        ResponseCode = 401
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
	profile, err := GetProfile(accessToken)
	if err != nil {
		respondWithError(w, FailedProfileFetch, err.Error())
	}
	json.NewEncoder(w).Encode(profile)
}

// Test returns a sample LinkedIn Profile JSON object
func Test(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tt := params["testType"]
	if tt == "profile" {
		json.NewEncoder(w).Encode(strings.Replace(sampleProfile, "\n", "", -1))
	} else if tt == "seedUser" {
		users, err := fbClient.Ref("/users")
		if err != nil {
			fmt.Println(err)
			respondWithError(w, FailedDBCall, err.Error())
			return
		}

		umap := make(map[string]User, len(cachedUsers))

		for _, u := range cachedUsers {
			u.Profile.FormattedName = u.Profile.FirstName + " " + u.Profile.LastName
			u.Profile.ID = u.ID
			umap[u.ID] = u
		}

		if err := users.Update(umap); err != nil {
			fmt.Println(err)
		}
	} else if tt == "distance" {
		var testLoc Geolocation
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &testLoc); err != nil {
			bodyString := string(bodyBytes)
			fmt.Println(bodyString)
			return // no match was returned
		} //40.4259° N, 86.9081° W
		res := InRadius(Geolocation{Long: -86.9081, Lat: 40.4259}, testLoc, 20)
		rj := make(map[string]bool)
		rj["res"] = res
		respondWithJSON(w, 200, rj)
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

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
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
		// fmt.Print("Trying to get matches")
		MatchList, err := GetMatches(userID, PMatchList)
		// fmt.Print("Got matches")
		// fmt.Print(MatchList)
		if err != nil {
			return // unable to get anyone to match
		}
		json.NewEncoder(w).Encode(MatchList)
		return
	}
	respondWithError(w, Unauthorized, "Unauthorized")
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
		return
	}
	respondWithError(w, Unauthorized, "Unauthorized")
}

// InitiateMeetover called to begin the meetover appointment betwen two users
func InitiateMeetover(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		initiatorId := r.Header.Get("Identity")
		params := mux.Vars(r)
		requestedId := params["otherId"]

		if err := AddThread(initiatorId, requestedId); err != nil {
			respondWithError(w, FailedDBCall, "Could not create chat thread")
			fmt.Println("Failed to create the thread " + initiatorId + ", " + requestedId)
			return
		}

		// Send a push notification to the requested user
		formattedName, err := fbClient.Ref("/users/" + initiatorId + "/profile/formattedName")
		if err != nil {
			json.NewEncoder(w).Encode("Could not send push notification")
			fmt.Println(err.Error())
			return
		}
		var name string
		if err = formattedName.Value(&name); err != nil {
			json.NewEncoder(w).Encode("Could not send push notification")
			fmt.Println(err.Error())
			return
		}

		title := "New MeetOver Request"
		body := name + " would like to MeetOver"
		pushNotification := PushNotification{
			ID:    requestedId,
			Title: title,
			Body:  body,
		}
		err = SendPushNotification(&pushNotification)
		if err != nil {
			json.NewEncoder(w).Encode("Could not send push notification")
			fmt.Println(err.Error())
			return
		}
	}
	json.NewEncoder(w).Encode("Success")
}

// SendPush sends sends a push notification for a verified user
func SendPush(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		var pushNotification PushNotification
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &pushNotification); err != nil {
			fmt.Println("Unable to send push notification")
			return
		}
		SendPushNotification(&pushNotification)
	}
	json.NewEncoder(w).Encode("Success")
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
