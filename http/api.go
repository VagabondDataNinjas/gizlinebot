package http

import (
	"fmt"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/VagabondDataNinjas/gizlinebot/storage"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Api struct {
	Storage storage.Storage
	LineBot *linebot.Client
	Port    int
}

func NewApi(port int, s storage.Storage, lb *linebot.Client) *Api {
	return &Api{
		Storage: s,
		Port:    port,
		LineBot: lb,
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

	e.GET("/api/webform/questions", func(c echo.Context) error {
		return c.JSON(http.StatusOK, qs)
	})

	e.POST("/api/webform/answer", AnswerHandlerBuilder(a.Storage))

	// @TODO add authentication
	e.POST("/api/admin/send-msg", SendLineMsgHandlerBuilder(a.Storage, a.LineBot))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", a.Port)))
	return nil
}
