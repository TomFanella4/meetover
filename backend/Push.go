package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

// PushNotification struct
type PushNotification struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SendPushNotification sends a push notification to the speicifed recipient
func SendPushNotification(pushNotification *PushNotification) error {
	endpoint := "https://exp.host/--/api/v2/push/send"
	expoPushToken, err := GetExpoPushToken(pushNotification.ID)
	if err != nil {
		return errors.New("Failed to send push notification to " + pushNotification.ID)
	}
	data := fmt.Sprintf(`{"to":"%s", "title": "%s", "body": "%s"}`,
		expoPushToken, pushNotification.Title, pushNotification.Body)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return errors.New("Unable to create Expo Push Token API call")
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return errors.New("Unable to call Expo Push Token API call")
	}

	fmt.Println("[+] Push notification sent to " + pushNotification.ID)
	return nil
}
