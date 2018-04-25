package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"meetover/backend/firebase"
	"meetover/backend/matching"
	"meetover/backend/user"
)

func getTestUsers(rawFile string) []user.User {
	raw, err := ioutil.ReadFile(rawFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var users []user.User
	json.Unmarshal(raw, &users)
	return users
}

func testMatch() {
	p := fmt.Println
	userFile := "./data/FinalTestUsers.json"
	var users []user.User
	users = getTestUsers(userFile)
	for i := 0; i < 10; i++ {
		ri := rand.Intn(len(users))
		ru := users[ri]
		p("Requesting user: ")
		p(ru.Profile.Summary)
		p("matches: ")
		mt, err := matching.GetMatches(ru.Profile.ID, ru, users)
		if err != nil {
			p(err.Error())
		}
		for j := 0; j < 3; j++ {
			fmt.Printf("[%d, %f] %s\n", j, mt[j].Dist, mt[j].Usr.Summary)
		}
	}
}

// our main function
func main() {

	// router := router.NewRouter()

	// Initialiaze database, chat, and static storage
	firebase.InitializeFirebase()
	firebase.InitializeFiles()

	// ML
	matching.InitMLModel(matching.WordModelContextWindow, matching.WordModelDimension)
	rand.Seed(time.Now().Unix())

	testMatch()
	// port, deployMode := os.LookupEnv("PORT")
	// if deployMode {
	// 	fmt.Println(http.ListenAndServe(":"+port, router))
	// } else {
	// 	fmt.Println("running in debug mode")
	// 	fmt.Println(http.ListenAndServe(":8080", router))
	// }

}
