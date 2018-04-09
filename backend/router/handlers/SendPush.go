package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SendPush sends sends a push notification for a verified user
func SendPush(w http.ResponseWriter, r *http.Request) {
	if CheckAuthorized(w, r) {
		var pushNotification PushNotification
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &pushNotification); err != nil {
			fmt.Println("Unable to send push notification")
			respondWithError(w, FailedSendPush, "Could not send push notification")
			return
		}
		err := SendPushNotification(&pushNotification)
		if err != nil {
			respondWithError(w, FailedSendPush, "Could not send push notification")
			fmt.Println(err.Error())
			return
		}
		resp := ServerResponse{Success, "Success", true}
		json.NewEncoder(w).Encode(resp)
	}
}
