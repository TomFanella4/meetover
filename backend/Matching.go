package main

import "errors"

//TODO:
// - take user ID and find k cachedUsers that are:
//		- in the same location
//		- will be closest in similarity metric that is picked

// GetMatches returns k user profiles that the UserID will want to meet

func GetMatches(UserID string, neighbors []User) ([]User, error) {
	callingUser := GetUser(userID)
	emptyUSer := User{}
	if emptyUSer == emptyUSer {
		return []User{}, errors.New("Unable to fetch calling user")
	}


	return cachedUsers, nil
}
