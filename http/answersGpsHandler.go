package http

import (
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/labstack/echo"
)

type AnswersGpsRequest struct {
	UserId   string        `json:"user_id"`
	Location AnswerGpsItem `json:"location"`
}

type AnswerGpsItem struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func AnswerGpsHandlerBuilder(s * storage.Sql) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(AnswersGpsRequest)

		if err := c.Bind(payload); err != nil {
			return err
		}

		profile, err := s.GetUserProfile(payload.UserId)
		if err != nil {
			return err
		}
		if profile.UserId == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "reason": "user_id not found"})
		}

		// @TODO check that the question_ids exist in the questions table

		err = s.UserAddGpsAnswer(domain.AnswerGps{
			UserId:  payload.UserId,
			Lat:     payload.Location.Lat,
			Lon:     payload.Location.Lon,
			Channel: "web",
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}
