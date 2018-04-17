package main

import (
	"fmt"
	"strings"

	"meetover/backend/matching"
)

// our main function
func main() {

	t := "I'm going to present three projects in rapid fire. I don't have much time to do it. And I want to reinforce three ideas with that rapid-fire presentation.The first is what I like to call a hyper-rational process. It's a process that takes rationality almost to an absurd level, and it transcends all the baggage that normally comes with what people would call, sort of a rational conclusion to something. And it concludes in something that you see here, that you actually wouldn't expect as being the result of rationality.The second â€” the second is that this process does not have a signature. There is no authorship. Architects are obsessed with authorship. This is something that has editing and it has teams, but in fact, we no longer see within this process, the traditional master architect creating a sketch that his minions carry out.And the third is that it challenges â€” and this is, in the length of this, very hard to support why, connect all these things â€” but it challenges the high modernist notion of flexibility. High modernists said we will create sort of singular spaces that are generic, almost anything can happen within them. I call it sort of "
	ta := strings.Split(t, " ")
	fmt.Println("Before:\n" + t)
	fmt.Println("After:\n" + matching.StripStopWords(ta))

	// router := router.NewRouter()

	// // Initialiaze database, chat, and static storage
	// firebase.InitializeFirebase()
	// firebase.InitializeFiles()

	// // ML
	// matching.InitMLModel(matching.WordModelContextWindow, matching.WordModelDimension)
	// rand.Seed(time.Now().Unix())

	// port, deployMode := os.LookupEnv("PORT")
	// if deployMode {
	// 	fmt.Println(http.ListenAndServe(":"+port, router))
	// } else {
	// 	fmt.Println("running in debug mode")
	// 	fmt.Println(http.ListenAndServe(":8080", router))
	// }

}
