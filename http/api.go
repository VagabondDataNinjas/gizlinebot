package http

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

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

// cache-busting middleware
func NoCacheMW(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		return next(c)
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

	e.Group("/", NoCacheMW).Static("", "../gizsurvey/build")
	e.Group("/static", NoCacheMW).Static("", "../gizsurvey/build/static")
	e.Group("/media", NoCacheMW).Static("", "../media")

	e.GET("/api/user/wipe/:userid", WipeUserHandlerBuilder(a.Storage, a.LineBot), NoCacheMW)

	e.GET("/api/user/prices/:userid", PriceHandler(a.Storage), NoCacheMW) // @TODO remove

	e.GET("/api/islands", IslandsHandler(a.Storage), NoCacheMW)

	e.POST("/api/webform/answer", AnswerHandlerBuilder(a.Storage, a.LineBot))
	e.POST("/api/webform/answer-gps", AnswerGpsHandlerBuilder(a.Storage))

	// @TODO add authentication
	e.POST("/api/admin/send-msg", SendLineMsgHandlerBuilder(a.Storage, a.LineBot))

	e.GET("/api/admin/download/data.csv", DownloadHandlerBuilder(a.Storage), NoCacheMW)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", a.Conf.Port)))
	return nil
}

func PriceHandler(s storage.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		userId := c.Param("userid")
		lp, err := s.GetUserNearbyPrices(userId)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success", "locs": lp,
		})
	}
}

func sendPriceList(bot *linebot.Client, s storage.Storage, userId string) error {
	lp, err := s.GetUserNearbyPrices(userId)
	if err != nil {
		return err
	}

	tplStr, err := s.GetPriceTplMsg()
	if err != nil {
		return err
	}
	tmpl, err := template.New("priceMsg").Parse(tplStr)
	if err != nil {
		return err
	}

	type TplVars struct {
		Location string
		Price    float64
	}

	// @TODO check empty list of lp
	var msg string
	for _, loc := range lp {
		buf := new(bytes.Buffer)
		tplVars := TplVars{loc.Name, loc.Price}
		err = tmpl.Execute(buf, tplVars)
		if err != nil {
			return err
		}

		msg += buf.String()
		msg += "\n"
	}

	if _, err = bot.PushMessage(userId, linebot.NewTextMessage(msg)).Do(); err != nil {
		return err
	}

	return nil
}
