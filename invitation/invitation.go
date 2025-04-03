package invitation

// import (
//     "github.com/labstack/echo/v4"
//     "net/http"
// 	"pt-report-backend/db"
// )

type Invitation struct {
	EDR_id      string `db:"EDR_id"`
	Wallet_type string `db:"Wallet_type"`
}
// type Service struct {

// }

// func (s *Service) Sendinvitation(c echo.Context) error {

// 	startDate := c.QueryParam("start_date")
// 		endDate := c.QueryParam("end_date")
// 		dateType := c.QueryParam("date_type")
		
// 		// ตรวจสอบว่ามีค่าหรือไม่
// 		if startDate == "" || endDate == "" {
// 			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing date parameters"})
// 		}

// 		customer, err := GetAllCustomers(db.DB, startDate, endDate)

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			
// 		}
// 		invitations, err := GetAllInvitation(db.DB, startDate, endDate,dateType)
		
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			
// 		}

// 		occupationCount, walletTypeCount := FilterCustomers(customer, invitations)

// 		return c.JSON(http.StatusOK, map[string]interface{}{
// 			"occupation": occupationCount,
// 			"wallet_type": walletTypeCount,
// 		})
// }
