package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// LIAPI - linkedIn api site
// https://developer.linkedin.com/docs/fields/basic-profile
// https://developer.linkedin.com/docs/signin-with-linkedin
const LIAPI string = "https://api.linkedin.com"

// ATokenResponse from Code exchange with LI
type ATokenResponse struct {
	AToken string `json:"access_token"`
	Expiry uint   `json:"expires_in"`
}

// Profile is the JSON object we user to store our user data
type Profile struct {
	CurrentShare struct {
		Attribution struct {
			Share struct {
				Author struct {
					FirstName string `json:"firstName"`
					ID        string `json:"id"`
					LastName  string `json:"lastName"`
				} `json:"author"`
				Comment string `json:"comment"`
				ID      string `json:"id"`
			} `json:"share"`
		} `json:"attribution"`
		Author struct {
			FirstName string `json:"firstName"`
			ID        string `json:"id"`
			LastName  string `json:"lastName"`
		} `json:"author"`
		Comment string `json:"comment"`
		ID      string `json:"id"`
		Source  struct {
			ServiceProvider struct {
				Name string `json:"name"`
			} `json:"serviceProvider"`
		} `json:"source"`
		Timestamp  int64 `json:"timestamp"`
		Visibility struct {
			Code string `json:"code"`
		} `json:"visibility"`
	} `json:"currentShare"`
	EmailAddress  string `json:"emailAddress"`
	FirstName     string `json:"firstName"`
	FormattedName string `json:"formattedName"`
	Headline      string `json:"headline"`
	ID            string `json:"id"`
	Industry      string `json:"industry"`
	LastName      string `json:"lastName"`
	Location      struct {
		Country struct {
			Code string `json:"code"`
		} `json:"country"`
		Name string `json:"name"`
	} `json:"location"`
	NumConnections int    `json:"numConnections"`
	PictureURL     string `json:"pictureUrl"`
	Positions      struct {
		Total  int `json:"_total"`
		Values []struct {
			Company struct {
				ID       int    `json:"id"`
				Industry string `json:"industry"`
				Name     string `json:"name"`
				Size     string `json:"size"`
				Type     string `json:"type"`
			} `json:"company"`
			ID        int  `json:"id"`
			IsCurrent bool `json:"isCurrent"`
			Location  struct {
				Name string `json:"name"`
			} `json:"location"`
			StartDate struct {
				Month int `json:"month"`
				Year  int `json:"year"`
			} `json:"startDate"`
			Summary string `json:"summary"`
			Title   string `json:"title"`
		} `json:"values"`
	} `json:"positions"`
	Summary       string `json:"summary"`
	ShareLocation bool   `json:"shareLocation"`
	Greeting      string `json:"greeting"`
}

// User is user on MeetOver
type User struct {
	ID           string         `json:"uid,omitempty"`
	Location     Geolocation    `json:"location,omitempty"`
	AccessToken  ATokenResponse `json:"accessToken"`
	Profile      Profile        `json:"profile"`
	IsSearching  bool           `json:"isSearching"`
	IsMatchedNow bool           `json:"isMatched"` // set directly from the mobile app
}

// cachedUsers - temp var for testing
var cachedUsers []User

// ExchangeToken does the auhentication using client code and secret
func ExchangeToken(TempClientCode string, RedirectURI string) (ATokenResponse, error) {

	cid, found := os.LookupEnv("LI_CLIENT_ID")
	if !found {
		return ATokenResponse{}, errors.New("Unable to get client id from env var")
	}
	csecret, found := os.LookupEnv("LI_CLIENT_SECRET")
	if !found {
		return ATokenResponse{}, errors.New("Unable to get client secret from env var")
	}
	ruri, found := os.LookupEnv("LI_REDIRECT_URI")
	if found {
		fmt.Println("[+] Using server env redirect uri")
	} else {
		fmt.Println("[+] Using client param redirect uri")
		ruri = RedirectURI
	}

	code := url.QueryEscape(TempClientCode)
	ruri = url.QueryEscape(ruri)
	params := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&"+
		"client_id=%s&client_secret=%s", code, ruri, cid, csecret)

	endpoint := "https://www.linkedin.com/oauth/v2/accessToken"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(params)))
	if err != nil {
		return ATokenResponse{}, errors.New("Unable to create LI API call")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "www.linkedin.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ATokenResponse{}, errors.New("Unable to call LI API")
	}

	// Fill the record with the data from the JSON response
	var record ATokenResponse
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err := json.Unmarshal(bodyBytes, &record); err != nil {
		fmt.Println("Unexpected Token Exchange Response: ")
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return ATokenResponse{}, errors.New("Unable to LI token JSON")
	}

	if len(record.AToken) > 0 {
		return record, nil
		// stash to firebase save( data ( = token), flag)
	}
	return ATokenResponse{}, errors.New("No access token returned from linkedIn API")
}

// GetProfile uses access token and REST call to get the user's linkedIn profile
func GetProfile(AccessToken string) (Profile, error) {
	// Fill the record with the data from the JSON
	var record Profile
	// QueryEscape escapes the parama
	items := "(id,first-name,last-name,maiden-name,formatted-name,phonetic-first-name," +
		"phonetic-last-name,formatted-phonetic-name,headline,location," +
		"industry,current-share,num-connections,num-connections-capped," +
		"summary,specialties,positions,picture-url,email-address)"

	ATokenQE := url.QueryEscape(AccessToken)
	url := fmt.Sprintf("%s/v1/people/~:%s?oauth2_access_token=%s&format=json", LIAPI, items, ATokenQE)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Profile{}, errors.New("Unable to form HTTP request")
	}
	// For control over HTTP client headers,redirect policy, and other settings,
	client := &http.Client{}
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		return Profile{}, errors.New("Unable to make REST call to get profile data")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err := json.Unmarshal(bodyBytes, &record); err != nil {
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return Profile{}, errors.New("Unexpected Response Getting LI profile: " + bodyString)
	}

	// Set initial profile values
	record.ShareLocation = true

	return record, nil
}

// InitUser Updates access token if user exists or adds a new User as user in firebase
func InitUser(p Profile, aTokenResp ATokenResponse) (bool, error) {
	// Check if token exists
	userExists := false
	tokenRef, err := fbClient.Ref("/users/" + p.ID + "/accessToken")
	if err != nil {
		fmt.Println("Failed to save user profile to Firebase in InitUser()")
		return userExists, errors.New("Failed to save user profile: \n" + err.Error())
	}
	token := ATokenResponse{}
	if err := tokenRef.Value(&token); err != nil {
		log.Fatalf("Could not get user in Firebase: %v\n", err)
	}

	// Update token or create new user
	if token.AToken != "" {
		defer tokenRef.Set(aTokenResp)
		userExists = true
	} else {
		userRef, err := fbClient.Ref("/users/" + p.ID)
		if err != nil {
			fmt.Println("Failed to save user profile to Firebase in InitUser()")
			return false, errors.New("Failed to save user profile: \n" + err.Error())
		}

		user := User{
			ID:          p.ID,
			AccessToken: aTokenResp,
			Profile:     p,
		}
		defer userRef.Set(user)
	}
	return userExists, nil
}

// LoadTestUsers gets test users generated from a random data generator
func LoadTestUsers() {
	raw, err := ioutil.ReadFile("./data/FinalTestUsers.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var testUsers []User
	if err = json.Unmarshal(raw, &testUsers); err != nil {
		fmt.Println("[-] Unable to load test users")
		fmt.Println(err.Error())
		return
	}
	cachedUsers = append(cachedUsers, testUsers...)
}
