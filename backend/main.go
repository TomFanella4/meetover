package main

import (
	"fmt"
	"net/http"
	"os"
)

// our main function
func main() {

	// Test Data
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "Chicago", State: "IL", Area: "ORD"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "Chicago", State: "IL", Area: "ORD"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday", Address: &Address{City: "New York", State: "NY", Area: "JFK"}})

	router := NewRouter()

	InitializeFirebase()

	port, deployMode := os.LookupEnv("PORT")
	if deployMode {
		fmt.Println(http.ListenAndServe(":"+port, router))
	} else {
		fmt.Println("running in dubug mode")
		fmt.Println(http.ListenAndServe(":8080", router))
	}
}
