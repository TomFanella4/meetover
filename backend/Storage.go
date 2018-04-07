package main

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
	"github.com/zabawaba99/firego"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var fbApp *firebase.App
var fbClient *firego.Firebase

const separator = "|"

// GetProspectiveUsers Get the list of people for matching in the area
func GetProspectiveUsers(coords *Geolocation, radius int, lastUpdate int) ([]User, error) {
	// TODO:
	// Filter users by Geolocation and radius
	userMap := map[string]User{}
	users := []User{}
	userRef, err := fbClient.Ref("/users")
	if err != nil {
		return []User{}, err
	}
	if err := userRef.Value(&userMap); err != nil {
		return []User{}, err
	}

	for k := range userMap {
		users = append(users, userMap[k])
	}

	return users, nil
}

// random - helper for tests
func random(min, max int) int {
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
func GetUser(uid string) (User, error) {
	userRef, err := fbClient.Ref("/users/" + uid)
	if err != nil {
		return User{}, err
	}

	var user User
	if err := userRef.Value(&user); err != nil {
		return User{}, err
	}
	if user.ID == "" {
		return User{}, errors.New("Failed to get user from Firebase")
	}
	return user, nil
}

// FetchUsers gets te Users from a list of uid's
func FetchUsers(uids []string) ([]User, error) {
	res := []User{}
	for _, uid := range uids {
		u, err := GetUser(uid)
		if err != nil {
			return []User{}, err
		}
		res = append(res, u)
	}
	return res, nil
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

// GetExpoPushToken Returns the push token for a specified user
func GetExpoPushToken(ID string) (string, error) {
	expoPushToken, err := fbClient.Ref("/users/" + ID + "/expoPushToken")
	if err != nil {
		return "", err
	}

	var token string
	if err := expoPushToken.Value(&token); err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New("expoPushToken not found")
	}

	return token, nil
}

// CreateThreadForUser Adds a threadlist entry for a user w/ID1
func CreateThreadForUser(ID1 string, threadID string, ID2 string) error {
	threadList, err := fbClient.Ref("/users/" + ID1 + "/threadList/" + threadID)
	if err != nil {
		return err
	}

	otherUserName, err := fbClient.Ref("/users/" + ID2 + "/profile/formattedName")
	if err != nil {
		return err
	}

	var name string
	if err := otherUserName.Value(&name); err != nil {
		return err
	}

	threadInfo := map[string]interface{}{
		"_id":  threadID,
		"name": name,
	}
	if err := threadList.Update(threadInfo); err != nil {
		return err
	}

	return nil
}

// AddThread Creates a messaging thread between two users
func AddThread(P1 string, P2 string) error {
	ID1, ID2 := "", ""

	if P1 == P2 {
		return errors.New("Cannot start a thread with only one user")
	} else if P1 < P2 {
		ID1, ID2 = P1, P2
	} else {
		ID1, ID2 = P2, P1
	}

	threadID := ID1 + separator + ID2

	if err := CreateThreadForUser(ID1, threadID, ID2); err != nil {
		return err
	}
	if err := CreateThreadForUser(ID2, threadID, ID1); err != nil {
		return err
	}

	thread, err := fbClient.Ref("/messages/" + threadID + "/0")
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

// CheckAuthorized checks if a user is authorized to make a request
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
