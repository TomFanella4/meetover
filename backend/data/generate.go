package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

func updateJSONFile(newJSON interface{}, fileName string) {
	bytes, err := json.Marshal(newJSON)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = ioutil.WriteFile(fileName, bytes, 0644)
	return
}

func getUsers(rawFile string) []User {
	raw, err := ioutil.ReadFile(rawFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var users []User
	json.Unmarshal(raw, &users)
	return users
}

func generateTestUsers(rawFile string, sinkFile string) {
	csvFile, err := os.Open("./jobs.csv")
	if err != nil {
		fmt.Println(err.Error())
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	headers, error := reader.Read() // headers
	if error == io.EOF {
		fmt.Println("OEF in dataset")
	}
	fmt.Println("Headers: ")
	fmt.Println(headers) // 1 - desc, 3 - title, 4 - skills
	r := rand.Intn
	users := getUsers(rawFile)
	textLength := 250
	for i, u := range users {
		line, error := reader.Read()
		if error == io.EOF {
			fmt.Println("OEF in dataset")
		}
		description := line[1]
		title := line[2]
		skills := line[3]
		n := len(description)
		profile := u.Profile
		start := r(n-1) / 2
		greeting := description[start:]
		if len(greeting) > textLength {
			profile.Greeting = greeting[:textLength]
		} else {
			profile.Greeting = greeting
		}
		profile.Headline = title
		profile.Industry = strings.Replace(skills, ",", " , ", -1)
		profile.FormattedName = profile.FirstName + " " + profile.LastName
		start = r(n - 1)
		summary := description[start:]
		if len(summary) > textLength {
			profile.Summary = summary[:textLength]
		} else {
			profile.Summary = summary
		}
		users[i] = u
		users[i].Profile = profile
	}
	fmt.Println(users[25])
	updateJSONFile(users, sinkFile)
}

// Geolocation - latitide and longitude and last time of update
type Geolocation struct {
	Lat       float64 `json:"lat,omitempty"`
	Long      float64 `json:"long,omitempty"`
	TimeStamp int64   `json:"timestamp,omitempty"`
}

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
	Location     *Geolocation   `json:"location,omitempty"`
	AccessToken  ATokenResponse `json:"accessToken"`
	Profile      Profile        `json:"profile"`
	IsSearching  bool           `json:"isSearching"`
	IsMatchedNow bool           `json:"isMatched"` // set directly from the mobile app
}
