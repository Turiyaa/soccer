package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	loggly "github.com/jamespearly/loggly"
)

// Response structs
type Response struct {
	Count   int `json:"count"`
	Filters struct {
		Limit int `json:"limit"`
	} `json:"filters"`

	// name of competitions
	Competition struct {
		ID   int `json:"id"`
		Area struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"area"`
		Name        string    `json:"name"`
		Code        string    `json:"code"`
		Plan        string    `json:"plan"`
		LastUpdated time.Time `json:"lastUpdated"`
	} `json:"competition"`

	// seasonYY/YY
	Season struct {
		ID              int         `json:"id"`
		StartDate       string      `json:"startDate"`
		EndDate         string      `json:"endDate"`
		CurrentMatchday int         `json:"currentMatchday"`
		Winner          interface{} `json:"winner"`
	} `json:"season"`

	//players info
	Scorers []struct {
		Player struct {
			ID             int         `json:"id"`
			Name           string      `json:"name"`
			FirstName      string      `json:"firstName"`
			LastName       interface{} `json:"lastName"`
			DateOfBirth    string      `json:"dateOfBirth"`
			CountryOfBirth string      `json:"countryOfBirth"`
			Nationality    string      `json:"nationality"`
			Position       string      `json:"position"`
			ShirtNumber    int         `json:"shirtNumber"`
			LastUpdated    time.Time   `json:"lastUpdated"`
		} `json:"player"`

		// name of the club
		Team struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"team"`
		// socreas
		NumberOfGoals int `json:"numberOfGoals"`
	} `json:"scorers"`
}

func printTop10Scorer(t time.Time) {
	//loogly token // TODO: set in the enviroment for security
	os.Setenv("LOGGLY_TOKEN", "4156602e-1451-4806-a3ad-80c982025fb1")
	tag := "SoccerScore"

	// Instantiate the loogly Client
	looglyClient := loggly.New(tag)

	// API end point urs (add nultiples)
	url := "https://api.football-data.org/v2/competitions/PD/scorers"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Auth-Token", "e42cb6a6ecc949c8897e06d284a55e05")

	// Valid EchoSend (message echoed to console and no error returned)
	logerr := looglyClient.EchoSend("info", "accessing json objects")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err:", logerr)
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(body, &responseObject)

	if responseObject.Count == 0 {
		// Valid Send (no error returned)
		logerr = looglyClient.Send("error", "object not found")
		fmt.Println("err:", logerr)
	}
	fmt.Print("Top ")
	fmt.Print(responseObject.Count)
	fmt.Print(" Scorers")
	fmt.Println()
	fmt.Println(responseObject.Competition.Name)
	fmt.Println("Name" + "\t\t\t" + "Goals")

	for i := 0; i < len(responseObject.Scorers); i++ {
		fmt.Print(responseObject.Scorers[i].Player.Name + "\t\t")
		fmt.Print(responseObject.Scorers[i].NumberOfGoals)
		fmt.Println()

	}

}

// functio to execut method with given time
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

// main function to start the program
func main() {
	doEvery(10*time.Second, printTop10Scorer)
}
