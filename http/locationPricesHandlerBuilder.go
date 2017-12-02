package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/labstack/echo"
)

func LocationPricesHandlerBuilder(s *storage.Sql) func(e echo.Context) error {
	return func(c echo.Context) error {
		resp, err := s.GetLocationPrices()
		if err != nil {
			log.Errorf("Error getting location prices: %s", err)
			return err
		}
		return c.JSON(http.StatusOK, resp)
	}
}
