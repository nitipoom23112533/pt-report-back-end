package api

import (
	"net/http"
	// "pt-report-backend/auth"
	"pt-report-backend/invitation"
	"pt-report-backend/transaction"

	"github.com/labstack/echo/v4"
)

type TransactionRoute struct {
	TransactionService *transaction.Service
	InvitationService *invitation.Service
	
}

func (r *TransactionRoute)Group(g *echo.Group)  {
	// g.Use(auth.Auth())
	g.GET( "/transaction", r.sendTransaction)
}

func (r *TransactionRoute)sendTransaction(c echo.Context) error {
	// รับ start_date และ end_date จาก Query Parameters
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
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

	// ดึงข้อมูล transaction
	transactions, err := r.TransactionService.GetAllTransaction(startDate, endDate)
	// transactions, err := r.TransactionService.GetCachedTransactions(startDate, endDate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// ประมวลผลข้อมูล
	// occupationCount, walletTypeCount := r.TransactionService.FilterCustomers(customer, transactions)
	occupationCount, walletTypeCount := r.InvitationService.FilterCustomers(customer, transactions,"transaction")


	// ส่ง response กลับ
	return c.JSON(http.StatusOK, map[string]interface{}{
		"occupation":  occupationCount,
		"wallet_type": walletTypeCount,
	})
}