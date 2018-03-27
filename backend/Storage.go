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

const separator = "|"

// GetProspectiveUsers Get the list of people for matching in the area
func GetProspectiveUsers(coords Geolocation, radius int, lastUpdate int) ([]User, error) {
	// TODO:
	// returns a list of people within radius of coords
	// that updated their location within lastUpdate hours
	return people, nil
}

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

func CreateThreadForUser(ID1 string, threadId string, ID2 string) error {
	threadList, err := fbClient.Ref("/users/" + ID1 + "/threadList/" + threadId)
	if err != nil {
		return err
	}

	otherUserName, err := fbClient.Ref("/users/" + ID2 + "/profile/firstName")
	if err != nil {
		return err
	}

	var name string
	if err := otherUserName.Value(&name); err != nil {
		return err
	}

	threadInfo := map[string]interface{}{
		"_id":  threadId,
		"name": name,
	}
	if err := threadList.Update(threadInfo); err != nil {
		return err
	}

	return nil
}

func AddThread(P1 string, P2 string) error {
	ID1, ID2 := "", ""

	if P1 == P2 {
		return errors.New("Cannot start a thread with only one user")
	} else if P1 < P2 {
		ID1, ID2 = P1, P2
	} else {
		ID1, ID2 = P2, P1
	}

	threadId := ID1 + separator + ID2

	if err := CreateThreadForUser(ID1, threadId, ID2); err != nil {
		return err
	}
	if err := CreateThreadForUser(ID2, threadId, ID1); err != nil {
		return err
	}

	thread, err := fbClient.Ref("/messages/" + threadId + "/0")
	if err != nil {
		return err
	}

	initialMessage := map[string]interface{}{
		"_id":       0,
		"text":      "Welcome to the chat!",
		"createdAt": time.Now().UTC().Format(time.RFC3339),
		"system":    true,
	}
	if err := thread.Update(initialMessage); err != nil {
		return err
	}

	return nil
}

func addGeolocation(coord Geolocation) {
	addGeo := make(map[string]interface{})

	loc := coord
	// TODO: look for the user and add/update the
	// Geolocation json WITHIN the User struct

	addGeo[loc.ID] = loc
	geo, err := fbClient.Ref("/Geo")

	if err != nil {
		fmt.Println("Adding Geo to DB error")
		return
	}
	defer geo.Update(addGeo)
}
