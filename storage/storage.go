package storage

import (
	"github.com/VagabondDataNinjas/gizlinebot/domain"
)

type Storage interface {
	// used for after-the-fact debugging
	AddRawLineEvent(eventType, rawMsg string) error
	GetUserProfile(userId string) (domain.UserProfile, error)
	MarkProfileBotSurveyInited(userId string) error
	GetUsersWithoutAnswers(sinceUnixTs int64) (userIds []string, err error)
	AddUserProfile(userId, displayName string) error
	UserHasAnswers(userId string) (bool, error)
	UserGetLastAnswer(userId string) (domain.Answer, error)
	UserAddAnswer(domain.Answer) error
	UserAddGpsAnswer(domain.AnswerGps) error
	GetQuestions() (*domain.Questions, error)
	GetWelcomeMsgs(tplVars *WelcomeMsgTplVars) ([]string, error)
	WipeUser(userId string) error
	// cleanup any connection / file descriptors to the storage
	Close() error
}
