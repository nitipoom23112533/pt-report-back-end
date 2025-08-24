package api

import (
	"pt-report-backend/invitation"
	"pt-report-backend/transaction"
	"pt-report-backend/survey"

	"github.com/labstack/echo/v4"
)

type API struct {
	InvitationRoute *InvitationRoute
	TransactionRoute *TransactionRoute
	SurveyRoute *SurveyRoute
}

func NewAPI(invitationService *invitation.Service, transactionService *transaction.Service,surveyService *survey.SurveyService) *API {
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

	}
}

func (api *API) Group(g *echo.Group)  {
	ptReportGroup := g.Group("pt-report")
	api.InvitationRoute.Group(ptReportGroup)
	api.TransactionRoute.Group(ptReportGroup)
	api.SurveyRoute.Group(ptReportGroup)

}


