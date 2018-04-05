package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

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
