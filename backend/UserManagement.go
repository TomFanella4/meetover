package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var sampleProfile = `
{
  "currentShare": {
    "attribution": {"share": {
      "author": {
        "firstName": "MaryAnn",
        "id": "bxQrUhmst0",
        "lastName": "Gibney"
      },
      "comment": "Big thanks to the always vibrant Geoff Nyheim and Amazon Web Services colleagues Roshni Joshi, Michael Dowling ☁, Peter Tannenwald, Robert Heitzler, and @shane for coaching and mentoring DePaul University Center for Sales Leadership students on career and resume skills this afternoon! Plus DePaul alum Jeremy Paul for sharing your experience! #CSLyourself",
      "id": "s6362854320802713600"
    }},
    "author": {
      "firstName": "James",
      "id": "HJSNGIIRCj",
      "lastName": "Baker"
    },
    "comment": "MIT ",
    "id": "s6361684091376668672",
    "source": {"serviceProvider": {"name": "FLAGSHIP"}},
    "timestamp": 1516743681760,
    "visibility": {"code": "anyone"}
  },
  "emailAddress": "abc.google.com",
  "firstName": "James",
  "formattedName": "James Baker",
  "headline": "Software Engineering Manager",
  "id": "HJSNGIIRCj",
  "industry": "Computer Software",
  "lastName": "Baker",
  "location": {
    "country": {"code": "us"},
    "name": "United States"
  },
  "numConnections": 411,
  "numConnectionsCapped": false,
  "pictureUrl": "https://media.licdn.com/mpr/mprx/0_1lpYsrJfPnDiMc-pBPEPZ-hfnNriM6Cg9vEYxXza1zwCM91rsle0tT4fvFjCMnPlnbXKlGc70KuGVGW7JGJCtXs_sKu_VGBjsGJgML7mllMfrquT9P71JlTxVLnPPGh21TgtORGj196",
  "positions": {
    "_total": 1,
    "values": [{
      "company": {
        "id": 3846,
        "industry": "Higher Education",
        "name": "MIT",
        "size": "10001+",
        "type": "Educational"
      },
      "id": 827836295,
      "isCurrent": true,
      "location": {"name": "New Jersey, New York"},
      "startDate": {
        "month": 5,
        "year": 2016
      },
      "summary": "Current:\n\nDesigning a supply chain management system using a blockchain oriented approach. Implemented a prototype using the Hyperledger project by IBM and discovering interesting use cases and research aspects for the system.\nVisit https://freedom.cs.purdue.edu/ for more info.\n\nPrevious Project:\n\nWorked in a team to design a communication protocol using cryptography primitives and cryptocurrencies such as Bitcoin and Ether. Some key knowledge I acquired for the project includes cryptographic hashing, watermarking (robust, fragile and semi-fragile), oblivious transfer protocol, blockchain primitives, bitcoin scripts, smart contracts, traitor tracing and other related fields in order to design the intended system and provide adequate applications to it.",
      "title": "Research Assistant"
    }]
  },
  "summary": "Senior (Graduation May 2018) at MIT pursuing a BSc in Computer Science looking to apply my knowledge in the field of software development, system security and blockchain solutions through full-time opportunities."
}`

// ExchangeToken does the auhentication using client code and secret
func ExchangeToken(TempClientCode string) ATokenResponse {

	cid, found := os.LookupEnv("LI_CLIENT_ID")
	if !found {
		log.Fatal("Unable to get client id from env var")
	}
	csecret, found := os.LookupEnv("LI_CLIENT_SECRET")
	if !found {
		log.Fatal("Unable to get client secret from env var")
	}
	ruri, found := os.LookupEnv("LI_REDIRECT_URI")
	if !found {
		log.Fatal("Unable to get redirect uri from env var")
	}
	code := url.QueryEscape(TempClientCode)
	ruri = url.QueryEscape(ruri)
	content := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&"+
		"client_id=%s&client_secret=%s", code, ruri, cid, csecret)

	endpoint := "https://www.linkedin.com/oauth/v2/accessToken"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(content)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "www.linkedin.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)
	// log.Print(bodyString)

	// Fill the record with the data from the JSON response
	var record ATokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	return record
}

// GetLiProfile uses access token and REST call to get the user's linkedIn profile
func GetLiProfile(AccessToken string) LiProfile {
	// Fill the record with the data from the JSON
	var record LiProfile
	// QueryEscape escapes the parama
	items := "(id,first-name,last-name,maiden-name,formatted-name,phonetic-first-name," +
		"phonetic-last-name,formatted-phonetic-name,headline,location," +
		"industry,current-share,num-connections,num-connections-capped," +
		"summary,specialties,positions,picture-url,email-address)"

	at := url.QueryEscape(AccessToken)
	itm := url.QueryEscape(items)
	url := fmt.Sprintf("%s/v1/people/~:%s?oauth2_access_token=%s&format=json", LIAPI, itm, at)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return record
	}

	// For control over HTTP client headers,redirect policy, and other settings,
	client := &http.Client{}

	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("client.do failed ", err)
		return record
	}
	defer resp.Body.Close()
	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	// temprary place holder
	// Use json.Decode for reading streams of JSON data
	if err := json.Unmarshal([]byte(sampleProfile), &record); err != nil {
		log.Println(err)
	}
	return record
}
