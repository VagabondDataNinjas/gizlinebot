package survey

import (
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

type Initiator struct {
	Bot        *linebot.Client
	Storage    *storage.Sql
	Survey     *Survey
	GlobalVars *domain.GlobalTplVars
}

func NewInitiator(surv *Survey, s *storage.Sql, bot *linebot.Client, globalVars *domain.GlobalTplVars) *Initiator {
	return &Initiator{
		Bot:        bot,
		Storage:    s,
		Survey:     surv,
		GlobalVars: globalVars,
	}
}

func (i *Initiator) Monitor(delay int64, errc chan error) {
	for c := time.Tick(5 * time.Second); ; <-c {
		userIds, err := i.Storage.UsersSurveyNotStarted(delay)
		if err != nil {
			log.Errorf("Error pooling for answers: %s", err)
			errc <- err
			continue
		}

		if len(userIds) == 0 {
			continue
		}

		for _, userId := range userIds {
			questionVars := &QuestionTemplateVars{
				UserId:   userId,
				Hostname: i.GlobalVars.Hostname,
			}
			question, err := i.Survey.GetNextQuestion(userId, questionVars)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("Sending linebot survey to user: %s", userId)
			if _, err = i.Bot.PushMessage(userId, linebot.NewTextMessage(question.Text)).Do(); err != nil {
				log.Error(err)
				continue
			}

			if err = i.Storage.ToggleUserSurvey(userId, true); err != nil {
				log.Error(err)
				continue
			}
		}
	}
}
