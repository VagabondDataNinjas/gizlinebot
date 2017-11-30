package http

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	logrusmiddleware "github.com/andreiashu/echo-logrusmiddleware"
	log "github.com/sirupsen/logrus"

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
	Storage *storage.Sql
	LineBot *linebot.Client
	Surv    *survey.Survey
	Conf    *ApiConf
}

func NewApi(s *storage.Sql, lb *linebot.Client, surv *survey.Survey, conf *ApiConf) *Api {
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
	e.Logger = logrusmiddleware.Logger{log.StandardLogger()}
	e.Use(logrusmiddleware.Hook())
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

	// lineimgs/groots.png/700
	e.GET("/lineimgs/:filename", LineImgHandler(), NoCacheMW)
	e.GET("/lineimgs/:filename/:size", LineImgHandler(), NoCacheMW)

	e.GET("/api/islands", IslandsHandler(a.Storage), NoCacheMW)

	e.POST("/api/webform/answer", AnswerHandlerBuilder(a.Storage, a.LineBot))
	e.POST("/api/webform/answer-gps", AnswerGpsHandlerBuilder(a.Storage))

	// @TODO add authentication
	e.POST("/api/admin/send-msg", SendLineMsgHandlerBuilder(a.Storage, a.LineBot))
	e.POST("/api/admin/send/custom/question", SendLineCustomQuestionHandlerBuilder(a.Storage, a.LineBot))

	e.GET("/api/admin/download/data.csv", DownloadHandlerBuilder(a.Storage), NoCacheMW)
	e.GET("/api/admin/download/lineevents.csv", LineEventsDownloadHandlerBuilder(a.Storage), NoCacheMW)

	log.Fatal(e.Start(fmt.Sprintf(":%d", a.Conf.Port)))
	return nil
}

func PriceHandler(s *storage.Sql) func(c echo.Context) error {
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

func userPriceList(userId string, s *storage.Sql) (priceList string, err error) {
	lp, err := s.GetUserNearbyPrices(userId)
	if err != nil {
		return "", err
	}

	if len(lp) == 0 {
		log.Infof("No prices found for user %s", userId)
		return "", nil
	}

	tplStr, err := s.GetPriceTplMsg()
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("priceMsg").Parse(tplStr)
	if err != nil {
		return "", err
	}

	type TplVars struct {
		Location string
		Price    float64
	}

	var msg string
	for _, loc := range lp {
		buf := new(bytes.Buffer)
		tplVars := TplVars{loc.Name, loc.Price}
		err = tmpl.Execute(buf, tplVars)
		if err != nil {
			return "", err
		}

		msg += buf.String()
		msg += "\n"
	}

	return msg, nil
}

func sendPriceList(bot *linebot.Client, s *storage.Sql, userId string) error {
	priceList, err := userPriceList(userId, s)
	if err != nil {
		return err
	}
	if priceList == "" {
		return nil
	}

	if _, err = bot.PushMessage(userId, linebot.NewTextMessage(priceList)).Do(); err != nil {
		return err
	}
	return nil
}

func sendThankYouMsg(bot *linebot.Client, s *storage.Sql, userId string) error {
	questions, err := s.GetQuestions()
	if err != nil {
		log.Errorf("Error getting questions: %s", err)
		return err
	}

	lastQ, err := questions.Last()
	if err != nil {
		log.Errorf("Error getting last question: %s", err)
		return err
	}

	if _, err = bot.PushMessage(userId, linebot.NewTextMessage(lastQ.Text)).Do(); err != nil {
		log.Errorf("Error sending last thankyou / question: %s", err)
		return err
	}
	return nil
}
