package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"text/template"

	"github.com/VagabondDataNinjas/gizlinebot/domain"

	"github.com/labstack/echo"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/line/line-bot-sdk-go/linebot"
)

func LineWebhookHandlerBuilder(surv *survey.Survey, s storage.Storage, bot *linebot.Client, globalVars *domain.GlobalTplVars) func(c echo.Context) error {
	return func(c echo.Context) error {
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
			questionTplVars := &survey.QuestionTemplateVars{
				UserId:   userId,
				Hostname: globalVars.Hostname,
			}

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
				log.Printf("User %s follow event", userId)
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
					question, err := surv.GetNextQuestion(userId, questionTplVars)
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
						UserId:   userId,
						Hostname: globalVars.Hostname,
					}
					welcomeMsgs, err := s.GetWelcomeMsgs(welcomeTplVars)
					if err != nil {
						log.Print(err)
						continue
					}
					log.Printf("User %s: sending welcome messages", userId)
					for _, welcomeMsg := range welcomeMsgs {
						// check if the message is a video one
						isVideoMsg, _ := regexp.MatchString(".*.mp4", welcomeMsg)
						if isVideoMsg {
							vidMsgRegex := regexp.MustCompile("\\|")
							vidAndPreview := vidMsgRegex.Split(welcomeMsg, -1)
							if len(vidAndPreview) != 2 {
								log.Printf("Unexpected video message format. Got: \"%s\"", welcomeMsg)
								continue
							}
							if _, err := bot.PushMessage(userId, linebot.NewVideoMessage(vidAndPreview[0], vidAndPreview[1])).Do(); err != nil {
								log.Print(err)
								continue
							}
						} else {
							if _, err = bot.PushMessage(userId, linebot.NewTextMessage(welcomeMsg)).Do(); err != nil {
								log.Print(err)
								continue
							}
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
					answer, err := surv.RecordAnswer(userId, message.Text, "line")
					if err != nil {
						log.Print(err)
						break
					}

					if answer.QuestionId == "price" {
						err = sendPriceList(bot, s, userId)
						if err != nil {
							log.Print(err)
						}
					}

					question, err := surv.GetNextQuestion(userId, questionTplVars)
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

func sendPriceList(bot *linebot.Client, s storage.Storage, userId string) error {
	lp, err := s.GetUserNearbyPrices(userId)
	if err != nil {
		return err
	}

	tplStr, err := s.GetPriceTplMsg()
	if err != nil {
		return err
	}
	tmpl, err := template.New("priceMsg").Parse(tplStr)
	if err != nil {
		return err
	}

	type TplVars struct {
		Location string
		Price    float64
	}

	// @TODO check empty list of lp
	var msg string
	for _, loc := range lp {
		buf := new(bytes.Buffer)
		tplVars := TplVars{loc.Name, loc.Price}
		err = tmpl.Execute(buf, tplVars)
		if err != nil {
			return err
		}

		msg += buf.String()
		msg += "\n"
	}

	if _, err = bot.PushMessage(userId, linebot.NewTextMessage(msg)).Do(); err != nil {
		return err
	}

	return nil
}
