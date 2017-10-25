package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/labstack/echo"
)

type SendLineMsgRequest struct {
	UserIds []string `json:"user_ids"`
	Message string   `json:"message"`
}

func SendLineMsgHandlerBuilder(s storage.Storage, lineBot *linebot.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(SendLineMsgRequest)

		if err := c.Bind(payload); err != nil {
			return err
		}

		// @TODO validate UserIds exist
		// @TODO segment users list in max 150 users @see line documentation
		// @TODO use multicast
		warnings := []string{}
		lineMsg := linebot.NewTextMessage(payload.Message)
		for _, userId := range payload.UserIds {
			if _, err := lineBot.PushMessage(userId, lineMsg).Do(); err != nil {
				warn := fmt.Sprintf("Got error when seding msg to %s: %s", userId, err)
				log.Print(warn)
				warnings = append(warnings, warn)
				continue
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success", "warnings": warnings,
		})
	}
}
