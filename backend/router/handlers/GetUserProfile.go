package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meetover/backend/firebase"
	"net/http"
)

// GetUserProfiles will give back a json object of user's Profile
func GetUserProfiles(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		var ids []string
		profiles := make(map[string]UserProfileResponse)
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &ids); err != nil {
			fmt.Println(err.Error())
			respondWithError(w, FailedProfileFetch, "Unable to get user profiles")
			return
		}

		for _, id := range ids {
			user, err := firebase.GetUser(id)
			if err != nil {
				fmt.Println(err.Error())
				respondWithError(w, FailedProfileFetch, "Failed to get User")
				return
			}
			var profile UserProfileResponse
			profile.Profile = user.Profile
			profile.Location = user.Location
			profiles[id] = profile
		}
		json.NewEncoder(w).Encode(profiles)
	}
}
