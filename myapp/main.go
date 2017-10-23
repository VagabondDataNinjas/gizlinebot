package myapp

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/smeeklai/gizlinebot/domain"
	"github.com/smeeklai/gizlinebot/line"
	"github.com/smeeklai/gizlinebot/storage"
	"github.com/smeeklai/gizlinebot/survey"
)

func init() {
	connectionName := mustGetenv("GIZLB_SQL_HOST")
	user := mustGetenv("GIZLB_SQL_USER")
	password := os.Getenv("GIZLB_SQL_PASS") // NOTE: password may be empty
	dbName := mustGetenv("GIZLB_SQL_DB")
	dbPort := os.Getenv("GIZLB_SQL_PORT") // NOTE: dbPort may be empty
	localDB, err := strconv.ParseBool(mustGetenv("LOCAL_DB"))
	if err != nil {
		fmt.Sprintf("\nLOCAL_DB value error")
	}

	var s *storage.Sql
	if localDB {
		s, err = storage.NewSql(user + ":" + password + "@(" + connectionName + ":" + dbPort + ")/" + dbName)
		if err != nil {
			fmt.Sprintf("\nCould not open db: %v", err)
		}
	} else {
		s, err = storage.NewGoogleSql(user, password, connectionName, dbName)
		if err != nil {
			fmt.Sprintf("\nCould not open db: %v", err)
		}
	}

	qs := domain.NewQuestions()
	questions := [][]string{
		{"Q1", `Thank you for following us!
			If you'd like to complete the survey online please go to https://google.com
			Otherwise you can complete the form here.
			You can start by replying back with your location (area or island name):`},
		{"Q2", "How much do you pay for diesel in your area?"},
		{"Q3", "What is your occupation?"},
		{"Q4", "What is your line id?"},
		{"Q5", "What did you do today?"},
		{"Q6", "Thank you for all your help! We might ask you more questions in the future"},
	}

	for _, question := range questions {
		err = qs.Add(question[0], question[1])
		if err != nil {
			fmt.Sprintf("\nError: %s\n", err)
		}
	}
	surv := survey.NewSurvey(s, qs)

	err = line.ServeAppEngine(s, surv, mustGetenv("CHANNEL_SECRET"), mustGetenv("CHANNEL_TOKEN"))
	if err != nil {
		fmt.Sprintf("\nError: %s\n", err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}
