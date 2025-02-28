package main

import (
	"log"
	"net/http"
	"pt-report-backend/auth"
	"pt-report-backend/db"
	"pt-report-backend/invitation"
	"pt-report-backend/transaction"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	err := db.InitDB(db.Config{
		Username:     "niti",
		Password:     "VO/tguyYm(HIKl0n",
		Server:       "35.186.157.84",
		DatabaseName: "pt_record",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer db.DB.Close()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	// ตัวอย่าง Route
	// e.Use(auth.Auth()) // ใช้งาน JWT Middleware

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Backend with Echo!")
	})

	e.GET("/invitation",func(c echo.Context) error {
		// รับ start_date และ end_date จาก Query Parameters
		startDate := c.QueryParam("start_date")
		endDate := c.QueryParam("end_date")
		dateType := c.QueryParam("date_type")

		// ตรวจสอบว่ามีค่าหรือไม่
		if startDate == "" || endDate == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing date parameters"})
		}

		customer, err := invitation.GetAllCustomers(db.DB, startDate, endDate)

		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			
		}
		invitations, err := invitation.GetAllInvitation(db.DB, startDate, endDate,dateType)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			
		}

		occupationCount, walletTypeCount := invitation.FilterCustomers(customer, invitations)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"occupation": occupationCount,
			"wallet_type": walletTypeCount,
		})
	},auth.Auth())

	e.GET("/transaction",func(c echo.Context) error {

		// รับ start_date และ end_date จาก Query Parameters
		startDate := c.QueryParam("start_date")
		endDate := c.QueryParam("end_date")

		// ตรวจสอบว่ามีค่าหรือไม่
		if startDate == "" || endDate == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing date parameters"})
		}

		customer, err := invitation.GetAllCustomers(db.DB, startDate, endDate)

		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			
		}
		transactions, err := transaction.GetAllTransaction(db.DB,startDate,endDate)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			
		}

		occupationCount, walletTypeCount := transaction.FilterCustomers(customer, transactions)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"occupation": occupationCount,
			"wallet_type": walletTypeCount,
		})
	},auth.Auth())

	e.GET("/login", auth.VerifyTokenHandler)
	// เริ่มเซิร์ฟเวอร์
	e.Logger.Fatal(e.Start(":8080"))
}
