package user

import (
	"meetover/backend/location"
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
	ID           string               `json:"uid,omitempty"`
	Location     location.Geolocation `json:"location,omitempty"`
	AccessToken  ATokenResponse       `json:"accessToken"`
	Profile      Profile              `json:"profile"`
	IsSearching  bool                 `json:"isSearching"`
	IsMatchedNow bool                 `json:"isMatched"` // set directly from the mobile app
}
