package http

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/pkg/errors"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
)

type SendLineCustomQuestionRequest struct {
	QuestionId string `json:"question_id"`
	Text       string `json:"text"`
	ReplyText  string `json:"reply_text"`
}

func SendLineCustomQuestionHandlerBuilder(s *storage.Sql, lineBot *linebot.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		pyld := new(SendLineCustomQuestionRequest)

		if err := c.Bind(pyld); err != nil {
			return err
		}

		// add the message to the storage
		if err := s.AddCustomQuestion(pyld.QuestionId, pyld.Text, pyld.ReplyText); err != nil {
			return errors.Wrap(err, "Failed to add custom question")
		}

		// fetch all users from DB
		profiles, err := s.GetAllActiveUserProfiles()
		if err != nil {
			return errors.Wrap(err, "Failed to fetch user profiles")
		}

		warnings := make([]string, 0)
		for _, profile := range profiles {
			msgStr, err := templateCustomQuestion(pyld.Text, profile, s)
			if err != nil {
				warn := fmt.Sprintf("Got error for custom question template for userId \"%s\": %s", profile.UserId, err)
				log.Error(warn)
				warnings = append(warnings, warn)
				continue
			}

			lineMsg := linebot.NewTextMessage(msgStr)
			logfields := log.Fields{
				"DisplayName": profile.DisplayName,
				"UserId":      profile.UserId,
				"QuestionId":  pyld.QuestionId,
			}
			log.WithFields(logfields).Info("Sending custom question to line")
			if _, err := lineBot.PushMessage(profile.UserId, lineMsg).Do(); err != nil {
				warn := fmt.Sprintf("Got error when seding msg to %s: %s", profile.UserId, err)
				log.WithFields(logfields).Error(warn)
				warnings = append(warnings, warn)
				continue
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success", "warnings": warnings,
		})
	}
}

func templateCustomQuestion(msgTpl string, userProfile domain.UserProfile, s *storage.Sql) (msg string, err error) {
	loc, err := s.FindUserLocation(userProfile.UserId)
	if err != nil {
		return "", err
	}

	tpl, err := template.New("customQuestion").Parse(msgTpl)
	if err != nil {
		return "", err
	}

	type TplVars struct {
		Location    string
		DisplayName string
	}
	buf := new(bytes.Buffer)
	tplVars := TplVars{
		Location:    loc.NameThai,
		DisplayName: userProfile.DisplayName,
	}
	err = tpl.Execute(buf, tplVars)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
