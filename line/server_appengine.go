package line

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/smeeklai/gizlinebot/storage"
	"github.com/smeeklai/gizlinebot/survey"
)

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
		for _, event := range events {
			userId := event.Source.UserID
			eventString, err := json.Marshal(event)
			if err != nil {
				log.Printf("[err] Could not marshal event: %+v; err: %s", event, err)
				aelog.Infof(ctx, "[err] Could not marshal event: %+v; err: %s", event, err)
			} else {
				log.Printf("\nEvent string: %s\n", eventString)
				err = storage.AddRawLineEvent(string(event.Type), string(eventString))
				if err != nil {
					log.Printf("[err] Could not store event: %+v; err: %s", event, err)
					aelog.Infof(ctx, "[err] Could not store event: %+v; err: %s", event, err)
				}
			}

			if event.Type == linebot.EventTypeFollow {
				userProfileResp, err := bot.GetProfile(userId).Do()
				if err != nil {
					log.Print(err)
					aelog.Infof(ctx, "err: %s\n", err)
					continue
				}

				err = storage.AddUserProfile(userProfileResp.UserID, userProfileResp.DisplayName)
				if err != nil {
					fmt.Printf("AddUserProfile err: %s\n", err)
					aelog.Infof(ctx, "AddUserProfile err: %s\n", err)
					continue
				}

				question, err := surv.GetNextQuestion(userId)
				if err != nil {
					log.Print(err)
					aelog.Infof(ctx, "err: %s\n", err)
					continue
				}
				if _, err = bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
					log.Print(err)
					aelog.Infof(ctx, "err: %s\n", err)
					continue
				}
			}

			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					err = surv.RecordAnswer(userId, message.Text)
					if err != nil {
						log.Print(err)
						aelog.Infof(ctx, "err: %s\n", err)
						break
					}

					question, err := surv.GetNextQuestion(userId)
					if err != nil {
						log.Print(err)
						aelog.Infof(ctx, "err: %s\n", err)
						break
					}
					fmt.Printf("\nUser [id: %s] said: %s", userId, message.Text)

					if _, err = bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
						log.Print(err)
						aelog.Infof(ctx, "err: %s\n", err)
					}
				}
			}
		}
	})
	http.Handle("/linewebhook", handler)
	http.HandleFunc("/", testHandler)

	return nil
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
