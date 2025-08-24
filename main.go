package main

import (
	"log"
	"pt-report-backend/api"
	"pt-report-backend/db"
	"pt-report-backend/invitation"
	"pt-report-backend/transaction"
	"pt-report-backend/survey"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {

	err := db.InitDB(db.Config{
		Username:     "pt-report",
		Password:     "6EuFmuTj3d83hgQ/",
		Server:       "35.186.157.84",
		DatabaseName: "pt_record",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer db.DB.Close()

	invitationService := invitation.NewService()
	transactionService := transaction.NewService()
	surveyService := survey.NewSurveyService()

	// err = invitationService.CreateIndex()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// datedu,err := invitationService.GetDuration()
	// if err != nil {
	// 	log.Printf("Preload failed: %v", err)
	// }
	// // startDate := "2025-01-01"
	// startDate := datedu.Start_date.Format("2006-01-02")
	// endDate := time.Now().Format("2006-01-02")

	// go func() {
	// 	_,err = transactionService.PreloadTransactionCache(startDate,endDate)
	// 	if err != nil {
	// 		log.Printf("Preload transactions failed: %v", err)
	// 	}

	// 	_, err = invitationService.PreloadInvitationsCache(startDate, endDate)
	// 	if err != nil {
	// 		log.Printf("Preload invitations failed: %v", err)
	// 	}

	// 	_, err = invitationService.PreloadCustomers(startDate, endDate)
	// 	if err != nil {
	// 		log.Printf("Preload customers failed: %v", err)
	// 	}
	// 	log.Println("All Preload Success")
	// }()

	// // เรียก scheduler preload ทุกตี 1
	// scheduleDailyPreload(invitationService, transactionService)

	e := echo.New()
	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.Use(middleware.CORS())
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins:     []string{"https://pt-report-fcccf.web.app"},
	// 	AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	// 	AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	// 	AllowCredentials: true,
	// }))
	e.Pre(middleware.RemoveTrailingSlash())

	apiClient := api.NewAPI(invitationService,transactionService,surveyService)
	apiClient.Group(e.Group(""))

	// เริ่มเซิร์ฟเวอร์
	e.Logger.Fatal(e.Start(":8080"))
}

// func scheduleDailyPreload(invitationService *invitation.Service, transactionService *transaction.Service) {
// 	go func() {
// 		for {
// 			now := time.Now()
// 			next1AM := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, now.Location())
// 			if now.After(next1AM) {
// 				next1AM = next1AM.Add(24 * time.Hour)
// 			}
// 			wait := time.Until(next1AM)
// 			log.Printf("Waiting %v until next preload at 01:00\n", wait)

// 			time.Sleep(wait)

// 			// เรียก preload ที่ต้องการ
// 			// startDate := "2025-01-01"
// 			datedu,err := invitationService.GetDuration()
// 			if err != nil {
// 				log.Printf("Preload failed: %v", err)
// 			}

// 			startDate := datedu.Start_date.Format("2006-01-02")
// 			endDate := time.Now().Format("2006-01-02")


// 			log.Println("Starting daily preload at 01:00...")
// 			if _, err := transactionService.PreloadTransactionCache(startDate, endDate); err != nil {
// 				log.Printf("Preload transactions failed: %v", err)
// 			}
// 			if _, err := invitationService.PreloadInvitationsCache(startDate, endDate); err != nil {
// 				log.Printf("Preload invitations failed: %v", err)
// 			}
// 			if _, err := invitationService.PreloadCustomers(startDate, endDate); err != nil {
// 				log.Printf("Preload customers failed: %v", err)
// 			}
// 			log.Println("Daily preload complete")

// 			time.Sleep(24 * time.Hour) // รออีก 24 ชั่วโมงถัดไป
// 		}
// 	}()
// }
