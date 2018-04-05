package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"meetover/backend/firebase"

	"github.com/gorilla/mux"
)

// InitiateMeetover called to begin the meetover appointment betwen two users
func InitiateMeetover(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		initiatorID := r.Header.Get("Identity")
		params := mux.Vars(r)
		requestedID := params["otherID"]

		if err := firebase.AddThread(initiatorID, requestedID); err != nil {
			respondWithError(w, FailedDBCall, "Could not create chat thread")
			fmt.Println("Failed to create the thread " + initiatorID + ", " + requestedID)
			return
		}

		// Send a push notification to the requested user
		name, err := firebase.GetFormattedName(initiatorID)
		if err != nil {
			respondWithError(w, FailedDBCall, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}

		title := "New MeetOver Request"
		body := name + " would like to MeetOver"
		pushNotification := PushNotification{
			ID:    requestedID,
			Title: title,
			Body:  body,
		}
		err = SendPushNotification(&pushNotification)
		if err != nil {
			respondWithError(w, FailedSendPush, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}
		resp := ServerResponse{Success, "Success", true}
		json.NewEncoder(w).Encode(resp)
	}
}
