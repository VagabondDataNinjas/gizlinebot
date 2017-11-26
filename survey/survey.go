package survey

import (
	"bytes"
	"html/template"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

type QuestionTemplateVars struct {
	UserId   string
	Hostname string
}

type Survey struct {
	Storage   *storage.Sql
	Questions *domain.Questions
}

func NewSurvey(storage *storage.Sql, questions *domain.Questions) (survey *Survey) {
	return &Survey{
		Storage:   storage,
		Questions: questions,
	}
}

func (s *Survey) GetNextQuestion(userId string, tplVars *QuestionTemplateVars) (question *domain.Question, err error) {
	q, err := s.getNextQuestionRaw(userId)
	if err != nil {
		return question, err
	}
	tmpl, err := template.New("lineQuestion").Parse(q.Text)
	if err != nil {
		return question, err
	}
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, tplVars)
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

	profile, err := s.Storage.GetUserProfile(userId)
	if err != nil {
		return question, err
	}

	if !profile.SurveyStarted || !has {
		return s.Questions.At(0)
	}

	answer, err := s.Storage.UserGetLastAnswer(userId)
	if err != nil {
		return question, err
	}

	return s.Questions.Next(answer.QuestionId)
}

func (s *Survey) RecordAnswerRaw(userId, questionId, answerText, channel string) error {
	answer := domain.Answer{
		UserId:     userId,
		QuestionId: questionId,
		Answer:     answerText,
		Channel:    channel,
	}
	return s.Storage.UserAddAnswer(answer)
}

func (s *Survey) RecordAnswer(userId, answerText, channel string) (domain.Answer, error) {
	has, err := s.Storage.UserHasAnswers(userId)
	if err != nil {
		return domain.Answer{}, err
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
		if err = s.Storage.UserAddAnswer(answer); err != nil {
			return answer, err
		}
		return answer, nil
	}

	prevAnswer, err := s.Storage.UserGetLastAnswer(userId)
	if err != nil {
		return answer, err
	}
	currentQ, err := s.Questions.Next(prevAnswer.QuestionId)
	if err != nil {
		// if the user already answered all the questions
		// record against a n 'na' question id
		if err == domain.ErrQuestionsNoNext {
			answer.QuestionId = "na"
			if err = s.Storage.UserAddAnswer(answer); err != nil {
				return answer, err
			}
			return answer, nil
		}

		return answer, err
	}
	answer.QuestionId = currentQ.Id
	s.Storage.UserAddAnswer(answer)
	return answer, nil
}

func (s *Survey) RecordGpsAnswer(userId string, lat, lon float64, address string, channel string) (err error) {
	answer := domain.AnswerGps{
		UserId:  userId,
		Lat:     lat,
		Lon:     lon,
		Address: address,
		Channel: channel,
	}
	return s.Storage.UserAddGpsAnswer(answer)
}
