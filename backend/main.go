package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// ClientID - temp
var ClientID string

// Person is user on MeetOver
type Person struct {
	ID          string   `json:"id,omitempty"`
	Firstname   string   `json:"firstname,omitempty"`
	Lastname    string   `json:"lastname,omitempty"`
	Address     *Address `json:"address,omitempty"`
	AccessToken string   `json:"accesstoken"`
}

// Address is a our location metric
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	Area  string `json:"area,omitempty"`
}

// ServerResponse is what we give to our clients
type ServerResponse struct {
	Code    int    `json:"id"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

var people []Person

// our main function
func main() {

	// Test Data
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "Chicago", State: "IL", Area: "ORD"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "Chicago", State: "IL", Area: "ORD"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday", Address: &Address{City: "New York", State: "NY", Area: "JFK"}})

	router := mux.NewRouter()
	router.HandleFunc("/location/{coords}", GetAddress).Methods("GET")
	router.HandleFunc("/people/{id}", GetPeople).Methods("GET")
	router.HandleFunc("/login/{code}", VerifyUser).Methods("POST")
	router.HandleFunc("/userinfo/{code}", GetUserProfile).Methods("POST")
	router.HandleFunc("/match/{ouser}", Match).Methods("POST")

	port, deployMode := os.LookupEnv("PORT")
	if !deployMode {
		log.Fatal(http.ListenAndServe(":"+port, router))
	} else {
		log.Fatal(http.ListenAndServe(":8080", router))
	}
}

// GetUserProfile will give back a json object of user's LinkedIn Profile
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accessCode := params["code"]
	profile := GetLiProfile(accessCode)
	json.NewEncoder(w).Encode(profile)
}

// GetAddress will give back a json object based on coordinates
func GetAddress(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	//  location = getFromDataSet( params["coords"].split(",") )
	location := Address{City: "Chicago", State: "IL", Area: "ORD"}
	json.NewEncoder(w).Encode(location)
}

// GetPeople returns all users
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// VerifyUser will get a code object to obtain an access token
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	TempUserCode := params["code"]
	AToken := ExchangeToken(TempUserCode, "https://meetover407.herokuapp.com")
	profile := GetLiProfile(AToken)
	json.NewEncoder(w).Encode(profile)
}

// Match will set a flag to notify the system the suer is matched
func Match(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
