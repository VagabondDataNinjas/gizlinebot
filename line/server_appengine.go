package line

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

type LineAppEngine struct {
	Bot          *linebot.Client
	Storage      storage.Storage
	Survey       *survey.Survey
	Ctx          context.Context
	DuringSurvey bool
}

func ServeAppEngine(storage storage.Storage, surv *survey.Survey, secret, token string) error {
	handler, err := httphandler.New(secret, token)
	if err != nil {
		return err
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		ctx := appengine.NewContext(r)
		bot, err := handler.NewClient(linebot.WithHTTPClient(urlfetch.Client(ctx)))
		if err != nil {
			log.Printf("\nError: %s\n", err)
			aelog.Errorf(ctx, "%v", err)
			return
		}
		ls := &LineAppEngine{
			Bot:          bot,
			Storage:      storage,
			Survey:       surv,
			Ctx:          ctx,
			DuringSurvey: false,
		}
		ls.test(events)
	})
	http.Handle("/linewebhook", handler)
	http.HandleFunc("/", testHandler)

	return nil
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func (ls *LineAppEngine) checkErr(err error) {
	if err != nil {
		log.Print(err)
		aelog.Infof(ls.Ctx, "err: %s\n", err)
	}
}

func (ls *LineAppEngine) startSurvey(userId string) (replyMsg string, err error) {
	question, err := ls.Survey.GetNextQuestion(userId)
	if err != nil {
		return "", err
	}
	return question.Text, nil
}

func (ls *LineAppEngine) test(events []*linebot.Event) {
	for _, event := range events {
		userId := event.Source.UserID
		eventString, err := json.Marshal(event)
		if err != nil {
			log.Printf("[err] Could not marshal event: %+v; err: %s", event, err)
			aelog.Infof(ls.Ctx, "[err] Could not marshal event: %+v; err: %s", event, err)
		} else {
			log.Printf("\nEvent string: %s\n", eventString)
			err = ls.Storage.AddRawLineEvent(string(event.Type), string(eventString))
			if err != nil {
				log.Printf("[err] Could not store event: %+v; err: %s", event, err)
				aelog.Infof(ls.Ctx, "[err] Could not store event: %+v; err: %s", event, err)
			}
		}

		if event.Type == linebot.EventTypeFollow {
			userProfileResp, err := ls.Bot.GetProfile(userId).Do()
			if err != nil {
				log.Print(err)
				aelog.Infof(ls.Ctx, "err: %s\n", err)
				continue
			}

			err = ls.Storage.AddUserProfile(userProfileResp.UserID, userProfileResp.DisplayName)
			if err != nil {
				fmt.Printf("AddUserProfile err: %s\n", err)
				aelog.Infof(ls.Ctx, "AddUserProfile err: %s\n", err)
				continue
			}
		}

		if event.Type == linebot.EventTypeMessage {
			var replyMsg string
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				registeredStatus, err := ls.Storage.GetUserRegisteredStatus(userId)
				if err != nil {
					log.Print(err)
					aelog.Infof(ls.Ctx, "err: %s\n", err)
					break
				}
				if registeredStatus {
					log.Printf("\nSurvey status: %v", ls.DuringSurvey)
					if ls.DuringSurvey {
						err = ls.Survey.RecordAnswer(userId, message.Text)
						if err != nil {
							log.Print(err)
							aelog.Infof(ls.Ctx, "err: %s\n", err)
							break
						}

						question, err := ls.Survey.GetNextQuestion(userId)
						if err != nil {
							// ls.DuringSurvey = false
							log.Print(err)
							aelog.Infof(ls.Ctx, "err: %s\n", err)
							break
						}
						replyMsg = question.Text
					} else {
						switch message.Text {
						case "start":
							ls.DuringSurvey = true
							replyMsg, err = ls.startSurvey(userId)
							if err != nil {
								log.Print(err)
								aelog.Infof(ls.Ctx, "err: %s\n", err)
								break
							}
						default:
							replyMsg = `Type "start" to start filling survey`
						}
					}
				} else {
					replyMsg = `You haven't finished telling us about your profile yet
						Please tell us about yourself at https://test.com before proceeding to the next step 
						Otherwise you can complete the form here by typing "register" (3 hearts)`
				}
				_, err = ls.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do()
				ls.checkErr(err)
			}
		}
	}
}
