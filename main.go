package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	loggly "github.com/jamespearly/loggly"
	uuid "github.com/satori/go.uuid"
)

type ObjPlayer struct {
	UuID           uuid.UUID   `json:"uuid"`
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
}
type ObjArea struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ObjCompetition struct {
	ID          int       `json:"id"`
	Area        ObjArea   `json:"area"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Plan        string    `json:"plan"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type ObjSeason struct {
	ID              int         `json:"id"`
	StartDate       string      `json:"startDate"`
	EndDate         string      `json:"endDate"`
	CurrentMatchday int         `json:"currentMatchday"`
	Winner          interface{} `json:"winner"`
}

type ObjTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Count   int `json:"count"`
	Filters struct {
		Limit int `json:"limit"`
	} `json:"filters"`
	Competition ObjCompetition `json:"competition"`
	Season      ObjSeason      `json:"season"`
	Scorers     []struct {
		Player        ObjPlayer `json:"player"`
		Team          ObjTeam   `json:"team"`
		NumberOfGoals int       `json:"numberOfGoals"`
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
	insertIntoDynamoDB(responseObject)

}

// functio to execut method with given time
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func insertIntoDynamoDB(res Response) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// insert players into db
	for i := 0; i < len(res.Scorers); i++ {
		u1 := uuid.Must(uuid.NewV4())
		res.Scorers[i].Player.UuID = u1
		av, err := dynamodbattribute.MarshalMap(res.Scorers[i].Player)
		if err != nil {
			fmt.Println("Got error marshalling map:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Top10"),
		}

		_, err = svc.PutItem(input)

		if err != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Successfully added player")
	}
}

// main function to start the program
func main() {
	doEvery(10*time.Second, printTop10Scorer)
}
