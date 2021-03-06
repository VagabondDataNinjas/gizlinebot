package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

type AnswersRequest struct {
	UserId  string       `json:"user_id"`
	Answers []AnswerItem `json:"answers"`
}

type AnswerItem struct {
	QuestionId string `json:"question_id"`
	Answer     string `json:"answer"`
}

func AnswerHandlerBuilder(s *storage.Sql, bot *linebot.Client) func(c echo.Context) error {
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

		if !profile.SurveyStarted {
			s.ToggleUserSurvey(profile.UserId, true)
		}
		// @TODO check that the question_ids exist in the questions table

		for _, answer := range payload.Answers {
			if answer.QuestionId == "price" {
				err = sendPriceList(bot, s, profile.UserId)
				if err != nil {
					log.Errorf("Error sending prices to user: %s - %s", profile.UserId, err)
				}

				go func(bot *linebot.Client, s *storage.Sql, userId string) {
					time.Sleep(5 * time.Second)
					err = sendThankYouMsg(bot, s, profile.UserId)
					if err != nil {
						log.Errorf("Error sending thank you msg to user: %s - %s", profile.UserId, err)
					}
				}(bot, s, profile.UserId)
			}

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
