package domain

import "time"

type AnswerGps struct {
	Id        uint
	UserId    string
	Lat       float64
	Lon       float64
	Channel   string
	Timestamp time.Time
}
