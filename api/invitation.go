package api

import (
	"net/http"
	"pt-report-backend/invitation"
	"github.com/labstack/echo/v4"
	"pt-report-backend/auth"

)

type InvitationRoute struct {
	InvitationService *invitation.Service
}
func (r *InvitationRoute) Group (g *echo.Group)  {
	g.Use(auth.Auth())
	g.GET("/invitation", r.sendinvitation)
}

func (r *InvitationRoute) sendinvitation(c echo.Context) error {

	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	dateType := c.QueryParam("date_type")
	selected1InvPProfile := c.QueryParam("selected1InvPProfile")
	selectedAllProfile := c.QueryParam("selectedAllProfile")


	// ตรวจสอบว่ามีค่าหรือไม่
	if startDate == "" || endDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing date parameters"})
	}

	customer, err := r.InvitationService.GetAllCustomers(startDate, endDate,selectedAllProfile)
	// customer, err := r.InvitationService.GetCachedCustomers(startDate, endDate,selectedAllProfile)


	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}
	invitations, err := r.InvitationService.GetAllInvitation(startDate, endDate, dateType,selected1InvPProfile)
	// invitations, err := r.InvitationService.GetCachedInvitations(startDate,endDate,dateType,selected1InvPProfile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}

	occupationCount, walletTypeCount := r.InvitationService.FilterCustomers(customer, invitations,"invitation")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"occupation":  occupationCount,
		"wallet_type": walletTypeCount,
	})
}
