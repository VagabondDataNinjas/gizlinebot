package http

import (
	"bytes"
	"encoding/csv"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/labstack/echo"
)

func LineEventsDownloadHandlerBuilder(s *storage.Sql) func(e echo.Context) error {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "text/csv")
		c.Response().Header().Set("Content-Description", "File Transfer")
		c.Response().Header().Set("Content-Disposition", "attachment; filename=linevents.csv")
		c.Response().WriteHeader(http.StatusOK)

		// create the CSV head
		csvHead := []string{
			"UserId",
			"DisplayName",
			"EventType",
			"EventTime",
		}

		events, err := s.GetRawLineEvents()
		if err != nil {
			log.Errorf("Error getting raw events data: %s", err)
			return err
		}
		b := &bytes.Buffer{}
		csvW := csv.NewWriter(b)
		csvRow := make([]string, len(csvHead))
		for _, event := range events {
			// insert the cell data at the correct index based on csv head row
			csvRow[0] = event.UserId
			csvRow[1] = event.DisplayName
			csvRow[2] = event.EventType
			csvRow[3] = event.EventTime
			if err = respWFlushRow(csvRow, csvW, c, b); err != nil {
				return err
			}
		}
		return nil
	}
}

func respWFlushRow(csvRow []string, csvW *csv.Writer, c echo.Context, b *bytes.Buffer) error {
	if err := csvW.Write(csvRow); err != nil {
		return err
	}
	csvW.Flush()
	if _, err := c.Response().Write(b.Bytes()); err != nil {
		return err
	}

	c.Response().Flush()
	b.Reset()
	return nil
}
