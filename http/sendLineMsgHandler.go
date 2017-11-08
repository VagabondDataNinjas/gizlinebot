package http

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

type SendLineMsgRequest struct {
	UserIds  []string `json:"user_ids"`
	Messages []string `json:"messages"`
}

type MsgVars struct {
	UserId string
}

func SendLineMsgHandlerBuilder(s storage.Storage, lineBot *linebot.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(SendLineMsgRequest)

		if err := c.Bind(payload); err != nil {
			return err
		}

		// @TODO validate UserIds exist
		// @TODO validate that the message templates compile
		// @TODO segment users list in max 150 users @see line documentation
		// @TODO use multicast
		warnings := []string{}
		for _, userId := range payload.UserIds {
			for _, message := range payload.Messages {
				vars := MsgVars{
					UserId: userId,
				}
				tmpl, err := template.New("lineMsg").Parse(message)
				if err != nil {
					// @TODO log
					return err
				}
				buf := new(bytes.Buffer)
				err = tmpl.Execute(buf, vars)
				if err != nil {
					return err
				}

				lineMsg := linebot.NewTextMessage(buf.String())
				if _, err := lineBot.PushMessage(userId, lineMsg).Do(); err != nil {
					warn := fmt.Sprintf("Got error when seding msg to %s: %s", userId, err)
					log.Error(warn)
					warnings = append(warnings, warn)
					continue
				}
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success", "warnings": warnings,
		})
	}
}
