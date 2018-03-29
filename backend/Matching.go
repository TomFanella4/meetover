package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	// "github.com/ynqa/word-embedding/builder"
)

// MatchResponse returned to the UI when /match is hit
type MatchResponse struct {
	Matches []MatchValue `json:"matches"`
}

// 5abc5152c2d9048b32bfc917

// MatchValue represents each porspecive user in their distance from the caller
type MatchValue struct {
	U User    `json:"user"`
	D float64 `json:"distance"`
}

// WordModel is the vector representation of the words in the corpus file
var WordModel map[string][]float64

// ParVecParam - number of words considered for similarity
var ParVecParam = 50

// GetMatches returns an ordered list of user uid's from closest to furthest to the caller
func GetMatches(UserID string, neighbors []User) (MatchResponse, error) {
	callingUser, err := GetUser(UserID)
	if err != nil {
		return MatchResponse{}, errors.New("Unable to fetch calling user")
	}
	fmt.Println("Got user")
	fmt.Println(callingUser)
	order := GetOrder(callingUser, neighbors, WordModel)
	return order, nil
}

// GetOrder - preprocess for sortMap
func GetOrder(caller User, prospUsers []User, model map[string][]float64) MatchResponse {
	callerStr := userToString(caller)
	callerVec := parToVector(callerStr, model)
	prospUsers = removeCaller(caller, prospUsers)
	var mr MatchResponse
	mr.Matches = []MatchValue{}
	for _, pu := range prospUsers {
		prospStr := userToString(pu)
		prospVec := parToVector(prospStr, model)
		fmt.Println("Running Nested Distance")
		distance := nestedDistance(callerVec, prospVec)
		fmt.Println("Finished Nested Distance")
		mr.Matches = append(mr.Matches, MatchValue{pu, distance})
	}
	return mr
}

// sortMap - returns uid's with shortest distance first
func sortMap(m map[string]float64) []string {
	reverseMap := map[float64]string{}
	distances := []float64{}
	for uid, d := range m {
		reverseMap[d] = uid
		distances = append(distances, d)
	}
	sort.Float64s(distances)
	res := []string{}
	for _, d := range distances {
		res = append(res, reverseMap[d])
	}
	fmt.Println(res)
	return res
}

// nestedDistance - distance metric between par vectors
func nestedDistance(src []float64, dst []float64) float64 {
	d := 0.0
	for _, si := range src {
		for _, di := range dst {
			d += math.Abs(si - di)
		}
	}
	return d
}

// removeCaller takes the calling user out of prospective match list
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

// parToVector converts a string representation of user to numeric vector using
// the given word embeddings model
func parToVector(userStr string, model map[string][]float64) []float64 {
	res := []float64{}
	par := strings.Split(userStr, " ")
	n := ParVecParam
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

// userToString converts the user object to a paragraph for vector translation
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

// InitMLModel check if model has been created or creates it
func InitMLModel() {
	modelFile := "./ml/meetOver.model"
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		// corpusFile := "./ml/corpus.dat"
		// createModel(modelFile, corpusFile)
		fmt.Println("Model does not exist")
	}
	WordModel = readModel(modelFile)
}

// createModel uses the word2vec algo to create word embeddings
// func createModel(destinationFileName string, corpusFile string) {
// 	if _, err := os.Stat(corpusFile); os.IsNotExist(err) {
// 		fmt.Println("[-] Corpus file not found. No model created")
// 	}
// 	b := builder.NewWord2VecBuilder()
// 	b.SetDimension(5).
// 		SetWindow(20).
// 		SetModel("cbow").
// 		SetOptimizer("ns").
// 		SetNegativeSampleSize(5).
// 		SetVerbose()
// 	m, err := b.Build()
// 	if err != nil {
// 		fmt.Println("[-] Unable to build word2vec neural net")
// 	}
// 	inputFile1, _ := os.Open(corpusFile)
// 	f1, err := m.Preprocess(inputFile1)
// 	if err != nil {
// 		fmt.Println("Failed to Preprocess.")
// 	}
// 	// Start to Train.
// 	m.Train(f1)
// 	f1.Close()
// 	// Save word vectors to a text file.
// 	m.Save(destinationFileName)
// }

// readModel converts the generated model to an in-memory object
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
