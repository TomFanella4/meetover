package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
)

// RawUsers - no. of unprocessed users generate by the json generator tool

func updateJSONFile(newJSON interface{}, fileName string) {
	bytes, err := json.Marshal(newJSON)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = ioutil.WriteFile(fileName, bytes, 0644)
	return
}
func getJobJSON(fileName string) JobData {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	// fmt.Printf("Read: %s\n", string(file))
	var jobd JobData
	json.Unmarshal(file, &jobd)
	// fmt.Printf("Results: %v\n", jsontype)
	return jobd
}

func getRawUsers(rawFile string) []User {
	raw, err := ioutil.ReadFile(rawFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var users []User
	json.Unmarshal(raw, &users)
	return users
}

// GenTestUsers -
func GenTestUsers(sourceFile, sinkFile string) {

	numTech := 200
	rawUsers := getRawUsers(sourceFile)
	techUsers := genTechUsers(rawUsers[:numTech])
	fmt.Println(techUsers)
}

func jobDataToUser(jd JobData, rawUser User) User {
	js := jobSummary(jd)
	r := rand.Intn
	start := r(len(js.Description)) / 2
	rawUser.Profile.Greeting = js.Description[start:]
	rawUser.Profile.Headline = js.Titles[r(len(js.Titles))]
	rawUser.Profile.Industry = js.Skills[r(len(js.Skills))]
	rawUser.Profile.FormattedName = rawUser.Profile.FirstName + " " + rawUser.Profile.LastName
	start = r(len(js.Description)) / 2
	rawUser.Profile.Summary = js.Description[start:]
	return rawUser
}

func otherJobFiles() []string {
	files, err := ioutil.ReadDir("./otherJobs")
	if err != nil {
		log.Fatal(err)
	}
	res := []string{}
	for _, f := range files {
		res = append(res, "./otherJobs/"+f.Name())
	}
	return res
}
func genNonTechUsers(rawUsers []User) []User {
	nu := 0
	files := otherJobFiles()
	for _, f := range files {
		jd := getJobJSON(f) // make 100 users for each type of job data
		for ; nu%100 > 0 || nu == 0; nu++ {
			rawUsers[nu] = jobDataToUser(jd, rawUsers[nu])
		}
	}
	return rawUsers
}
func genTechUsers(rawUsers []User) []User {
	csvFile, err := os.Open("./jobs.csv")
	if err != nil {
		fmt.Println(err.Error())
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	headers, error := reader.Read() // headers
	if error == io.EOF {
		fmt.Println("OEF in dataset")
	}
	fmt.Println("Headers: ")
	fmt.Println(headers) // 1 - desc, 3 - title, 4 - skills
	r := rand.Intn
	textLength := 250
	for i, u := range rawUsers {
		line, error := reader.Read()
		if error == io.EOF {
			fmt.Println("OEF in dataset")
		}
		description := line[1]
		title := line[2]
		skills := line[3]
		n := len(description)
		profile := u.Profile
		start := r(n-1) / 2
		greeting := description[start:]
		if len(greeting) > textLength {
			profile.Greeting = greeting[:textLength]
		} else {
			profile.Greeting = greeting
		}
		profile.Headline = title
		profile.Industry = strings.Replace(skills, ",", " , ", -1)
		profile.FormattedName = profile.FirstName + " " + profile.LastName
		start = r(n - 1)
		summary := description[start:]
		if len(summary) > textLength {
			profile.Summary = summary[:textLength]
		} else {
			profile.Summary = summary
		}
		rawUsers[i] = u
		rawUsers[i].Profile = profile
	}
	return rawUsers
}

// Geolocation - latitide and longitude and last time of update
type Geolocation struct {
	Lat       float64 `json:"lat,omitempty"`
	Long      float64 `json:"long,omitempty"`
	TimeStamp int64   `json:"timestamp,omitempty"`
}

// ATokenResponse from Code exchange with LI
type ATokenResponse struct {
	AToken string `json:"access_token"`
	Expiry uint   `json:"expires_in"`
}

// Profile is the JSON object we user to store our user data
type Profile struct {
	CurrentShare struct {
		Attribution struct {
			Share struct {
				Author struct {
					FirstName string `json:"firstName"`
					ID        string `json:"id"`
					LastName  string `json:"lastName"`
				} `json:"author"`
				Comment string `json:"comment"`
				ID      string `json:"id"`
			} `json:"share"`
		} `json:"attribution"`
		Author struct {
			FirstName string `json:"firstName"`
			ID        string `json:"id"`
			LastName  string `json:"lastName"`
		} `json:"author"`
		Comment string `json:"comment"`
		ID      string `json:"id"`
		Source  struct {
			ServiceProvider struct {
				Name string `json:"name"`
			} `json:"serviceProvider"`
		} `json:"source"`
		Timestamp  int64 `json:"timestamp"`
		Visibility struct {
			Code string `json:"code"`
		} `json:"visibility"`
	} `json:"currentShare"`
	EmailAddress  string `json:"emailAddress"`
	FirstName     string `json:"firstName"`
	FormattedName string `json:"formattedName"`
	Headline      string `json:"headline"`
	ID            string `json:"id"`
	Industry      string `json:"industry"`
	LastName      string `json:"lastName"`
	Location      struct {
		Country struct {
			Code string `json:"code"`
		} `json:"country"`
		Name string `json:"name"`
	} `json:"location"`
	NumConnections int    `json:"numConnections"`
	PictureURL     string `json:"pictureUrl"`
	Positions      struct {
		Total  int `json:"_total"`
		Values []struct {
			Company struct {
				ID       int    `json:"id"`
				Industry string `json:"industry"`
				Name     string `json:"name"`
				Size     string `json:"size"`
				Type     string `json:"type"`
			} `json:"company"`
			ID        int  `json:"id"`
			IsCurrent bool `json:"isCurrent"`
			Location  struct {
				Name string `json:"name"`
			} `json:"location"`
			StartDate struct {
				Month int `json:"month"`
				Year  int `json:"year"`
			} `json:"startDate"`
			Summary string `json:"summary"`
			Title   string `json:"title"`
		} `json:"values"`
	} `json:"positions"`
	Summary       string `json:"summary"`
	ShareLocation bool   `json:"shareLocation"`
	Greeting      string `json:"greeting"`
}

// User is user on MeetOver
type User struct {
	ID           string         `json:"uid,omitempty"`
	Location     *Geolocation   `json:"location,omitempty"`
	AccessToken  ATokenResponse `json:"accessToken"`
	Profile      Profile        `json:"profile"`
	IsSearching  bool           `json:"isSearching"`
	IsMatchedNow bool           `json:"isMatched"` // set directly from the mobile app
}

// JobSummary -
type JobSummary struct {
	Titles      []string `json:"titles"`
	Description string   `json:"description"`
	Skills      []string `json:"skills"`
}

func jobSummary(jd JobData) JobSummary {
	var res JobSummary
	res.Description = ""
	res.Skills = []string{}
	res.Titles = jd.Occupation.SampleOfReportedJobTitles.Title
	for _, v := range jd.Tasks.Task {
		res.Description += v.Name + " "
	}
	for _, v := range jd.TechnologySkills.Category {
		res.Skills = append(res.Skills, v.Title.Name)
		for _, v2 := range v.Example {
			res.Skills = append(res.Skills, v2.Name)
		}
	}
	for _, v := range jd.ToolsTechnology.Technology.Category {
		res.Skills = append(res.Skills, v.Title.Name)
		for _, v2 := range v.Example {
			res.Skills = append(res.Skills, v2.Name)
		}
	}
	for _, v := range jd.Knowledge.Element {
		res.Skills = append(res.Skills, v.Name)
		res.Description += v.Description + " "
	}
	for _, v := range jd.Skills.Element {
		res.Skills = append(res.Skills, v.Name)
		res.Description += v.Description + " "
	}
	for _, v := range jd.Abilities.Element {
		res.Skills = append(res.Skills, v.Name)
		res.Description += v.Description + " "
	}
	for _, v := range jd.WorkActivities.Element {
		res.Skills = append(res.Skills, v.Name)
		res.Description += v.Description + " "
	}
	for _, v := range jd.DetailedWorkActivities.Activity {
		res.Description += v.Name + " "
	}
	for _, v := range jd.RelatedOccupations.Occupation {
		res.Titles = append(res.Titles, v.Title)
	}
	for _, v := range jd.AdditionalInformation.Source {
		res.Description += v.Name + " "
	}
	return res
}

// JobData -
type JobData struct {
	Occupation struct {
		Code  string `json:"code"`
		Title string `json:"title"` // job name
		Tags  struct {
			BrightOutlook bool `json:"bright_outlook"`
			Green         bool `json:"green"`
		} `json:"tags"`
		Description               string `json:"description"` //d
		SampleOfReportedJobTitles struct {
			Title []string `json:"title"` // title list
		} `json:"sample_of_reported_job_titles"`
	} `json:"occupation"`
	Tasks struct {
		Task []struct {
			ID      int    `json:"id"`
			Green   bool   `json:"green"`
			Related string `json:"related"`
			Name    string `json:"name"` // d
		} `json:"task"`
	} `json:"tasks"`
	TechnologySkills struct {
		Category []struct {
			Related string `json:"related"`
			Title   struct {
				ID   int    `json:"id"`
				Name string `json:"name"` //skill
			} `json:"title"`
			Example []struct {
				HotTechnology int    `json:"hot_technology,omitempty"`
				Name          string `json:"name"` // skill
			} `json:"example"`
		} `json:"category"`
	} `json:"technology_skills"`
	ToolsTechnology struct {
		Tools struct {
			Category []struct {
				Related string `json:"related"`
				Title   struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"title"`
				Example []struct {
					Name string `json:"name"` // skill
				} `json:"example"`
			} `json:"category"`
		} `json:"tools"`
		Technology struct {
			Category []struct {
				Related string `json:"related"`
				Title   struct {
					ID   int    `json:"id"`
					Name string `json:"name"` // skill
				} `json:"title"`
				Example []struct {
					HotTechnology int    `json:"hot_technology,omitempty"`
					Name          string `json:"name"` // skill
				} `json:"example"`
			} `json:"category"`
		} `json:"technology"`
	} `json:"tools_technology"`
	Knowledge struct {
		Element []struct {
			ID          string `json:"id"`
			Related     string `json:"related"`
			Name        string `json:"name"`        // s
			Description string `json:"description"` // desc
		} `json:"element"`
	} `json:"knowledge"`
	Skills struct {
		Element []struct {
			ID          string `json:"id"`
			Related     string `json:"related"`
			Name        string `json:"name"`        // s
			Description string `json:"description"` // desc
		} `json:"element"`
	} `json:"skills"`
	Abilities struct {
		Element []struct {
			ID          string `json:"id"`
			Related     string `json:"related"`
			Name        string `json:"name"`        // s
			Description string `json:"description"` // d
		} `json:"element"`
	} `json:"abilities"`
	WorkActivities struct {
		Element []struct {
			ID          string `json:"id"`
			Related     string `json:"related"`
			Name        string `json:"name"`        // s
			Description string `json:"description"` // d
		} `json:"element"`
	} `json:"work_activities"`
	DetailedWorkActivities struct {
		Activity []struct {
			ID      string `json:"id"`
			Related string `json:"related"`
			Name    string `json:"name"` // d
		} `json:"activity"`
	} `json:"detailed_work_activities"`
	RelatedOccupations struct {
		Occupation []struct {
			Href  string `json:"href"`
			Code  string `json:"code"`
			Title string `json:"title"` // title
			Tags  struct {
				BrightOutlook bool `json:"bright_outlook"`
				Green         bool `json:"green"`
			} `json:"tags"`
		} `json:"occupation"`
	} `json:"related_occupations"`
	AdditionalInformation struct {
		Source []struct {
			URL  string `json:"url"`
			Name string `json:"name"` //desc
		} `json:"source"`
	} `json:"additional_information"`
}
