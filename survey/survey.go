package survey

import (
	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

type Survey struct {
	Storage   storage.Storage
	Questions *domain.Questions
}

func NewSurvey(storage storage.Storage, questions *domain.Questions) (survey *Survey) {
	return &Survey{
		Storage:   storage,
		Questions: questions,
	}
}

func (s *Survey) GetNextQuestion(userId string) (question *domain.Question, err error) {
	has, err := s.Storage.UserHasAnswers(userId)
	if err != nil {
		return question, err
	}
	if !has {
		return s.Questions.At(0)
	}

	answer, err := s.Storage.UserGetLastAnswer(userId)
	if err != nil {
		return question, err
	}

	return s.Questions.Next(answer.QuestionId)
}

func (s *Survey) RecordAnswer(userId, answerText string) (err error) {
	has, err := s.Storage.UserHasAnswers(userId)
	if err != nil {
		return err
	}

	answer := domain.Answer{
		UserId: userId,
		Answer: answerText,
	}
	// if the user has not answered any of the questions
	// record this answer against the first question
	if !has {
		cq, _ := s.Questions.At(0)
		answer.QuestionId = cq.Id
		s.Storage.UserAddAnswer(answer)
		return nil
	}

	prevAnswer, err := s.Storage.UserGetLastAnswer(userId)
	if err != nil {
		return err
	}
	currentQ, err := s.Questions.Next(prevAnswer.QuestionId)
	if err != nil {
		// if the user already answered all the questions
		// record this answer against the last question id
		if err == domain.ErrQuestionsNoNext {
			lastQ, _ := s.Questions.Last()
			answer.QuestionId = lastQ.Id
			s.Storage.UserAddAnswer(answer)
			return nil
		}

		return err
	}
	answer.QuestionId = currentQ.Id
	s.Storage.UserAddAnswer(answer)
	return nil
}
