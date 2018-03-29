package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	// . "github.com/bugra/kmeans"
	"github.com/ynqa/word-embedding/builder"
	// "gonum.org/v1/gonum/mat"
	// "meetover/backend/main"
)

func main() {
	rawFile := "rawTestUsers.json"
	sinkFile := "MLTestUsers.json"
	generateTestUsers(rawFile, sinkFile)

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

func getUser(uid string, users []User) User {
	for _, u := range users {
		if u.ID == uid {
			return u
		}
	}
	return User{}
}
func distanceSort(caller User, prospUsers []User, model map[string][]float64) []string {
	distances := make(map[string]float64)
	callerStr := userToString(caller)
	callerVec := parToVector(callerStr, model)
	prospUsers = removeCaller(caller, prospUsers)
	for _, pu := range prospUsers {
		prospStr := userToString(pu)
		prospVec := parToVector(prospStr, model)
		distances[pu.ID] = nestedDistance(callerVec, prospVec)
	}
	return sortMap(distances)
}

// returns uid's with shortest distance first
func sortMap(m map[string]float64) []string {
	reverseMap := map[float64]string{}
	distances := []float64{}
	for uid, d := range m {
		reverseMap[d] = uid
		distances = append(distances, d)
	}
	sort.Float64s(distances)
	printMap(m)
	res := []string{}
	for _, d := range distances {
		res = append(res, reverseMap[d])
	}
	fmt.Println(res)
	return res
}
func printMap(m map[string]float64) {
	for k, v := range m {
		fmt.Print(k)
		fmt.Printf(" : %f\n", v)
	}
}
func nestedDistance(src []float64, dst []float64) float64 {
	d := 0.0
	for _, si := range src {
		for _, di := range dst {
			d += math.Abs(si - di)
		}
	}
	return d
}

func removeCaller(caller User, prospUsers []User) []User {
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
func parToVector(userStr string, model map[string][]float64) []float64 {
	res := []float64{}
	par := strings.Split(userStr, " ")
	n := 200
	i := 0
	for i < n {
		l := len(par)
		randomWord := par[random(0, l)]
		if val, found := model[randomWord]; found {
			res = append(val, res...)
			i++
		}
	}
	return res
}

// random - helper for tests
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func userToString(u User) string {
	var res string
	res = u.Profile.Greeting + " " + u.Profile.Headline + " " + u.Profile.Summary + " " + u.Profile.Industry
	for _, pos := range u.Profile.Positions.Values {
		res += pos.Company.Industry + " "
		res += pos.Company.Name + " "
		res += pos.Summary + " "
		res += pos.Title
	}
	return res
}

func readModel(modelFile string) map[string][]float64 {
	content, err := ioutil.ReadFile(modelFile)
	if err != nil {
		log.Fatal(err)
	}
	model := make(map[string][]float64)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		vector := strings.Split(line, " ")
		word := vector[0]
		n := len(vector)
		if n < 2 {
			vector = vector[1:]
		} else {
			vector = vector[1 : len(vector)-1]
		}
		floatVector := make([]float64, len(vector))
		for jj := range vector {
			floatVector[jj], err = strconv.ParseFloat(vector[jj], 64)
		}
		model[word] = floatVector
	}
	return model
}
func createModel(destinationFileName string) {
	b := builder.NewWord2VecBuilder()

	b.SetDimension(5).
		SetWindow(20).
		SetModel("cbow").
		SetOptimizer("ns").
		SetNegativeSampleSize(5).
		SetVerbose()

	m, err := b.Build()

	if err != nil {
		// Failed to build word2vec.
	}

	inputFile1, _ := os.Open("combined.dat")
	// inputFile2, _ := os.Open("transcripts_clean.csv")

	f1, err := m.Preprocess(inputFile1)
	// f2, err := m.Preprocess(inputFile2)

	if err != nil {
		fmt.Println("Failed to Preprocess.")
	}

	// Start to Train.
	m.Train(f1)
	// m.Train(f2)
	f1.Close()
	// f2.Close()

	// Save word vectors to a text file.
	m.Save(destinationFileName)

}
