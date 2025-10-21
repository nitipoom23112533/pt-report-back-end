package api

import (
	"pt-report-backend/invitation"
	"pt-report-backend/transaction"
	"pt-report-backend/survey"
	"pt-report-backend/survey-responses"
	"github.com/labstack/echo/v4"
	"pt-report-backend/auth"
)

type API struct {
	InvitationRoute *InvitationRoute
	TransactionRoute *TransactionRoute
	SurveyRoute *SurveyRoute
	SurveyResRoute *surveyResRoute
}

func NewAPI(invitationService *invitation.Service, transactionService *transaction.Service,surveyService *survey.SurveyService,SurveyResService *surveyresponses.SurveyResService) *API {
	return &API{
		InvitationRoute: &InvitationRoute{
			InvitationService: invitationService,
		},
		TransactionRoute: &TransactionRoute{
			TransactionService: transactionService,
			InvitationService:  invitationService,
		},
		SurveyRoute: &SurveyRoute{
			SurveyService: surveyService,
		},
		SurveyResRoute: &surveyResRoute{
			SurveyResService: SurveyResService,
		},
	}
}

func (api *API) Group(g *echo.Group)  {
	// ptReportGroup := g.Group("pt-report")
	ptReportGroup := g.Group("/pt-report")
	ptReportGroup.Use(auth.Auth()) // echojwt ของคุณ
	api.InvitationRoute.Group(ptReportGroup)
	api.TransactionRoute.Group(ptReportGroup)
	api.SurveyRoute.Group(ptReportGroup)

	ptReportPublic := g.Group("/pt-report")
	api.SurveyResRoute.Group(ptReportPublic)

}


