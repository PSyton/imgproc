package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"imgproc/internal/processing"
)

type imageProcessor interface {
	Process() (string, error)
}

func (s *Server) processImageHandler(c echo.Context) error {
	processor, err := s.createProcessor(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if processor == nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	result, err := processor.Process()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusCreated, result)
}

func (s *Server) createProcessor(c echo.Context) (proc imageProcessor, err error) {
	req := c.Request()
	ctype := req.Header.Get(echo.HeaderContentType)
	if req.ContentLength == 0 {
		proc, err = processing.NewURLProcessor(c.QueryParam("url"), s.tools)
	} else {
		switch {
		case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):
			proc, err = processing.NewJSONProcessor(req.Body, s.tools)
		case strings.HasPrefix(ctype, echo.MIMEMultipartForm):
			form, e := c.MultipartForm()
			if e != nil {
				err = e
				return
			}
			files, ok := form.File["image"]
			if ok && len(files) == 1 {
				proc, err = processing.NewFormProcessor(files[0], s.tools)
			} else {
				if !ok {
					err = fmt.Errorf("No image field n form")
				} else {
					err = fmt.Errorf("Only single file allowed")
				}
			}
		}
	}

	return
}
