package main

import (
	"fmt"
	"net/http"
	"os"
)

// our main function
func main() {

	router := NewRouter()
	InitializeFirebase()

	// test
	LoadTestUsers()

	// ML
	InitMLModel()
	port, deployMode := os.LookupEnv("PORT")
	if deployMode {
		fmt.Println(http.ListenAndServe(":"+port, router))
	} else {
		fmt.Println("running in debug mode")
		fmt.Println(http.ListenAndServe(":8080", router))
	}

}
