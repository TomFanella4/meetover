package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	// vectorTest()
	// rawFile := "rawTestUsers.json"
	// sinkFile := "FinalTestUsers.json"
	// GenTestUsers(rawFile, sinkFile)
	nt := genNonTechUsers(getRawUsers("rawTestUsers.json"))
	updateJSONFile(nt, "nt.json")
	// updateJSONFile(js, "testjs.json")

	// modelFile := "meetOver.model"
	// createModel(modelFile)
	// model := readModel(modelFile)

	// cachedUsers := getUsers(sinkFile)
	// n := len(cachedUsers)
	// start := random(0, n/2)
	// end := random((n/2)+1, n)
	// prospUsers := cachedUsers[start:end]
	// randomCaller := cachedUsers[random(0, n-1)]
	// order := distanceSort(randomCaller, prospUsers, model)
	// closest := order[0]
	// furthest := order[len(order)-1]
	// fmt.Println()
	// fmt.Println("caller: " + userToString(randomCaller))
	// fmt.Println()
	// fmt.Println()
	// fmt.Println("closest: " + userToString(getUser(closest, cachedUsers)))
	// fmt.Println()
	// fmt.Println()
	// fmt.Println("furthest: " + userToString(getUser(furthest, cachedUsers)))

}
