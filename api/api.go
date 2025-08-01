package api

import (
	"pt-report-backend/invitation"
	"pt-report-backend/transaction"

	"github.com/labstack/echo/v4"
)

type API struct {
	InvitationRoute *InvitationRoute
	TransactionRoute *TransactionRoute
}

// func NewAPI() *API {
// 	InvitationService := invitation.NewService()
// 	TransactionService := transaction.NewService()
	
// 	return &API{
// 		InvitationRoute: &InvitationRoute{
// 			InvitationService: InvitationService,
// 		},
// 		TransactionRoute: &TransactionRoute{
// 			TransactionService: TransactionService,
// 			InvitationService: InvitationService,

// 		},

// 	}
// }

func NewAPI(invitationService *invitation.Service, transactionService *transaction.Service) *API {
	return &API{
		InvitationRoute: &InvitationRoute{
			InvitationService: invitationService,
		},
		TransactionRoute: &TransactionRoute{
			TransactionService: transactionService,
			InvitationService:  invitationService,
		},
	}
}

func (api *API) Group(g *echo.Group)  {
	ptReportGroup := g.Group("pt-report")
	api.InvitationRoute.Group(ptReportGroup)
	api.TransactionRoute.Group(ptReportGroup)

	// g.POST("/refresh-cache", api.refreshCache)
}

// func (api *API) refreshCache(c echo.Context) error {

// }

