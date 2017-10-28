package http

import (
	"log"
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/domain"

	"github.com/VagabondDataNinjas/gizlinebot/storage"

	"github.com/labstack/echo"
)

type AnswersRequest struct {
	UserId  string       `json:"user_id"`
	Answers []AnswerItem `json:"answers"`
}

type AnswerItem struct {
	QuestionId string `json:"question_id"`
	Answer     string `json:"answer"`
}

func AnswerHandlerBuilder(s storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(AnswersRequest)

		if err := c.Bind(payload); err != nil {
			log.Printf("Bind error: %s", err)
			return err
		}

		profile, err := s.GetUserProfile(payload.UserId)
		if err != nil {
			log.Printf("Error fetching user profile %s: %s", payload.UserId, err)
			return err
		}
		if profile.UserId == "" {
			log.Printf("Missing user id %s", payload.UserId)
			return c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "reason": "user_id not found"})
		}

		// @TODO check that the question_ids exist in the questions table

		for _, answer := range payload.Answers {
			err := s.UserAddAnswer(domain.Answer{
				UserId:     payload.UserId,
				QuestionId: answer.QuestionId,
				Answer:     answer.Answer,
				Channel:    "web",
			})
			if err != nil {
				return err
			}
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}
