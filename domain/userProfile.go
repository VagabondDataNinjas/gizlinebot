package domain

type UserProfile struct {
	UserId        string
	DisplayName   string
	Timestamp     int
	SurveyStarted bool
	// if false: user blocked the bot
	Active bool
}
