package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"meetover/backend/firebase"
	"meetover/backend/matching"
	"meetover/backend/router"
)

// our main function
func main() {

	router := router.NewRouter()

	// Initialiaze database, chat, and static storage
	firebase.InitializeFirebase()
	firebase.InitializeFiles()

	// ML
	matching.InitMLModel(matching.WordModelContextWindow, matching.WordModelDimension)
	rand.Seed(time.Now().Unix())

	port, deployMode := os.LookupEnv("PORT")
	if deployMode {
		fmt.Println(http.ListenAndServe(":"+port, router))
	} else {
		fmt.Println("running in debug mode")
		fmt.Println(http.ListenAndServe(":8080", router))
	}

}
