package matching

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"meetover/backend/firebase"
	"meetover/backend/location"
	"meetover/backend/user"
)

// Relative path to data directory. Relative from location of backend binary
var DataDir = "./data/"

// MatchValue represents each perspective user in their distance from the caller
type MatchValue struct {
	Usr  user.Profile         `json:"profile"`
	Dist float64              `json:"distance"`
	Loc  location.Geolocation `json:"location"`
}

type byDistance []MatchValue

func (b byDistance) Len() int {
	return len(b)
}

func (b byDistance) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byDistance) Less(i, j int) bool {
	return b[i].Dist < b[j].Dist
}

// GetMatches returns an ordered list of user uid's from closest to furthest to the caller
func GetMatches(UserID string, neighbors []user.User) ([]MatchValue, error) {
	callingUser, err := firebase.GetUser(UserID)
	if err != nil {
		return nil, errors.New("Unable to fetch calling user with uid: " + UserID)
	}
	order := GetOrder(callingUser, neighbors, WordModel)
	return order, nil
}

// GetOrder - preprocess for sortMap
func GetOrder(caller user.User, prospUsers []user.User, model map[string][]float64) []MatchValue {
	callerStr := userToString(caller)
	callerVec := parToVector(callerStr, model)
	prospUsers = removeCaller(caller, prospUsers)
	matches := []MatchValue{}
	start := time.Now()
	for _, pu := range prospUsers {
		prospStr := userToString(pu)
		prospVec := parToVector(prospStr, model)
		distance := nestedDistance(callerVec, prospVec)
		matches = append(matches, MatchValue{pu.Profile, distance, pu.Location})
	}
	sort.Sort(byDistance(matches))
	elapsed := time.Since(start)
	fmt.Println("Destance Calculation took: " + elapsed.String())
	return matches
}

// removeCaller takes the calling user out of prospective match list
func removeCaller(caller user.User, prospUsers []user.User) []user.User {
	s := -1
	for i, u := range prospUsers {
		if u.ID == caller.ID {
			s = i
		}
	}
	if s < 0 {
		return prospUsers
	}
	return append(prospUsers[:s], prospUsers[s+1:]...)

}

// random - helper for tests
func random(min, max int) int {
	return rand.Intn(max-min) + min
}
