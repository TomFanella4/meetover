package firebase

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/zabawaba99/firego"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var fbApp *firebase.App
var fbClient *firego.Firebase
var fbBucket *storage.BucketHandle

const separator = "|"

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
	storageConfig := &firebase.Config{
		StorageBucket: "meetoverdb.appspot.com",
	}
	opt := option.WithCredentialsFile("./firebase-config.json")
	app, err := firebase.NewApp(context.Background(), storageConfig, opt)
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

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}

	fbApp = app
	fbClient = firego.New("https://meetoverdb.firebaseio.com/", conf.Client(oauth2.NoContext))
	fbBucket = bucket
}

// InitializeFiles Initializes the files
func InitializeFiles() error {
	fmt.Println("[+] Initializing files...")
	modelFile := "./data/meetOver.model"
	// Check fs for model
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		// Check Firebase for model
		fmt.Println("[-] Model file not found. Downloading model from Firebase...")
		ctx := context.Background()
		modelObject := fbBucket.Object("meetOver.model")
		modelReader, err := modelObject.NewReader(ctx)

		if err != nil {
			// Check fs for corpus
			fmt.Println("[-] Model not found in Firebase. Checking for Corpus file...")
			corpusFile := "./data/corpus.dat"

			if _, err = os.Stat(corpusFile); os.IsNotExist(err) {
				// Check Firebase for corpus
				fmt.Println("[-] Corpus file not found. Downloading corpus from Firebase...")
				corpusObject := fbBucket.Object("corpus.dat")
				corpusReader, err1 := corpusObject.NewReader(ctx)
				if err1 != nil {
					// Corpus could not be found
					fmt.Println("[-] Corpus file not found in Firebase.")
					return err1
				}

				defer corpusReader.Close()
				content, err1 := ioutil.ReadAll(corpusReader)
				if err1 != nil {
					return err1
				}

				err = ioutil.WriteFile(corpusFile, content, 0644)
				if err != nil {
					return err
				}

				fmt.Println("[+] Loaded corpus file from Firebase.")
				return nil
			}

			fmt.Println("[+] Found corpus file.")
			return nil
		}

		// Download model from Firebase
		defer modelReader.Close()
		content, err := ioutil.ReadAll(modelReader)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(modelFile, content, 0644)
		if err != nil {
			return err
		}
		fmt.Println("[+] Loaded model file from Firebase.")
		return nil
	}
	fmt.Println("[+] Found model file.")
	return nil
}
