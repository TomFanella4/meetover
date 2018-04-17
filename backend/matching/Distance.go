package matching

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"meetover/backend/user"
	"regexp"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// getStopWords - gets the stop words from file
func getStopWords(fileName string) []string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	res := strings.Split(str, "\n")
	// end := len(res) - 1
	// return res[:end]
	return res
}

// stripStopWords -
func StripStopWords(str []string) string {
	str = strings.ToLower(str)
	stopFile := DataDir + "stopwords.txt"
	stopWords := getStopWords(stopFile)
	// i := 172
	// sw := stopWords[i]
	// fmt.Println(len(sw))
	// fmt.Println("err: ---" + string(sw[4]) + "---")
	// fmt.Println("removing: *" + sw + "*")
	spaceRep, _ := regexp.Compile("[ \t\r\n\v\f]+")
	res := ""
	for _, sw := range stopWords {
		fmt.Print("removing " + sw)
		stopRep, err := regexp.Compile("[ \t\r\n\v\f]*" + sw + "[ \t\r\n\v\f]*")
		if err != nil {
			temp := stopRep.ReplaceAllString(str, " ")
			res = temp
		} else {
			fmt.Println("Unable to regex for: " + sw)
			fmt.Println(len(sw))
			fmt.Println(err.Error())
			break
		}

	}
	res = spaceRep.ReplaceAllString(res, " ")
	// r := strings.NewReplacer(sw, "")
	// res := r.Replace(str)
	return res
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

// random - helper for tests
func random(min, max int) int {
	return rand.Intn(max-min) + min
}
