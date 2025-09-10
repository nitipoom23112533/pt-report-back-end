package api

import (
	"log"
	// "pt-report-backend/auth"
	"pt-report-backend/survey"

	"github.com/labstack/echo/v4"
	"net/http"
)

type SurveyRoute struct{
		
	SurveyService *survey.SurveyService

}

func (r *SurveyRoute)Group(g *echo.Group)  {
	// g.Use(auth.Auth())
	g.GET( "/surveyDb", r.getSurveyDb)
	g.PATCH("/updatesurveyDb", r.updateSurveyDb)
}

func (r *SurveyRoute)getSurveyDb(c echo.Context) error {

	survey,err := r.SurveyService.GetSurvey()
	if err != nil {
		log.Println(err)
		return c.JSON(500, err)
	}

	return c.JSON(http.StatusOK,survey)
}

func (r *SurveyRoute)updateSurveyDb(c echo.Context) error {
	var s survey.Survey
	if err := c.Bind(&s); err != nil {
		log.Println(err.Error())
		return c.JSON(500, err)
	}

	err := r.SurveyService.UpdateSurvey(&s)
	if err != nil {
		log.Println(err)
		return c.JSON(500, err)
	}
	return c.JSON(http.StatusOK,s)

}