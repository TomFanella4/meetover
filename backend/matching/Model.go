package matching

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mpace965/word-embedding/builder"
)

// WordModel is the vector representation of the words in the corpus file
var WordModel map[string][]float64

// WordModelContextWindow -
var WordModelContextWindow = 20

// WordModelDimension -
var WordModelDimension = 8

// WordModelRandomParam - number of words considered for similarity
var WordModelRandomParam = 100

// InitMLModel check if model has been created or creates it
func InitMLModel(windowSize int, wordDimensions int) {
	modelFile := DataDir + "meetOver.model"
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		fmt.Println("Model does not exist. Creating Model")
		corpusFile := DataDir + "corpus.dat"
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
