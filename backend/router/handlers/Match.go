package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"meetover/backend/firebase"
	"meetover/backend/matching"
)

// Match will set a flag to notify the system the user is matched
func Match(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		// requesters uid
		params := mux.Vars(r)
		userID := params["uid"]
		// get requester's coords
		user, err := firebase.GetUser(userID)
		if err != nil {
			fmt.Println(err)
			return // Unable to get location
		}
		radius := 1.0     //1 km
		lastUpdate := 2 // 2hrs
		PMatchList, err := firebase.GetProspectiveUsers(user.Location, radius, lastUpdate)
		if err != nil {
			fmt.Println(err)
			return // unable to get anyone from db
		}
		MatchList, err := matching.GetMatches(userID, PMatchList)
		if err != nil {
			fmt.Println(err)
			return // unable to get anyone to match
		}
		mr := MatchResponse{}
		mr.Matches = MatchList
		json.NewEncoder(w).Encode(mr)
	}
}
