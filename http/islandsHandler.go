package http

import (
	"log"
	"net/http"

	"github.com/VagabondDataNinjas/gizlinebot/storage"

	"github.com/labstack/echo"
)

func IslandsHandler(s * storage.Sql) func(c echo.Context) error {
	return func(c echo.Context) error {
		locs, err := s.GetLocations()
		if err != nil {
			log.Printf("Error fetching locations: %s", err)
			return err
		}

		type respType struct {
			English string `json:"english"`
			Thai    string `json:"thai"`
		}
		resp := make([]respType, 0)
		for _, l := range locs {
			resp = append(resp, respType{l.Name, l.NameThai})
		}

		return c.JSON(http.StatusOK, resp)
	}
}
