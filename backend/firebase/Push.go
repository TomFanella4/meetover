package firebase

import "errors"

// GetExpoPushToken Returns the push token for a specified user
func GetExpoPushToken(ID string) (string, error) {
	expoPushToken, err := fbClient.Ref("/users/" + ID + "/expoPushToken")
	if err != nil {
		return "", err
	}

	var token string
	if err := expoPushToken.Value(&token); err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New("expoPushToken not found")
	}

	return token, nil
}
