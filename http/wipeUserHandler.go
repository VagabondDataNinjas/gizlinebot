package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

func WipeUserHandlerBuilder(s *storage.Sql, lineBot *linebot.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		userIds := strings.Split(c.Param("userid"), ",")

		if len(userIds) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "error", "err": "Missing user id param",
			})
		}

		warnings := make([]string, 0)
		for _, userId := range userIds {
			err := s.WipeUser(strings.TrimSpace(userId))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"status": "error", "err": err.Error(),
				})
			}

			log.WithFields(log.Fields{
				"UserId": userId,
			}).Info("user data wiped")
			lineMsg := linebot.NewTextMessage("Successfully removed your data from GROOTS DB.")
			if _, err := lineBot.PushMessage(userId, lineMsg).Do(); err != nil {
				warn := fmt.Sprintf("Got error when seding msg to %s: %s", userId, err)
				log.Warn(warn)
				warnings = append(warnings, warn)
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success", "warnings": warnings,
		})
	}
}
