package domain

import "time"

type Answer struct {
	Id         uint
	UserId     string
	QuestionId string
	Answer     string
	Channel    string
	Timestamp  time.Time
}
