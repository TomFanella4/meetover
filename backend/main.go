package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)


type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
    Address   *Address `json:"address,omitempty"`
}
type Address struct {
    City  string `json:"city,omitempty"`
    State string `json:"state,omitempty"`
	Area string `json:"area,omitempty"`
}

type ServerResponse struct {
	Code  int `json:"id"`
	Message string `json:"message"`
	Success  bool `json:"success"`
}
	
var people []Person


// our main function
func main() {

	// Test Data
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "Chicago", State: "IL", Area: "ORD"}})
people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "Chicago", State: "IL", Area: "ORD"}})
people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday", Address: &Address{City: "New York", State: "NY", Area: "JFK" }})


    router :=mux.NewRouter()
	router.HandleFunc("/location/{coords}", GetAddress).Methods("GET")
router.HandleFunc("/people/{id}", GetPeople).Methods("GET")
router.HandleFunc("/login/{cred}", VerifyUser).Methods("POST")
router.HandleFunc("/match/{ouser}", Match).Methods("POST")

    log.Fatal(http.ListenAndServe(":8080", router))
}


func GetAddress(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	//  location = getFromDataSet( params["coords"].split(",") )

	location := Address{ City: "Chicago", State: "IL", Area: "ORD"}
	json.NewEncoder(w).Encode(location)
}
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	params :=  mux.Vars(r)
	// var person Person
	//  _ = json.NewDecoder(r.Body).Decode(&person)
	//  person.ID = params["id"]
	//  people = append(people, person)
	cred :=  params["cred"]
	msg := cred + " Login Successful" 
	sr := ServerResponse{ Code: 200, Message: msg,  Success: true}
	json.NewEncoder(w).Encode(sr)

}
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
