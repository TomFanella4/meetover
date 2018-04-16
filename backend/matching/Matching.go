package matching

import (
	"errors"
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

	"meetover/backend/firebase"
	"meetover/backend/location"
	"meetover/backend/user"

	"github.com/mpace965/word-embedding/builder"
	"gonum.org/v1/gonum/mat"
)

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

// WordModel is the vector representation of the words in the corpus file
var WordModel map[string][]float64

// WordModelContextWindow -
var WordModelContextWindow = 20

// WordModelDimension -
var WordModelDimension = 8

// WordModelRandomParam - number of words considered for similarity
var WordModelRandomParam = 50

// Relative path to ml directory. Relative from location of backend binary
var mlDirectory = "./data/"

// GetMatches returns an ordered list of user uid's from closest to furthest to the caller
func GetMatches(UserID string, neighbors []user.User) ([]MatchValue, error) {
	callingUser, err := firebase.GetUser(UserID)
	if err != nil {
		return nil, errors.New("Unable to fetch calling user")
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

// nestedDistance - distance metric between par vectors
func nestedDistance(src []*mat.VecDense, dst []*mat.VecDense) float64 {
	d := 0.0
	for _, si := range src {
		for _, di := range dst {
			temp := mat.NewVecDense(WordModelDimension, nil)
			temp.SubVec(si, di)
			d += flattenVector(WordModelDimension, temp)
		}
	}
	return d
}
func flattenVector(rows int, vec mat.Matrix) float64 {
	res := 0.0
	for i := 0; i < rows; i++ {
		res += math.Abs(vec.At(i, 0))
	}
	return res
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

// stripStopWords -
func stripStopWords(str string) string {
	return ""
}

// parToVector converts a string representation of user to numeric vector using
// the given word embeddings model
func parToVector(userStr string, model map[string][]float64) []*mat.VecDense {
	res := []*mat.VecDense{}
	par := strings.Split(userStr, " ")
	n := WordModelRandomParam
	i := 0
	for i < n {
		l := len(par)
		randomWord := par[random(0, l)]
		randomWord = strings.TrimSpace(strings.ToLower(randomWord))
		if val, found := model[randomWord]; found {
			if len(val) == WordModelDimension {
				vec := mat.NewVecDense(WordModelDimension, val)
				res = append(res, vec)
				i++
			}
		}
	}
	return res
}

// userToString converts the user object to a paragraph for vector translation
func userToString(u user.User) string {
	var res string
	res = u.MatchStatus.Greeting + " " + u.Profile.Headline + " " + u.Profile.Summary + " " + u.Profile.Industry
	for _, pos := range u.Profile.Positions.Values {
		res += pos.Company.Industry + " "
		res += pos.Company.Name + " "
		res += pos.Summary + " "
		res += pos.Title
	}
	return res
}

// InitMLModel check if model has been created or creates it
func InitMLModel(windowSize int, wordDimensions int) {
	modelFile := mlDirectory + "meetOver.model"
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		fmt.Println("Model does not exist. Creating Model")
		corpusFile := mlDirectory + "corpus.dat"
		createModel(modelFile, corpusFile, windowSize, wordDimensions)
	}
	WordModel = readModel(modelFile)
}

// createModel uses the word2vec algo to create word embeddings
func createModel(destinationFileName string, corpusFile string, windowSize int, wordDimensions int) {
	if _, err := os.Stat(corpusFile); os.IsNotExist(err) {
		fmt.Println("[-] Corpus file not found. No model created")
		return
	}
	b := builder.NewWord2VecBuilder()
	b.SetDimension(wordDimensions).
		SetWindow(windowSize).
		SetModel("skip-gram").
		SetOptimizer("ns").
		SetNegativeSampleSize(15).
		SetVerbose()
	m, err := b.Build()
	if err != nil {
		fmt.Println("[-] Unable to build word2vec neural net")
	}
	inputFile1, _ := os.Open(corpusFile)
	f1, err := m.Preprocess(inputFile1)
	if err != nil {
		fmt.Println("Failed to Preprocess.")
	}
	// Start to Train.
	m.Train(f1)
	f1.Close()
	// Save word vectors to a text file.
	m.Save(destinationFileName)
}

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
			if err != nil {
				fmt.Println(err)
				return make(map[string][]float64)
			}
		}
		model[word] = floatVector
	}
	return model
}

// random - helper for tests
func random(min, max int) int {
	return rand.Intn(max-min) + min
}
