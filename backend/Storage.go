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
	"os"
	"time"

	firebase "firebase.google.com/go"
	"gopkg.in/zabawaba99/firego.v1"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var fbApp *firebase.App
var fbClient *firego.Firebase

// InitializeFirebase reads API keys to use db
func InitializeFirebase() {
	// This is a Firebase workaround. The Firebase go library MUST read
	// certificate credentials from a file. Also, Heroku doesn't have static
	// storage. So, we must create the credential file dynamically from an
	// environment variable containing the json
	config := []byte(os.Getenv("FIREBASE_CONFIG"))
	err := ioutil.WriteFile("./firebase-config.json", config, 0644)
	if err != nil {
		log.Fatalf("Could not write config file")
	}

	opt := option.WithCredentialsFile("./firebase-config.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	conf, err := google.JWTConfigFromJSON(
		config,
		"https://www.googleapis.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email",
	)
	if err != nil {
		log.Fatalf("error initializing Google OAuth Config")
	}

	fbApp = app
	fbClient = firego.New("https://meetoverdb.firebaseio.com/", conf.Client(oauth2.NoContext))
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

func addGeolocation(uid string, coord Coord, timeStamp int64) {
	addGeo := make(map[string]interface{})

	loc := Geolocation{
		ID:        uid,
		Coord:     coord,
		TimeStamp: timeStamp,
	}

	addGeo[loc.ID] = loc
	geo, err := fbClient.Ref("/Geo")

	if err != nil {
		fmt.Println("Adding Geo to DB error")
		return
	}
	defer geo.Update(addGeo)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
