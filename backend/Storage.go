package main

// TODO:
// JSON schema for user profile
// IM schema and maintenance

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var fbApp *firebase.App
var fbClient *firego.Firebase

// GetProspectiveUsers Get the list of cachedUsers for matching in the area
func GetProspectiveUsers(coords Geolocation, radius int, lastUpdate int) ([]User, error) {
	// TODO:
	// returns a list of cachedUsers within radius of coords
	// that updated their location within lastUpdate hours
	n := len(cachedUsers)
	start := random(0, n/2)
	end := random((n/2)+1, n)
	return cachedUsers[start:end], nil
}

// random - helper for tests
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// InitializeFirebase reads API keys to use db
func InitializeFirebase() {
	// This is a Firebase workaround. The Firebase go library MUST read
	// certificate credentials from a file. Also, Heroku doesn't have static
	// storage. So, we must create the credential file dynamically from an
	// environment variable containing the json

	config, deployMode := os.LookupEnv("FIREBASE_CONFIG")
	if deployMode {
		fmt.Println("Firebase config file fetched from ENV var")
		err := ioutil.WriteFile("./firebase-config.json", []byte(config), 0644)
		if err != nil {
			log.Fatalf("Could not write config file")
		}
	} else {
		fmt.Println("Firebase config file fetched from local dir")
		buf, err := ioutil.ReadFile("./firebase-config.json")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		config = string(buf)
	}
	opt := option.WithCredentialsFile("./firebase-config.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	conf, err := google.JWTConfigFromJSON(
		[]byte(config),
		"https://www.googleapis.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email",
	)
	if err != nil {
		log.Fatalf("error initializing Google OAuth Config")
	}

	fbApp = app
	fbClient = firego.New("https://meetoverdb.firebaseio.com/", conf.Client(oauth2.NoContext))
}

// GetUser returns the user with uid in firebase
func GetUser(uid string) {
	return User{}
}

// CreateCustomToken Creates firebase based IM access token for the user with LinkedIn user ID
func CreateCustomToken(ID string) (string, error) {
	client, err := fbApp.Auth(context.Background())
	if err != nil {
		return "", errors.New("error getting Auth client")
	}

	token, err := client.CustomToken(ID)
	if err != nil {
		return "", errors.New("error minting custom token")
	}

	return token, nil
}

// CheckAuthorized hepler for RefreshCustomToken
func CheckAuthorized(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Token")
	id := r.Header.Get("Identity")

	if token == "" || id == "" {
		respondWithError(w, Unauthorized, "You are not authorized to make this request")
		return false
	}

	user, err := fbClient.Ref("/users/" + id + "/accessToken")
	if err != nil {
		fmt.Println("Error fetching user " + id + " from Firebase for authentication")
		respondWithError(w, FailedDBCall, err.Error())
		return false
	}

	var aToken map[string]interface{}
	if err := user.Value(&aToken); err != nil {
		fmt.Println("Error fetching value of user " + id + " from Firebase for authentication")
		respondWithError(w, FailedDBCall, err.Error())
		return false
	}

	if aToken["access_token"] == token {
		return true
	}

	respondWithError(w, Unauthorized, "You are not authorized to make this request")
	return false
}
