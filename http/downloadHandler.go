package http

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/VagabondDataNinjas/gizlinebot/storage"
	"github.com/labstack/echo"
)

func DownloadHandlerBuilder(s * storage.Sql) func(e echo.Context) error {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "text/csv")
		c.Response().Header().Set("Content-Description", "File Transfer")
		c.Response().Header().Set("Content-Disposition", "attachment; filename=data.csv")
		c.Response().WriteHeader(http.StatusOK)

		b := &bytes.Buffer{}
		csvW := csv.NewWriter(b)

		questions, err := s.GetQuestions()
		if err != nil {
			return err
		}

		// create the CSV head
		csvHead := []string{
			"UserId",
			"Timestamp",
			"Channel",
		}
		for i := 0; i < questions.Len(); i++ {
			q, err := questions.At(i)
			if err != nil {
				return err
			}
			if q.Id == "gps" || q.Id == "thank_you" {
				// skip gps - we'll add it later
				continue
			}
			csvHead = append(csvHead, strings.Title(q.Id))
		}
		// add GPS related headers
		csvGpsIndex := len(csvHead)
		csvHead = append(csvHead, []string{"Address", "Lat", "Lon"}...)
		csvW.Write(csvHead)

		answerData, err := s.GetUserAnswerData()
		if err != nil {
			log.Errorf("Error getting user data: %s", err)
			return err
		}
		gpsData, err := s.GetGpsAnswerData()
		if err != nil {
			log.Errorf("Error getting gps data: %s", err)
			return err
		}

		csvRow := make([]string, len(csvHead))
		rowIndex := 0
		for _, answer := range answerData {
			// insert the cell data at the correct index based on csv head row
			for index, headCell := range csvHead {
				if strings.ToLower(headCell) != strings.ToLower(answer.QuestionId) {
					continue
				}
				if csvRow[0] == "" {
					csvRow[0] = answer.UserId
					csvRow[1] = time.Unix(int64(answer.Timestamp), 0).UTC().Format(time.RFC1123)
					csvRow[2] = answer.Channel
				} else if csvRow[index] != "" {
					if rowIndex < len(gpsData) {
						// add gps data
						csvRow[csvGpsIndex] = gpsData[rowIndex].Address
						csvRow[csvGpsIndex+1] = strconv.FormatFloat(gpsData[rowIndex].Lat, 'f', -1, 64)
						csvRow[csvGpsIndex+2] = strconv.FormatFloat(gpsData[rowIndex].Lon, 'f', -1, 64)
					}
					err = respWFlushRow(csvRow, csvW, c, b)
					if err != nil {
						return err
					}
					csvRow = make([]string, len(csvHead))
					rowIndex++
				}
				csvRow[index] = answer.Answer
			}
		}
		return respWFlushRow(csvRow, csvW, c, b)
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
