package main

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

// LIAPI - linkedIn api site
// https://developer.linkedin.com/docs/fields/basic-profile
// https://developer.linkedin.com/docs/signin-with-linkedin
const LIAPI string = "https://api.linkedin.com"

// ATokenResponse from Code exchange with LI
type ATokenResponse struct {
	AToken string `json:"access_token"`
	Expiry uint   `json:"expires_in"`
}

// LiProfile is the JSON object we get from the linkedin profile
type LiProfile struct {
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
	NumConnections       int    `json:"numConnections"`
	NumConnectionsCapped bool   `json:"numConnectionsCapped"`
	PictureURL           string `json:"pictureUrl"`
	Positions            struct {
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
	Summary string `json:"summary"`
}

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

// GetLiProfile uses access token and REST call to get the user's linkedIn profile
func GetLiProfile(AccessToken string) (LiProfile, error) {
	// Fill the record with the data from the JSON
	var record LiProfile
	// QueryEscape escapes the parama
	items := "(id,first-name,last-name,maiden-name,formatted-name,phonetic-first-name," +
		"phonetic-last-name,formatted-phonetic-name,headline,location," +
		"industry,current-share,num-connections,num-connections-capped," +
		"summary,specialties,positions,picture-url,email-address)"

	ATokenQE := url.QueryEscape(AccessToken)
	url := fmt.Sprintf("%s/v1/people/~:%s?oauth2_access_token=%s&format=json", LIAPI, items, ATokenQE)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LiProfile{}, errors.New("Unable to form HTTP request")
	}
	// For control over HTTP client headers,redirect policy, and other settings,
	client := &http.Client{}
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		return LiProfile{}, errors.New("Unable to make REST call to get profile data")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err := json.Unmarshal(bodyBytes, &record); err != nil {
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return LiProfile{}, errors.New("Unexpected Response Getting LI profile: " + bodyString)
	}
	return record, nil
}
