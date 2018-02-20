package main

// TODO:
// curl commands for firebase
// methods to store, publish, modify, delete
// JSON schema for user profile
// IM schema and maintenance

import (
  "context"
  "fmt"
  "os"
  "io/ioutil"

  firebase "firebase.google.com/go"
  // "firebase.google.com/go/auth"

  // "google.golang.org/api/iterator"
  "google.golang.org/api/option"
)

var firebaseApp *firebase.App

func InitializeFirebase() {
  // This is a Firebase workaround. The Firebase go library MUST read
  // certificate credentials from a file. Also, Heroku doesn't have static
  // storage. So, we must create the credential file dynamically from an
  // environment variable containing the json
  config := []byte(os.Getenv("FIREBASE_CONFIG"))
  err := ioutil.WriteFile("./firebase-config.json", config, 0644)
  if err != nil {
    fmt.Printf("Could not write config file")
  }

  opt := option.WithCredentialsFile("./firebase-config.json")
  fbApp, err := firebase.NewApp(context.Background(), nil, opt)
  if err != nil {
    fmt.Printf("error initializing app: %v\n", err)
  }

  firebaseApp = fbApp
}
