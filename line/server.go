package line

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineServer struct {
	Port    string
	Bot     *linebot.Client
	Storage storage.Storage
	Survey  *survey.Survey
}

func NewLineServer(port string, surv *survey.Survey, storage storage.Storage, secret, token string) (server *LineServer, err error) {
	bot, err := linebot.New(secret, token)
	if err != nil {
		return server, err
	}

	return &LineServer{
		Port:    port,
		Bot:     bot,
		Storage: storage,
		Survey:  surv,
	}, nil
}

func (ls *LineServer) Serve() error {
	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/linewebhook", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("\nRequest: %+v\n", req)
		events, err := ls.Bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		fmt.Printf("\nevents: %+v", events)
		for _, event := range events {
			userId := event.Source.UserID
			eventString, err := json.Marshal(event)
			if err != nil {
				log.Printf("[err] Could not marshal event: %+v; err: %s", event, err)
			} else {
				log.Printf("\nEvent string: %s\n", eventString)
				err = ls.Storage.AddRawLineEvent(string(event.Type), string(eventString))
				if err != nil {
					log.Printf("[err] Could not store event: %+v; err: %s", event, err)
				}
			}

			if event.Type == linebot.EventTypeFollow {
				userProfileResp, err := ls.Bot.GetProfile(userId).Do()
				if err != nil {
					log.Print(err)
					continue
				}

				err = ls.Storage.AddUserProfile(userProfileResp.UserID, userProfileResp.DisplayName)
				if err != nil {
					fmt.Printf("AddUserProfile err: %s\n", err)
					continue
				}

				question, err := ls.Survey.GetNextQuestion(userId)
				if err != nil {
					log.Print(err)
					continue
				}
				if _, err = ls.Bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
					log.Print(err)
					continue
				}
			}

			if event.Type == linebot.EventTypeMessage {
				// fmt.Printf("\ntype: %q \ntoken: %q \nstring: %q", string(event.Type), event.ReplyToken, string(eventString))
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					err = ls.Survey.RecordAnswer(userId, message.Text)
					if err != nil {
						log.Print(err)
						break
					}

					question, err := ls.Survey.GetNextQuestion(userId)
					if err != nil {
						log.Print(err)
						break
					}
					fmt.Printf("\nUser [id: %s] said: %s", userId, message.Text)

					if _, err = ls.Bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	log.Printf("Starting http server on port %s", ls.Port)
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+ls.Port, nil); err != nil {
		log.Fatal(err)
	}

	return nil
}
