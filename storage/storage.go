package storage

import (
	"github.com/VagabondDataNinjas/gizlinebot/domain"
)

type Storage interface {
	// used for after-the-fact debugging
	AddRawLineEvent(eventType, rawMsg string) error
	AddUserProfile(userId, displayName string) error
	UserHasAnswers(userId string) (bool, error)
	UserGetLastAnswer(userId string) (domain.Answer, error)
	UserAddAnswer(domain.Answer) error
	// cleanup any connection / file descriptors to the storage
	Close() error
}
