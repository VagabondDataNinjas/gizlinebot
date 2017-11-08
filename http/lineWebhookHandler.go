package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
)

func LineWebhookHandlerBuilder(surv *survey.Survey, s *storage.Sql, bot *linebot.Client, globalVars *domain.GlobalTplVars) func(c echo.Context) error {
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
				log.Errorf("Could not marshal event: %+v; err: %s", event, err)
			} else {
				err = s.AddRawLineEvent(string(event.Type), string(eventString))
				if err != nil {
					log.Errorf("Could not store event: %+v; err: %s", event, err)
				}
			}

			if event.Type == linebot.EventTypeUnfollow {
				err = onUnfollow(userId, s)
				log.Infof("Unfollow event: %s", userId)
				if err != nil {
					log.Error(err)
					continue
				}
			}

			if event.Type == linebot.EventTypeFollow {
				log.Infof("Follow event: %s", userId)
				userProfileResp, err := bot.GetProfile(userId).Do()
				if err != nil {
					log.Error(err)
					continue
				}

				err = s.AddUpdateUserProfile(userProfileResp.UserID, userProfileResp.DisplayName)
				if err != nil {
					log.Errorf("AddUpdateUserProfile err: %s", err)
					continue
				}

				profile, err := s.GetUserProfile(userId)
				if err != nil {
					log.Error(err)
					continue
				}

				if profile.SurveyStarted {
					question, err := surv.GetNextQuestion(userId, questionTplVars)
					if err != nil {
						log.Error(err)
						continue
					}
					if _, err = bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
						log.Error(err)
						continue
					}
				} else {
					log.Infof("User %s: sending welcome messages", userId)
					err = sendWelcomeMsgs(userId, bot, s, globalVars)
					if err != nil {
						log.Error(err)
						continue
					}
				}
			}

			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.LocationMessage:
					err = surv.RecordGpsAnswer(userId, message.Latitude, message.Longitude, message.Address, "line")
					if err != nil {
						log.Error(err)
						break
					}

				case *linebot.TextMessage:
					answer, err := surv.RecordAnswer(userId, message.Text, "line")
					if err != nil {
						log.Error(err)
						break
					}

					if answer.QuestionId == "price" {
						err = sendPriceList(bot, s, userId)
						if err != nil {
							log.Error(err)
						}
					}

					question, err := surv.GetNextQuestion(userId, questionTplVars)
					if err != nil {
						log.Error(err)
						break
					}

					if _, err = bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
						log.Error(err)
					}
				}
			}
		}
		return nil
	}
}

func onUnfollow(userId string, s *storage.Sql) error {
	p, err := s.GetUserProfile(userId)
	if err != nil {
		return err
	}

	p.Active = false
	p.SurveyStarted = false
	return s.UpdateUserProfile(p)
}

func sendWelcomeMsgs(userId string, bot *linebot.Client, s *storage.Sql, globalVars *domain.GlobalTplVars) error {
	welcomeTplVars := &storage.WelcomeMsgTplVars{
		UserId:   userId,
		Hostname: globalVars.Hostname,
	}
	welcomeMsgs, err := s.GetWelcomeMsgs(welcomeTplVars)
	if err != nil {
		return err
	}

	for _, welcomeMsg := range welcomeMsgs {
		// check if the message is a video one
		isVideoMsg, _ := regexp.MatchString(".*.mp4", welcomeMsg)
		if isVideoMsg {
			vidMsgRegex := regexp.MustCompile("\\|")
			vidAndPreview := vidMsgRegex.Split(welcomeMsg, -1)
			if len(vidAndPreview) != 2 {
				return errors.New(fmt.Sprintf("Unexpected video message format. Got: \"%s\"", welcomeMsg))
			}
			if _, err := bot.PushMessage(userId, linebot.NewVideoMessage(vidAndPreview[0], vidAndPreview[1])).Do(); err != nil {
				return err
			}
			continue
		}

		// check if the message is a button one
		isWebSurveyBtn, _ := regexp.MatchString("web-survey-btn", welcomeMsg)
		if isWebSurveyBtn {
			cfg, err := s.GetWebSurveyBtnConfig()
			if err != nil {
				return err
			}
			surveyUrl := globalVars.Hostname + "/?uid=" + userId
			surveyImgPath := globalVars.Hostname + "/lineimgs/" + cfg.ImageName
			action := linebot.NewURITemplateAction(cfg.Label, surveyUrl)
			buttonTemplate := linebot.NewButtonsTemplate(surveyImgPath, cfg.Title, cfg.Text, action)
			altText := cfg.Title
			if altText == "" {
				altText = "Groots"
			}
			buttonMsg := linebot.NewTemplateMessage(altText, buttonTemplate)
			if _, err = bot.PushMessage(userId, buttonMsg).Do(); err != nil {
				return err
			}
			continue
		}

		if _, err := bot.PushMessage(userId, linebot.NewTextMessage(welcomeMsg)).Do(); err != nil {
			return err
		}
	}
	return nil
}
