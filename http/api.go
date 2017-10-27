package http

import (
	"fmt"
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/domain"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/VagabondDataNinjas/gizlinebot/survey"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type ApiConf struct {
	Port int
	// hostname where the API is hosted
	GlobalVars *domain.GlobalTplVars
}

type Api struct {
	Storage storage.Storage
	LineBot *linebot.Client
	Surv    *survey.Survey
	Conf    *ApiConf
}

func NewApi(s storage.Storage, lb *linebot.Client, surv *survey.Survey, conf *ApiConf) *Api {
	return &Api{
		Storage: s,
		LineBot: lb,
		Surv:    surv,
		Conf:    conf,
	}
}

func (a *Api) Serve() error {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	e.Use(middleware.CORS())

	qs, err := a.Storage.GetQuestions()
	if err != nil {
		return err
	}
	e.POST("/linewebhook", LineWebhookHandlerBuilder(a.Surv, a.Storage, a.LineBot, a.Conf.GlobalVars))

	e.GET("/api/webform/questions", func(c echo.Context) error {
		return c.JSON(http.StatusOK, qs)
	})

	e.Static("/", "../gizsurvey/build")

	e.POST("/api/webform/answer", AnswerHandlerBuilder(a.Storage))
	e.POST("/api/webform/answer-gps", AnswerGpsHandlerBuilder(a.Storage))

	// @TODO add authentication
	e.POST("/api/admin/send-msg", SendLineMsgHandlerBuilder(a.Storage, a.LineBot))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", a.Conf.Port)))
	return nil
}
