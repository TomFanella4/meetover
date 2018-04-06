package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// our main function
func main() {

	router := NewRouter()

	// database & chat
	InitializeFirebase()

	// demo data
	// LoadTestUsers()

	// ML
	InitMLModel(WordModelContextWindow, WordModelDimension)
	rand.Seed(time.Now().Unix())

	port, deployMode := os.LookupEnv("PORT")
	if deployMode {
		fmt.Println(http.ListenAndServe(":"+port, router))
	} else {
		fmt.Println("running in debug mode")
		fmt.Println(http.ListenAndServe(":8080", router))
	}

}
