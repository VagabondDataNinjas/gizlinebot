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

func loginHandler(s *storage.Sql) func(c echo.Context) error {
	return func(c echo.Context) error {
		inputUser := c.FormValue("user")
		inputPass := c.FormValue("pass")

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

		log.Infof("User: %s; Pass: %s", user, pass)
		log.Infof("Iuser: %s; iPass: %s", user, pass)

		if user == "" || pass == "" || secret == "" {
			return errors.New("Setup error: missing adminUser, adminUser, or jwtSecret from DB setup")
		}

		if inputUser == user && inputPass == pass {
			// Create token
			token := jwt.New(jwt.SigningMethodHS256)

			// Set claims
			claims := token.Claims.(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

			// Generate encoded token and send it as response.
			t, err := token.SignedString([]byte(secret))
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"token": t,
			})
		}

		return echo.ErrUnauthorized
	}
}
