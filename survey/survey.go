package survey

import (
	"bytes"
	"html/template"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

type QuestionTemplateVars struct {
	UserId string
}

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
	q, err := s.getNextQuestionRaw(userId)
	if err != nil {
		return question, err
	}
	tmpl, err := template.New("lineQuestion").Parse(q.Text)
	if err != nil {
		return question, err
	}
	buf := new(bytes.Buffer)

	vars := QuestionTemplateVars{
		UserId: userId,
	}
	err = tmpl.Execute(buf, vars)
	if err != nil {
		return question, err
	}

	q.Text = buf.String()
	return q, nil
}

func (s *Survey) getNextQuestionRaw(userId string) (question *domain.Question, err error) {
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

func (s *Survey) RecordAnswer(userId, answerText, channel string) (err error) {
	has, err := s.Storage.UserHasAnswers(userId)
	if err != nil {
		return err
	}

	answer := domain.Answer{
		UserId:  userId,
		Answer:  answerText,
		Channel: channel,
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
