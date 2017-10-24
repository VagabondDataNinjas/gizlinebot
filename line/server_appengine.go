package line

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type LineAppEngine struct {
	BotHandler *httphandler.WebhookHandler
	Storage    storage.Storage
	Survey     *survey.Survey
}

func NewLineAppEngine(storage storage.Storage, surv *survey.Survey, secret, token string) (server *LineAppEngine, err error) {
	handler, err := httphandler.New(secret, token)
	if err != nil {
		return server, err
	}

	return &LineAppEngine{
		BotHandler: handler,
		Storage:    storage,
		Survey:     surv,
	}, nil
}

func (ls *LineAppEngine) ServeAppEngine() {
	http.HandleFunc("/linewebhook", ls.lineBotHandler)
	http.HandleFunc("/", ls.testHandler)
}

func (ls *LineAppEngine) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func (ls *LineAppEngine) lineBotHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	bot, err := ls.BotHandler.NewClient(linebot.WithHTTPClient(urlfetch.Client(ctx)))
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	// aelog.Infof(ctx, "\nevents: %+v", events)

	for _, event := range events {
		userId := event.Source.UserID
		eventString, err := json.Marshal(event)
		if err != nil {
			aelog.Errorf(ctx, "[err] Could not marshal event: %+v; err: %s", event, err)
			return
		} else {
			aelog.Infof(ctx, "\nEvent string: %s\n", eventString)
			err = ls.Storage.AddRawLineEvent(string(event.Type), string(eventString))
			if err != nil {
				aelog.Errorf(ctx, "[err] Could not store event: %+v; err: %s", event, err)
				return
			}
		}

		if event.Type == linebot.EventTypeFollow {
			userProfileResp, err := bot.GetProfile(userId).Do()
			if err != nil {
				log.Print(err)
				aelog.Warningf(ctx, "err: %s\n", err)
				continue
			}

			err = ls.Storage.AddUserProfile(userProfileResp.UserID, userProfileResp.DisplayName)
			if err != nil {
				aelog.Warningf(ctx, "AddUserProfile err: %s\n", err)
				continue
			}
		}

		if event.Type == linebot.EventTypeMessage {
			var replyMsg string
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				registeredStatus, err := ls.Storage.GetUserRegisteredStatus(userId)
				if err != nil {
					aelog.Warningf(ctx, "err: %s\n", err)
					break
				}
				if registeredStatus {
					err = ls.Survey.RecordAnswer(userId, message.Text)
					if err != nil {
						aelog.Warningf(ctx, "err: %s\n", err)
						break
					}

					question, err := ls.Survey.GetNextQuestion(userId)
					if err != nil {
						aelog.Warningf(ctx, "err: %s\n", err)
						break
					}
					replyMsg = question.Text
					// if ds {

					// } else {
					// 	switch message.Text {
					// 	case "start":
					// 		ds = true
					// 		replyMsg, err = ls.startSurvey(userId)
					// 		if err != nil {
					// 			log.Print(err)
					// 			// aelog.Infof(ls.Ctx, "err: %s\n", err)
					// 			break
					// 		}
					// 	default:
					// 		replyMsg = `Type "start" to start filling survey`
					// 	}
					// }
				} else {
					replyMsg = "You haven't finished telling us about your profile yet\n" +
						"Please tell us about yourself at https://test.com before proceeding to the next step\n" +
						"Otherwise you can complete the form here by typing 'register'"
				}
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do()
				if err != nil {
					aelog.Errorf(ctx, "err: %s\n", err)
					return
				}
			}
		}
	}
}

func (ls *LineAppEngine) startSurvey(userId string) (replyMsg string, err error) {
	question, err := ls.Survey.GetNextQuestion(userId)
	if err != nil {
		return "", err
	}
	return question.Text, nil
}
