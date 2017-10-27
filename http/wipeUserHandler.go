package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

func WipeUserHandlerBuilder(s storage.Storage, lineBot *linebot.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		if userId == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "error", "err": "Missing user id param",
			})
		}

		err := s.WipeUser(userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "error", "err": err.Error(),
			})
		}

		warnings := make([]string, 0)
		lineMsg := linebot.NewTextMessage("Successfully removed your profile and content. Block and add me again to start another test.")
		if _, err := lineBot.PushMessage(userId, lineMsg).Do(); err != nil {
			warn := fmt.Sprintf("Got error when seding msg to %s: %s", userId, err)
			log.Print(warn)
			warnings = append(warnings, warn)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success", "warnings": warnings,
		})
	}
}
