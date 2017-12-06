package http

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type LoginPayload struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func loginHandler(s *storage.Sql) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(LoginPayload)

		if err := c.Bind(payload); err != nil {
			log.Errorf("[loginHandler] Bind error: %s", err)
			return err
		}

		user, err := s.GetConfigVal("adminUser")
		if err != nil {
			return err
		}
		pass, err := s.GetConfigVal("adminPass")
		if err != nil {
			return err
		}
		secret, err := s.GetConfigVal("jwtSecret")
		if err != nil {
			return err
		}

		if user == "" || pass == "" || secret == "" {
			return errors.New("Setup error: missing adminUser, adminUser, or jwtSecret from DB setup")
		}

		if payload.User == user && payload.Pass == pass {
			// Create token
			token := jwt.New(jwt.SigningMethodHS256)

			// Set claims
			claims := token.Claims.(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

			// Generate encoded token
			t, err := token.SignedString([]byte(secret))
			if err != nil {
				return err
			}

			cookie := new(http.Cookie)
			cookie.Name = "token"
			cookie.Value = t
			cookie.Expires = time.Now().Add(87600 * time.Hour)
			c.SetCookie(cookie)

			return c.JSON(http.StatusOK, map[string]string{
				"status": "success",
			})
		}

		return echo.ErrUnauthorized
	}
}
