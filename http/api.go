package http

import (
	"fmt"
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/storage"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Api struct {
	Storage storage.Storage
	Port    int
}

func NewApi(port int, s storage.Storage) *Api {
	return &Api{
		Storage: s,
		Port:    port,
	}
}

func (a *Api) Serve() error {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	qs, err := a.Storage.GetQuestions()
	if err != nil {
		return err
	}

	e.GET("/api/webform/questions", func(c echo.Context) error {
		return c.JSON(http.StatusOK, qs)
	})

	e.POST("/api/webform/answer", AnswerHandlerBuilder(a.Storage))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", a.Port)))
	return nil
}
