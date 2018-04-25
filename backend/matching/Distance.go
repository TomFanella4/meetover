package matching

import (
	"fmt"
	"io/ioutil"
	"math"
	"meetover/backend/user"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// getStopWords - gets the stop words from file
func getStopWords(fileName string) map[string]bool {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	swarr := strings.Split(str, "\n")
	swm := make(map[string]bool)
	for _, sw := range swarr {
		swm[sw] = true
	}
	return swm
}

// StripStopWords -
func StripStopWords(str []string) []string {
	stopFile := DataDir + "stopwords.txt"
	stopWords := getStopWords(stopFile)
	for i, w := range str {
		if _, exists := stopWords[strings.ToLower(w)]; exists {
			// fmt.Println("found :" + w)
			str = append(str[:i], str[i+1:]...)
		}
	}
	return str
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
	// fmt.Printf("distance: %f\n", d)
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
	par = StripStopWords(par)
	for i := 0; i < WordModelRandomParam; {
		w := par[random(0, len(par))]
		w = strings.TrimSpace(strings.ToLower(w))
		if len(w) <= 1 {
			continue
		}
		if val, found := model[w]; found {
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
