package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meetover/backend/firebase"
	"net/http"

	"github.com/gorilla/mux"
)

// ProcessDecision Updates the status of the specified user pair
func ProcessDecision(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		requestedID := r.Header.Get("Identity")
		params := mux.Vars(r)
		initiatorID := params["otherID"]

		var decision MeetOverDecisionBody
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &decision); err != nil {
			fmt.Println(err.Error())
			respondWithError(w, FailedMeetOver, "Failed to set MeetOver decision")
			return
		}

		if err := firebase.SetThreadStatus(requestedID, decision.ThreadID, decision.Status); err != nil {
			fmt.Println(err.Error())
			respondWithError(w, FailedMeetOver, "Failed to set MeetOver decision")
			return
		}
		if err := firebase.SetThreadStatus(initiatorID, decision.ThreadID, decision.Status); err != nil {
			fmt.Println(err.Error())
			respondWithError(w, FailedMeetOver, "Failed to set MeetOver decision")
			return
		}

		// Send a push notification to the initiator
		name, err := firebase.GetFormattedName(requestedID)
		if err != nil {
			respondWithError(w, FailedDBCall, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}

		pushNotification := PushNotification{
			ID:    initiatorID,
			Title: name + " " + decision.Status + " request",
			Body:  "",
		}
		if err := SendPushNotification(&pushNotification); err != nil {
			fmt.Println("Unable to send push notification")
			respondWithError(w, FailedSendPush, "Could not send push notification")
			return
		}
		resp := ServerResponse{Success, "Success", true}
		json.NewEncoder(w).Encode(resp)
	}
}
