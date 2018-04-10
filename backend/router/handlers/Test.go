package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meetover/backend/location"
	"meetover/backend/firebase"
	"net/http"

	"github.com/gorilla/mux"
)

// Test returns a sample LinkedIn Profile JSON object
func Test(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tt := params["testType"]
	if tt == "distance" {
		var testLoc location.Geolocation
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &testLoc); err != nil {
			bodyString := string(bodyBytes)
			fmt.Println(bodyString)
			return // no match was returned
		} //40.4259° N, 86.9081° W
		res := location.InRadius(location.Geolocation{Longitude: -86.9081, Latitude: 40.4259}, testLoc, 20)
		rj := make(map[string]bool)
		rj["res"] = res
		respondWithJSON(w, 200, rj)
		return
	} else if tt == "query" {
		//NAECC Coords

		hour := (int64)(1)
		testCoords := location.Geolocation{
			Latitude: 40.4291,
			Longitude: -86.91686,
			TimeStamp: location.MakeTimestamp(hour),
			Accuracy: 100,
		}

		//Mock Radius & refresh
		radi := 20.0
		refre := 2

		firebase.GetProspectiveUsers(testCoords, radi , refre)
		/*testCoords.timeStamp = makeTimestamp(1)
		GetProspectiveUsers(testCoords, radi , refre)*/
	}
	resp := ServerResponse{Success, "Success", true}
	json.NewEncoder(w).Encode(resp)
}
