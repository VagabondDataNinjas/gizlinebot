package domain

import "time"

type AnswerGps struct {
	Id        uint
	UserId    string
	Lat       float64
	Lon       float64
	Address   string
	Channel   string
	Timestamp time.Time
}
