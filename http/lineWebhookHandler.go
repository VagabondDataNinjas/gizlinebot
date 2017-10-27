package http

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/labstack/echo"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/line/line-bot-sdk-go/linebot"
)

func LineWebhookHandlerBuilder(surv *survey.Survey, s storage.Storage, bot *linebot.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		fmt.Printf("\nLineWebHook!!!\n")
		events, err := bot.ParseRequest(c.Request())
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.JSON(400, map[string]string{"status": "error", "err": err.Error()})
				// w.WriteHeader(400)
			} else {
				// @TODO log error
				c.JSON(500, map[string]string{"status": "error"})
			}
		}
		for _, event := range events {
			userId := event.Source.UserID
			eventString, err := json.Marshal(event)
			if err != nil {
				log.Printf("[err] Could not marshal event: %+v; err: %s", event, err)
			} else {
				err = s.AddRawLineEvent(string(event.Type), string(eventString))
				if err != nil {
					log.Printf("[err] Could not store event: %+v; err: %s", event, err)
				}
			}

			if event.Type == linebot.EventTypeFollow {
				userProfileResp, err := bot.GetProfile(userId).Do()
				if err != nil {
					log.Print(err)
					continue
				}

				err = s.AddUserProfile(userProfileResp.UserID, userProfileResp.DisplayName)
				if err != nil {
					fmt.Printf("AddUserProfile err: %s\n", err)
					continue
				}

				hasAnswers, err := s.UserHasAnswers(userId)
				if err != nil {
					log.Print(err)
					continue
				}

				if hasAnswers {
					question, err := surv.GetNextQuestion(userId)
					if err != nil {
						log.Print(err)
						continue
					}
					if _, err = bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
						log.Print(err)
						continue
					}
				} else {
					welcomeTplVars := &storage.WelcomeMsgTplVars{
						UserId: userId,
					}
					welcomeMsgs, err := s.GetWelcomeMsgs(welcomeTplVars)
					if err != nil {
						log.Print(err)
						continue
					}
					for _, welcomeMsg := range welcomeMsgs {
						if _, err = bot.PushMessage(userId, linebot.NewTextMessage(welcomeMsg)).Do(); err != nil {
							log.Print(err)
							continue
						}
					}
				}
			}

			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.LocationMessage:
					err = surv.RecordGpsAnswer(userId, message.Latitude, message.Longitude, message.Address, "line")
					if err != nil {
						log.Print(err)
						break
					}
				case *linebot.TextMessage:
					err = surv.RecordAnswer(userId, message.Text, "line")
					if err != nil {
						log.Print(err)
						break
					}

					question, err := surv.GetNextQuestion(userId)
					if err != nil {
						log.Print(err)
						break
					}

					if _, err = bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
		return nil
	}
}
