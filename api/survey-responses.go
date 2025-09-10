package api

import(
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"pt-report-backend/survey-responses"
	"log"
	"io"
	"bytes"
	"strings"
	
)

type surveyResRoute struct {
	
SurveyResService *surveyresponses.SurveyResService
}

func (r *surveyResRoute) Group(g *echo.Group) {
	g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "sc5ng7VkXGcSx927TEGFFEvT6RBkq3fv", nil
		},
	}))
	g.POST("/api/responses", r.getSurveyResponses)
}

// func (r *surveyResRoute) getSurveyResponses(c echo.Context) error {
	

// 	var params struct {
// 		CustomerCode string `json:"customerCode"`
// 	}
// 	if err := c.Bind(&params); err != nil {
// 		log.Println(err.Error())
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}
// 	// log.Println(params.CustomerCode)
// 	err := r.SurveyResService.Responses(params.CustomerCode)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}
	 
// 	return c.JSON(http.StatusOK, nil)
// }

func (r *surveyResRoute) getSurveyResponses(c echo.Context) error {
    c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 1<<20)

    // log header
    log.Println("hdr content-type:", c.Request().Header.Get("Content-Type"))
    log.Println("hdr api-key len:", len(c.Request().Header.Get("api-key")))

    // read raw
    b, _ := io.ReadAll(c.Request().Body)
    log.Println("raw body:", string(b))
    // put back for Bind
    c.Request().Body = io.NopCloser(bytes.NewReader(b))

    var params struct {
        CustomerCode string `json:"customerCode"`
    }
    if err := c.Bind(&params); err != nil {
        log.Println("bind error:", err)
        return c.NoContent(http.StatusBadRequest)
    }
    if strings.TrimSpace(params.CustomerCode) == "" {
        log.Println("empty customerCode")
        return c.NoContent(http.StatusBadRequest)
    }

    if err := r.SurveyResService.Responses(params.CustomerCode); err != nil {
        log.Println("service error:", err)
        return c.NoContent(http.StatusInternalServerError)
    }
    return c.NoContent(http.StatusNoContent)
}
