package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/ynqa/word-embedding/builder"
	"gonum.org/v1/gonum/mat"
)

// MatchValue represents each perspective user in their distance from the caller
type MatchValue struct {
	Usr  Profile     `json:"profile"`
	Dist float64     `json:"distance"`
	Loc  Geolocation `json:"location"`
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

// GetMatches returns an ordered list of user uid's from closest to furthest to the caller
func GetMatches(UserID string, neighbors []User) (MatchResponse, error) {
	callingUser, err := GetUser(UserID)
	if err != nil {
		return MatchResponse{}, errors.New("Unable to fetch calling user")
	}
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
	start := time.Now()
	for _, pu := range prospUsers {
		prospStr := userToString(pu)
		prospVec := parToVector(prospStr, model)
		distance := nestedDistance(callerVec, prospVec)
		mr.Matches = append(mr.Matches, MatchValue{pu.Profile, distance, pu.Location})
	}
	sort.Sort(byDistance(mr.Matches))
	elapsed := time.Since(start)
	fmt.Println("Destance Calculation took: " + elapsed.String())
	return mr
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
func InitMLModel(windowSize int, wordDimensions int) {
	modelFile := "./data/meetOver.model"
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		fmt.Println("Model does not exist. Creating Model")
		corpusFile := "./data/corpus.dat"
		createModel(modelFile, corpusFile, windowSize, wordDimensions)
	}
	defer r.Close()
	WordModel = readModel(r)
	fmt.Println("Model Loaded")
}

// createModel uses the word2vec algo to create word embeddings
func createModel(modelObject *storage.ObjectHandle, corpusReader io.Reader, windowSize int, wordDimensions int) error {
	// Initialize word embeddings
	b := builder.NewWord2VecBuilder()
	b.SetDimension(wordDimensions).
		SetWindow(windowSize).
		SetModel("skip-gram").
		SetOptimizer("ns").
		SetNegativeSampleSize(15).
		SetVerbose()
	m, err := b.Build()
	if err != nil {
		return err
	}

	// Read corpus file
	fmt.Println("[+] Reading Corpus from Firebase Storage...")
	content, err := ioutil.ReadAll(corpusReader)
	if err != nil {
		return err
	}
	corpusByteReader := bytes.NewReader(content)

	fmt.Println("[+] Corpus loaded. Creating Model...")
	f1, err := m.Preprocess(corpusByteReader)
	if err != nil {
		return err
	}
	// Start to Train.
	m.Train(f1)
	f1.Close()

	// Save contents to disk temporarily because we can't access the model directly
	fmt.Println("[+] Model created. Saving Model to Firebase...")
	m.Save("temp")
	modelContent, err := ioutil.ReadFile("temp")
	if err != nil {
		return err
	}

	// Save word vectors to firebase storage.
	ctx := context.Background()
	modelWriter := modelObject.NewWriter(ctx)
	if _, err := fmt.Fprint(modelWriter, string(modelContent)); err != nil {
		return err
	}
	if err := modelWriter.Close(); err != nil {
		return err
	}
	if _, err := modelObject.Update(ctx, storage.ObjectAttrsToUpdate{
		ContentType: "application/octet-stream",
	}); err != nil {
		return err
	}

	// Remove temp contents
	if err := os.Remove("temp"); err != nil {
		return err
	}

	return nil
}

// readModel converts the generated model to an in-memory object
func readModel(modelReader io.Reader) map[string][]float64 {
	content, err := ioutil.ReadAll(modelReader)
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
