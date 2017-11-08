package http

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func LineImgHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		filename := c.Param("filename")
		size, err := strconv.Atoi(c.Param("size"))
		if err != nil {
			log.Error(err)
			return err
		}

		ext := path.Ext(filename)
		if ext == "" || (ext != ".png" && ext != ".jpg" && ext != ".jpeg") {
			return errors.Errorf("Invalid extension '%s' in filename '%s'", ext, filename)
		}

		filepath := "../line_imgs/" + filename
		reader, err := os.Open(filepath)
		if err != nil {
			log.Error(err)
			return err
		}
		defer reader.Close()
		image, _, err := image.Decode(reader)
		if err != nil {
			log.Errorf("Error trying to decode image file %s: %s", filepath, err)
			return err
		}

		newImg := imaging.Resize(image, size, 0, imaging.Lanczos)
		buff := new(bytes.Buffer)
		var contentType string
		if ext == ".png" {
			contentType = "image/png"
			err = png.Encode(buff, newImg)
			if err != nil {
				log.Errorf("Error trying to encode '%s' to png: %s", filepath, err)
				return err
			}
		} else {
			contentType = "image/jpeg"
			err = jpeg.Encode(buff, newImg, nil)
			if err != nil {
				log.Errorf("Error trying to encode '%s' to jpeg: %s", filepath, err)
				return err
			}
		}
		newImgReader := bytes.NewReader(buff.Bytes())
		return c.Stream(http.StatusOK, contentType, newImgReader)
	}
}
