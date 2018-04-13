package firebase

import (
	"errors"
	"time"
)

func createThreadForUser(ID1 string, threadID string, ID2 string, origin string) error {
	threadList, err := fbClient.Ref("/users/" + ID1 + "/threadList/" + threadID)
	if err != nil {
		return err
	}

	threadInfo := map[string]interface{}{
		"_id":    threadID,
		"userID": ID2,
		"origin": origin,
		"status": "pending",
	}
	if err := threadList.Update(threadInfo); err != nil {
		return err
	}

	return nil
}

// AddThread Creates a messaging thread between two users
func AddThread(P1 string, P2 string, initialMessage string) error {
	ID1, ID2 := "", ""
	origin1, origin2 := "", ""

	if P1 == P2 {
		return errors.New("Cannot start a thread with only one user")
	} else if P1 < P2 {
		ID1, ID2 = P1, P2
		origin1, origin2 = "sender", "receiver"
	} else {
		ID1, ID2 = P2, P1
		origin1, origin2 = "receiver", "sender"
	}

	threadID := ID1 + separator + ID2

	if err := createThreadForUser(ID1, threadID, ID2, origin1); err != nil {
		return err
	}
	if err := createThreadForUser(ID2, threadID, ID1, origin2); err != nil {
		return err
	}

	initiatorUser, err := GetUser(P1)
	if err != nil {
		return err
	}

	initiator := map[string]interface{}{
		"_id":    initiatorUser.Profile.ID,
		"avatar": initiatorUser.Profile.PictureURL,
		"name":   initiatorUser.Profile.FormattedName,
	}

	startMessageMap := map[string]interface{}{
		"_id":       0,
		"text":      "Start of the chat",
		"createdAt": time.Now().UTC().Format(time.RFC3339),
		"system":    true,
	}
	initialMessageMap := map[string]interface{}{
		"_id":       1,
		"text":      initialMessage,
		"createdAt": time.Now().UTC().Format(time.RFC3339),
		"user":      initiator,
	}

	messagesMap := map[string]interface{}{
		"0": startMessageMap,
		"1": initialMessageMap,
	}

	threadRef, err := fbClient.Ref("/messages/" + threadID)
	if err != nil {
		return err
	}

	if err := threadRef.Update(messagesMap); err != nil {
		return err
	}

	return nil
}

// SetThreadStatus Sets the status of a thread to the specified value
func SetThreadStatus(userID string, threadID string, status string) error {
	threadRef, err := fbClient.Ref("/users/" + userID + "/threadList/" + threadID)
	if err != nil {
		return err
	}
	statusMap := map[string]string{"status": status}
	threadRef.Update(statusMap)
	return nil
}
