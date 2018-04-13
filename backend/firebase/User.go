package firebase

import (
	"errors"
	"fmt"
	"log"

	"meetover/backend/location"
	"meetover/backend/user"
)

// InitUser Updates access token if user exists or adds a new User as user in firebase
func InitUser(p user.Profile, aTokenResp user.ATokenResponse) (bool, error) {
	// Check if token exists
	userExists := false
	tokenRef, err := fbClient.Ref("/users/" + p.ID + "/accessToken")
	if err != nil {
		fmt.Println("Failed to save user profile to Firebase in InitUser()")
		return userExists, errors.New("Failed to save user profile: \n" + err.Error())
	}
	token := user.ATokenResponse{}
	if err := tokenRef.Value(&token); err != nil {
		log.Fatalf("Could not get user in Firebase: %v\n", err)
	}

	// Update token or create new user
	if token.AToken != "" {
		defer tokenRef.Set(aTokenResp)
		userExists = true
	} else {
		userRef, err := fbClient.Ref("/users/" + p.ID)
		if err != nil {
			fmt.Println("Failed to save user profile to Firebase in InitUser()")
			return false, errors.New("Failed to save user profile: \n" + err.Error())
		}

		u := user.User{
			ID:          p.ID,
			AccessToken: aTokenResp,
			Profile:     p,
		}
		defer userRef.Set(u)
	}
	return userExists, nil
}

// GetUser returns the user with uid in firebase
func GetUser(uid string) (user.User, error) {
	userRef, err := fbClient.Ref("/users/" + uid)
	if err != nil {
		return user.User{}, err
	}

	var u user.User
	if err := userRef.Value(&u); err != nil {
		return user.User{}, err
	}
	if u.ID == "" {
		return user.User{}, errors.New("Failed to get user from Firebase")
	}
	return u, nil
}

// GetProspectiveUsers Get the list of people for matching in the area
func GetProspectiveUsers(coords location.Geolocation, radius float64, lastUpdate int) ([]user.User, error) {
	// TODO:
	// Filter users by Geolocation and radius
	userMap := map[string]user.User{}
	users := []user.User{}
	realMatches := []user.User{}

	//create oldest date acceptable
	oldestStamp := location.MakeTimestamp((int64)(lastUpdate))
	
	userRef, err := fbClient.Ref("/users")
	if err != nil {
		fmt.Println("error getting userRef\n", err)
		return []user.User{}, err
	}

	
	if err := userRef.Value(&userMap); err != nil {
		fmt.Println("error getting users\n", err)
		return []user.User{}, err
	}

	for k := range userMap {
		users = append(users, userMap[k])
	}


	for _, pmatch := range users {
		if (pmatch.Location.TimeStamp > oldestStamp) && location.InRadius(coords, pmatch.Location, radius) {
				realMatches = append(realMatches, pmatch)
		}
	}


	return realMatches, nil
}

// GetFormattedName Gets the formatted name of the user with the supplied uid
func GetFormattedName(uid string) (string, error) {
	formattedName, err := fbClient.Ref("/users/" + uid + "/profile/formattedName")
	if err != nil {
		return "", err
	}
	var name string
	if err = formattedName.Value(&name); err != nil {
		return "", err
	}

	return name, nil
}

// FetchUsers gets te Users from a list of uid's
func FetchUsers(uids []string) ([]user.User, error) {
	res := []user.User{}
	for _, uid := range uids {
		u, err := GetUser(uid)
		if err != nil {
			return []user.User{}, err
		}
		res = append(res, u)
	}
	return res, nil
}
