package api

import(
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"pt-report-backend/survey-responses"
	"log"
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

func (r *surveyResRoute) getSurveyResponses(c echo.Context) error {
	

	var params struct {
		CustomerCode string `json:"customerCode"`
	}
	if err := c.Bind(&params); err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if strings.TrimSpace(params.CustomerCode) != "" {
       	err := r.SurveyResService.Responses(params.CustomerCode)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		log.Println("Import response Success")
    }
	 
	return c.JSON(http.StatusOK, nil)
}
