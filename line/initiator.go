package line

import (
	"fmt"
	"log"
	"time"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Initiator struct {
	Bot     *linebot.Client
	Storage storage.Storage
	Survey  *survey.Survey
}

func NewInitiator(surv *survey.Survey, s storage.Storage, bot *linebot.Client) *Initiator {
	return &Initiator{
		Bot:     bot,
		Storage: s,
		Survey:  surv,
	}
}

func (i *Initiator) Monitor(delay int64, errc chan error) {
	for c := time.Tick(5 * time.Second); ; <-c {
		userIds, err := i.Storage.GetUsersWithoutAnswers(delay)
		if err != nil {
			// @TODO log
			fmt.Printf("\n\nError pooling for answers: %s", err)
			errc <- err
			continue
		}

		if len(userIds) == 0 {
			continue
		}

		for _, userId := range userIds {
			question, err := i.Survey.GetNextQuestion(userId)
			if err != nil {
				log.Print(err)
				continue
			}
			if _, err = i.Bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
				log.Print(err)
				continue
			}

			if err = i.Storage.MarkProfileBotSurveyInited(userId); err != nil {
				log.Print(err)
				continue
			}

		}
		fmt.Printf("\nUserIds: %+v", userIds)
	}
}
