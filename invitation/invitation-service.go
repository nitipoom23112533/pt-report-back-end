package invitation

import (
	"log"
	"github.com/jmoiron/sqlx"
    "fmt"
    "sync"
    "time"
    "pt-report-backend/db"
)
type Service struct {
    CacheMutex  sync.RWMutex
    Cache      []Invitation
    CacheCustomer      []Customer
    
}

func NewService() *Service {
	return &Service{}
}
func (s *Service) GetDuration() (Duration, error) {
    query := `SELECT start_date, end_date FROM duration WHERE id = 1;`
    var du Duration
    err := db.DB.Get(&du, query)
    if err != nil {
        log.Printf("Error fetching duration: %v", err)
        return du, err
    }

    return du, nil
}

func (s *Service)GetAllCustomers(db *sqlx.DB, startDate, endDate string,selectedAllProfile string) ([]Customer, error) {
    
    var customers []Customer
    if selectedAllProfile == "1" {
        query := `SELECT Customer_code,
	                MAX(Occupation) AS Occupation,
                    MAX(Customer_segment) AS Customer_segment,
                    MAX(Usage_segment) AS Usage_segment,
                    MAX(Age_range) AS Age_range,
                    MAX(Gender) AS Gender
	            FROM pt_customer
	            GROUP BY Customer_code;` // ดึงข้อมูลทั้งหมด
                
        err := db.Select(&customers, query) // ดึงข้อมูลทั้งหมดใส่ slice
        if err != nil {
            log.Printf("Error fetching customers: %v", err)
            return nil, err
        }
    }else{
        query := `SELECT Customer_code,
	                MAX(Occupation) AS Occupation,
                    MAX(Customer_segment) AS Customer_segment,
                    MAX(Usage_segment) AS Usage_segment,
                    MAX(Age_range) AS Age_range,
                    MAX(Gender) AS Gender
	            FROM pt_customer where Customer_date between ? and ?
	            GROUP BY Customer_code;` // ดึงข้อมูลทั้งหมด

        err := db.Select(&customers, query, startDate, endDate) // ดึงข้อมูลทั้งหมดใส่ slice
        if err != nil {
            log.Printf("Error fetching customers: %v", err)
            return nil, err
        }
    }
	return customers, nil
}
func (s *Service)GetAllInvitation(db *sqlx.DB, startDate, endDate string,dateType string,selected1InvPProfile string)([]Invitation,error)  {
    // Choose the correct date column based on dateType
    var dateColumn string
    if dateType == "invitationDate" {
        dateColumn = "IN_date"
    } else {
       dateColumn = "T_date"
    }
    var grpb string
    var Wallet_type string
    if selected1InvPProfile == "1" {
        grpb = "GROUP BY EDR_id"
        Wallet_type = "MAX(Wallet_type) as Wallet_type"
    }else{
        grpb = ""
        Wallet_type = "Wallet_type"
    }
    query := fmt.Sprintf(`SELECT EDR_id, %s
                            FROM pt_invitation 
                            WHERE %s BETWEEN ? AND ?
                            %s`, Wallet_type,dateColumn,grpb)
    var invitation []Invitation
    err := db.Select(&invitation,query, startDate, endDate)
    if err != nil {
        log.Printf("Error fetching invitation: %v", err)
		return nil, err
    }
    return invitation,err
}
func (s *Service)PreloadCustomers(startDate, endDate string) ([]Customer, error) {
    
    query := `SELECT Customer_code,
	                MAX(Occupation) AS Occupation,
                    MAX(Customer_segment) AS Customer_segment,
                    MAX(Usage_segment) AS Usage_segment,
                    MAX(Age_range) AS Age_range,
                    MAX(Gender) AS Gender,
                    MAX(Customer_date) AS Customer_date
	            FROM pt_customer
	            GROUP BY Customer_code ;` // ดึงข้อมูลทั้งหมด
    // query := `SELECT Customer_code,
	//                 MAX(Occupation) AS Occupation,
    //                 MAX(Customer_segment) AS Customer_segment,
    //                 MAX(Usage_segment) AS Usage_segment,
    //                 MAX(Age_range) AS Age_range,
    //                 MAX(Gender) AS Gender,
    //                 MAX(Customer_date) AS Customer_date
	//             FROM pt_customer where Customer_date between ? and ?
	//             GROUP BY Customer_code LIMIT 10;` // ดึงข้อมูลทั้งหมด

    var customers []Customer
    // err := db.DB.Select(&customers, query,startDate,endDate) // ดึงข้อมูลทั้งหมดใส่ slice
    err := db.DB.Select(&customers, query) // ดึงข้อมูลทั้งหมดใส่ slice

    if err != nil {
        log.Printf("Error fetching customers: %v", err)
        return nil, err
    }
    s.CacheMutex.Lock()
	s.CacheCustomer = customers
	s.CacheMutex.Unlock()

	log.Println("Customers cache refreshed")
	return customers, nil
}
// preload ข้อมูลจาก DB
func (s *Service) PreloadInvitationsCache(startDate, endDate string) ([]Invitation, error) {

    // query := `
	// 	SELECT EDR_id, Wallet_type,IN_date,T_date
	// 	FROM pt_invitation
	// 	WHERE T_date BETWEEN ? AND ? LIMIT 10
		// `
    query := `
		SELECT EDR_id, Wallet_type,IN_date,T_date
		FROM pt_invitation
		WHERE T_date BETWEEN ? AND ?
		`
	var invitations []Invitation
	err := db.DB.Select(&invitations, query, startDate, endDate)
	if err != nil {
		log.Printf("Error fetching invitation: %v", err)
		return nil, err
	}

	s.CacheMutex.Lock()
	s.Cache = invitations
	s.CacheMutex.Unlock()

	log.Println("Invitation cache refreshed")
	return invitations, nil
}

func (s *Service) GetCachedInvitations(startDate string, endDate string, dateType string, selected1InvPProfile string) ([]Invitation,error) {
    log.Println("GetCachedInvitations called")

    dateOnlyLayout := "2006-01-02"
    datetimeLayout := "2006-01-02T15:04:05-07:00"
    start, err := time.Parse(dateOnlyLayout, startDate)
    if err != nil {
        return nil, err
    }
    end, err := time.Parse(dateOnlyLayout, endDate)
    if err != nil {
        return nil, err
    }

    s.CacheMutex.RLock()
    defer s.CacheMutex.RUnlock()

    var filtered []Invitation
    mapEDRid := make(map[string]bool) // ใช้สำหรับเช็คว่า EDRid นี้เคยถูกใส่ไปหรือยัง
    if dateType == "invitationDate" {
        for _, inv := range s.Cache {
            invDate, err := time.Parse(datetimeLayout, inv.IN_date) // สมมติเก็บวันที่ใน field DateStr เป็น string
            if err != nil {
                log.Println(err)
            }
        
            if (invDate.Equal(start) || invDate.After(start)) && (invDate.Equal(end) || invDate.Before(end)) {
                if selected1InvPProfile == "1" {
                    if !mapEDRid[inv.EDR_id] {
                        mapEDRid[inv.EDR_id] = true // mark ว่าเจอแล้ว
                        filtered = append(filtered, inv)
                    }
                }else{
                    filtered = append(filtered, inv)
                }
                
            }
        }
    } else {
        for _, inv := range s.Cache {
            traDate, err := time.Parse(datetimeLayout, inv.T_date) // สมมติเก็บวันที่ใน field DateStr เป็น string
            if err != nil {
                log.Println(err)
            }
        
            if (traDate.Equal(start) || traDate.After(start)) && (traDate.Equal(end) || traDate.Before(end)) {
                // filtered = append(filtered, inv)
                if selected1InvPProfile == "1" {
                    if !mapEDRid[inv.EDR_id] {
                        mapEDRid[inv.EDR_id] = true // mark ว่าเจอแล้ว
                        filtered = append(filtered, inv)
                    }
                }else{
                    filtered = append(filtered, inv)
                }
            }
        }
    }
    
    return filtered, nil
}

func (s *Service) GetCachedCustomers(startDate string, endDate string,selectedAllProfile string) ([]Customer,error) {
    log.Println("GetCachedCustomers called")
	// s.CacheMutex.RLock()
	// defer s.CacheMutex.RUnlock()
	// return s.Cache,nil
    dateOnlyLayout := "2006-01-02"
    datetimeLayout := "2006-01-02T15:04:05-07:00"
    start, err := time.Parse(dateOnlyLayout, startDate)
    if err != nil {
        return nil, err
    }
    end, err := time.Parse(dateOnlyLayout, endDate)
    if err != nil {
        return nil, err
    }

    s.CacheMutex.RLock()
    defer s.CacheMutex.RUnlock()

    var filtered []Customer
    if selectedAllProfile == "1" {

        return s.CacheCustomer,nil

    } else {
        for _, cus := range s.CacheCustomer {
            traDate, err := time.Parse(datetimeLayout, cus.Customer_date) // สมมติเก็บวันที่ใน field DateStr เป็น string
            if err != nil {
                log.Println(err)
            }
        
            if (traDate.Equal(start) || traDate.After(start)) && (traDate.Equal(end) || traDate.Before(end)) {
                filtered = append(filtered, cus)
            }
        }
    }
    
    return filtered, nil
}


func (s *Service)FilterCustomers(customers []Customer, invitations []Invitation) (CountOccupation,Wallet_type) {
	var countOccupation CountOccupation
    var countWallet_type Wallet_type
	// สร้าง map เพื่อเก็บ Customer_code -> Usage_segment
    customerMap := make(map[string]Customer)
	for _, c := range customers {
		customerMap[c.CustomerCode] = c
	}
	// วนลูป invitations และตรวจสอบว่ามี Customer_code ตรงกันหรือไม่
	for _, inv := range invitations {
		customerAll, exists := customerMap[inv.EDR_id]

		if exists {
			// ถ้ามี ให้คำนวณ count ตาม Usage_segment
			switch customerAll.Usage_segment {
			case "1) Low":
				countOccupation.Low++

                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01L++
                case "PT02":
                    countWallet_type.PT02L++
                case "PT03":
                    countWallet_type.PT03L++
                case "PT05":
                    countWallet_type.PT05L++
                case "PT06":
                    countWallet_type.PT06L++
                case "PT08":
                    countWallet_type.PT08L++
                case "PT09":
                    countWallet_type.PT09L++
                case "PT10":
                    countWallet_type.PT10L++
                case "PT15":
                    countWallet_type.PT15L++
                case "PT16":
                    countWallet_type.PT16L++
                case "PT17":
                    countWallet_type.PT17L++
                case "PT18":
                    countWallet_type.PT18L++
                case "PT19":
                    countWallet_type.PT19L++
                }

			case "2) Medium":
				countOccupation.Medium++

                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01M++
                case "PT02":
                    countWallet_type.PT02M++
                case "PT03":
                    countWallet_type.PT03M++
                case "PT05":
                    countWallet_type.PT05M++
                case "PT06":
                    countWallet_type.PT06M++
                case "PT08":
                    countWallet_type.PT08M++
                case "PT09":
                    countWallet_type.PT09M++
                case "PT10":
                    countWallet_type.PT10M++
                case "PT15":
                    countWallet_type.PT15M++
                case "PT16":
                    countWallet_type.PT16M++
                case "PT17":
                    countWallet_type.PT17M++
                case "PT18":
                    countWallet_type.PT18M++
                case "PT19":
                    countWallet_type.PT19M++
                }
			case "3) High":
				countOccupation.High++

                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01H++
                case "PT02":
                    countWallet_type.PT02H++
                case "PT03":
                    countWallet_type.PT03H++
                case "PT05":
                    countWallet_type.PT05H++
                case "PT06":
                    countWallet_type.PT06H++
                case "PT08":
                    countWallet_type.PT08H++
                case "PT09":
                    countWallet_type.PT09H++
                case "PT10":
                    countWallet_type.PT10H++
                case "PT15":
                    countWallet_type.PT15H++
                case "PT16":
                    countWallet_type.PT16H++
                case "PT17":
                    countWallet_type.PT17H++
                case "PT18":
                    countWallet_type.PT18H++
                case "PT19":
                    countWallet_type.PT19H++
                }
			case "4) Login Only":
				countOccupation.Login_Only++

                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01LO++
                case "PT02":
                    countWallet_type.PT02LO++
                case "PT03":
                    countWallet_type.PT03LO++
                case "PT05":
                    countWallet_type.PT05LO++
                case "PT06":
                    countWallet_type.PT06LO++
                case "PT08":
                    countWallet_type.PT08LO++
                case "PT09":
                    countWallet_type.PT09LO++
                case "PT10":
                    countWallet_type.PT10LO++
                case "PT15":
                    countWallet_type.PT15LO++
                case "PT16":
                    countWallet_type.PT16LO++
                case "PT17":
                    countWallet_type.PT17LO++
                case "PT18":
                    countWallet_type.PT18LO++
                case "PT19":
                    countWallet_type.PT19LO++
                }
            case "5) Screen View Only (No Login)":
                countOccupation.Screen_View_Only++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01SVO++
                case "PT02":
                    countWallet_type.PT02SVO++
                case "PT03":
                    countWallet_type.PT03SVO++
                case "PT05":
                    countWallet_type.PT05SVO++
                case "PT06":
                    countWallet_type.PT06SVO++
                case "PT08":
                    countWallet_type.PT08SVO++
                case "PT09":
                    countWallet_type.PT09SVO++
                case "PT10":
                    countWallet_type.PT10SVO++
                case "PT15":
                    countWallet_type.PT15SVO++
                case "PT16":
                    countWallet_type.PT16SVO++
                case "PT17":
                    countWallet_type.PT17SVO++
                case "PT18":
                    countWallet_type.PT18SVO++
                case "PT19":
                    countWallet_type.PT19SVO++
                }

            case "6) Inactive":
                countOccupation.Inactive++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01I++
                case "PT02":
                    countWallet_type.PT02I++
                case "PT03":
                    countWallet_type.PT03I++
                case "PT05":
                    countWallet_type.PT05I++
                case "PT06":
                    countWallet_type.PT06I++
                case "PT08":
                    countWallet_type.PT08I++
                case "PT09":
                    countWallet_type.PT09I++
                case "PT10":
                    countWallet_type.PT10I++
                case "PT15":
                    countWallet_type.PT15I++
                case "PT16":
                    countWallet_type.PT16I++
                case "PT17":
                    countWallet_type.PT17I++
                case "PT18":
                    countWallet_type.PT18I++
                case "PT19":
                    countWallet_type.PT19I++
                }
            case "7) New User":
                countOccupation.New_User++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01NU++
                case "PT02":
                    countWallet_type.PT02NU++
                case "PT03":
                    countWallet_type.PT03NU++
                case "PT05":
                    countWallet_type.PT05NU++
                case "PT06":
                    countWallet_type.PT06NU++
                case "PT08":
                    countWallet_type.PT08NU++
                case "PT09":
                    countWallet_type.PT09NU++
                case "PT10":
                    countWallet_type.PT10NU++
                case "PT15":
                    countWallet_type.PT15NU++
                case "PT16":
                    countWallet_type.PT16NU++
                case "PT17":
                    countWallet_type.PT17NU++
                case "PT18":
                    countWallet_type.PT18NU++
                case "PT19":
                    countWallet_type.PT19NU++
                }
			default:
				countOccupation.NULL_Usage_segment++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01NULL++
                case "PT02":
                    countWallet_type.PT02NULL++
                case "PT03":
                    countWallet_type.PT03NULL++
                case "PT05":
                    countWallet_type.PT05NULL++
                case "PT06":
                    countWallet_type.PT06NULL++
                case "PT08":
                    countWallet_type.PT08NULL++
                case "PT09":
                    countWallet_type.PT09NULL++
                case "PT10":
                    countWallet_type.PT10NULL++
                case "PT15":
                    countWallet_type.PT15NULL++
                case "PT16":
                    countWallet_type.PT16NULL++
                case "PT17":
                    countWallet_type.PT17NULL++
                case "PT18":
                    countWallet_type.PT18NULL++
                case "PT19":
                    countWallet_type.PT19NULL++
                }
			}
            // Gender
            switch customerAll.Gender {
            case "F":
                countOccupation.Gender_F++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Female++
                case "PT02":
                    countWallet_type.PT02Female++
                case "PT03":
                    countWallet_type.PT03Female++
                case "PT05":
                    countWallet_type.PT05Female++
                case "PT06":
                    countWallet_type.PT06Female++
                case "PT08":
                    countWallet_type.PT08Female++
                case "PT09":
                    countWallet_type.PT09Female++
                case "PT10":
                    countWallet_type.PT10Female++
                case "PT15":
                    countWallet_type.PT15Female++
                case "PT16":
                    countWallet_type.PT16Female++
                case "PT17":
                    countWallet_type.PT17Female++
                case "PT18":
                    countWallet_type.PT18Female++
                case "PT19":
                    countWallet_type.PT19Female++
                }
            case "M":
                countOccupation.Gender_M++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Male++
                case "PT02":
                    countWallet_type.PT02Male++
                case "PT03":
                    countWallet_type.PT03Male++
                case "PT05":
                    countWallet_type.PT05Male++
                case "PT06":
                    countWallet_type.PT06Male++
                case "PT08":
                    countWallet_type.PT08Male++
                case "PT09":
                    countWallet_type.PT09Male++
                case "PT10":
                    countWallet_type.PT10Male++
                case "PT15":
                    countWallet_type.PT15Male++
                case "PT16":
                    countWallet_type.PT16Male++
                case "PT17":
                    countWallet_type.PT17Male++
                case "PT18":
                    countWallet_type.PT18Male++
                case "PT19":
                    countWallet_type.PT19Male++
                }
            default:
                countOccupation.NULL_Gender++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01GNULL++
                case "PT02":
                    countWallet_type.PT02GNULL++
                case "PT03":
                    countWallet_type.PT03GNULL++
                case "PT05":
                    countWallet_type.PT05GNULL++
                case "PT06":
                    countWallet_type.PT06GNULL++
                case "PT08":
                    countWallet_type.PT08GNULL++
                case "PT09":
                    countWallet_type.PT09GNULL++
                case "PT10":
                    countWallet_type.PT10GNULL++
                case "PT15":
                    countWallet_type.PT15GNULL++
                case "PT16":
                    countWallet_type.PT16GNULL++
                case "PT17":
                    countWallet_type.PT17GNULL++
                case "PT18":
                    countWallet_type.PT18GNULL++
                case "PT19":
                    countWallet_type.PT19GNULL++
                }
            }
            // age range
            switch customerAll.Age_range {
            case "01] ต่ำกว่า 22":
                countOccupation.Age_range_1++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Age1++
                case "PT02":
                    countWallet_type.PT02Age1++
                case "PT03":
                    countWallet_type.PT03Age1++
                case "PT05":
                    countWallet_type.PT05Age1++
                case "PT06":
                    countWallet_type.PT06Age1++
                case "PT08":
                    countWallet_type.PT08Age1++
                case "PT09":
                    countWallet_type.PT09Age1++
                case "PT10":
                    countWallet_type.PT10Age1++
                case "PT15":
                    countWallet_type.PT15Age1++
                case "PT16":
                    countWallet_type.PT16Age1++
                case "PT17":
                    countWallet_type.PT17Age1++
                case "PT18":
                    countWallet_type.PT18Age1++
                case "PT19":
                    countWallet_type.PT19Age1++
                }
            case "02] 22 - 25":
                countOccupation.Age_range_2++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Age2++
                case "PT02":
                    countWallet_type.PT02Age2++
                case "PT03":
                    countWallet_type.PT03Age2++
                case "PT05":
                    countWallet_type.PT05Age2++
                case "PT06":
                    countWallet_type.PT06Age2++
                case "PT08":
                    countWallet_type.PT08Age2++
                case "PT09":
                    countWallet_type.PT09Age2++
                case "PT10":
                    countWallet_type.PT10Age2++
                case "PT15":
                    countWallet_type.PT15Age2++
                case "PT16":
                    countWallet_type.PT16Age2++
                case "PT17":
                    countWallet_type.PT17Age2++
                case "PT18":
                    countWallet_type.PT18Age2++
                case "PT19":
                    countWallet_type.PT19Age2++
                }
            case "03] 26  - 30":
                countOccupation.Age_range_3++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Age3++
                case "PT02":
                    countWallet_type.PT02Age3++
                case "PT03":
                    countWallet_type.PT03Age3++
                case "PT05":
                    countWallet_type.PT05Age3++
                case "PT06":
                    countWallet_type.PT06Age3++
                case "PT08":
                    countWallet_type.PT08Age3++
                case "PT09":
                    countWallet_type.PT09Age3++
                case "PT10":
                    countWallet_type.PT10Age3++
                case "PT15":
                    countWallet_type.PT15Age3++
                case "PT16":
                    countWallet_type.PT16Age3++
                case "PT17":
                    countWallet_type.PT17Age3++
                case "PT18":
                    countWallet_type.PT18Age3++
                case "PT19":
                    countWallet_type.PT19Age3++
                }
            case "04] 31 - 40":
                countOccupation.Age_range_4++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Age4++
                case "PT02":
                    countWallet_type.PT02Age4++
                case "PT03":
                    countWallet_type.PT03Age4++
                case "PT05":
                    countWallet_type.PT05Age4++
                case "PT06":
                    countWallet_type.PT06Age4++
                case "PT08":
                    countWallet_type.PT08Age4++
                case "PT09":
                    countWallet_type.PT09Age4++
                case "PT10":
                    countWallet_type.PT10Age4++
                case "PT15":
                    countWallet_type.PT15Age4++
                case "PT16":
                    countWallet_type.PT16Age4++
                case "PT17":
                    countWallet_type.PT17Age4++
                case "PT18":
                    countWallet_type.PT18Age4++
                case "PT19":
                    countWallet_type.PT19Age4++
                }
            case "05] 41 - 45":
                countOccupation.Age_range_5++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Age5++
                case "PT02":
                    countWallet_type.PT02Age5++
                case "PT03":
                    countWallet_type.PT03Age5++
                case "PT05":
                    countWallet_type.PT05Age5++
                case "PT06":
                    countWallet_type.PT06Age5++
                case "PT08":
                    countWallet_type.PT08Age5++
                case "PT09":
                    countWallet_type.PT09Age5++
                case "PT10":
                    countWallet_type.PT10Age5++
                case "PT15":
                    countWallet_type.PT15Age5++
                case "PT16":
                    countWallet_type.PT16Age5++
                case "PT17":
                    countWallet_type.PT17Age5++
                case "PT18":
                    countWallet_type.PT18Age5++
                case "PT19":
                    countWallet_type.PT19Age5++
                }
            case "06] 46 - 50":
                countOccupation.Age_range_6++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Age6++
                case "PT02":
                    countWallet_type.PT02Age6++
                case "PT03":
                    countWallet_type.PT03Age6++
                case "PT05":
                    countWallet_type.PT05Age6++
                case "PT06":
                    countWallet_type.PT06Age6++
                case "PT08":
                    countWallet_type.PT08Age6++
                case "PT09":
                    countWallet_type.PT09Age6++
                case "PT10":
                    countWallet_type.PT10Age6++
                case "PT15":
                    countWallet_type.PT15Age6++
                case "PT16":
                    countWallet_type.PT16Age6++
                case "PT17":
                    countWallet_type.PT17Age6++
                case "PT18":
                    countWallet_type.PT18Age6++
                case "PT19":
                    countWallet_type.PT19Age6++
                }
            case "07] 51 - 60":
                countOccupation.Age_range_7++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Age7++
                case "PT02":
                    countWallet_type.PT02Age7++
                case "PT03":
                    countWallet_type.PT03Age7++
                case "PT05":
                    countWallet_type.PT05Age7++
                case "PT06":
                    countWallet_type.PT06Age7++
                case "PT08":
                    countWallet_type.PT08Age7++
                case "PT09":
                    countWallet_type.PT09Age7++
                case "PT10":
                    countWallet_type.PT10Age7++
                case "PT15":
                    countWallet_type.PT15Age7++
                case "PT16":
                    countWallet_type.PT16Age7++
                case "PT17":
                    countWallet_type.PT17Age7++
                case "PT18":
                    countWallet_type.PT18Age7++
                case "PT19":
                    countWallet_type.PT19Age7++
                }
            case "08] มากกว่า 60 ปี":
                countOccupation.Age_range_8++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Age8++
                case "PT02":
                    countWallet_type.PT02Age8++
                case "PT03":
                    countWallet_type.PT03Age8++
                case "PT05":
                    countWallet_type.PT05Age8++
                case "PT06":
                    countWallet_type.PT06Age8++
                case "PT08":
                    countWallet_type.PT08Age8++
                case "PT09":
                    countWallet_type.PT09Age8++
                case "PT10":
                    countWallet_type.PT10Age8++
                case "PT15":
                    countWallet_type.PT15Age8++
                case "PT16":
                    countWallet_type.PT16Age8++
                case "PT17":
                    countWallet_type.PT17Age8++
                case "PT18":
                    countWallet_type.PT18Age8++
                case "PT19":
                    countWallet_type.PT19Age8++
                }
            default:
                countOccupation.NULL_Age_range++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01AgeNULL++
                case "PT02":
                    countWallet_type.PT02AgeNULL++
                case "PT03":
                    countWallet_type.PT03AgeNULL++
                case "PT05":
                    countWallet_type.PT05AgeNULL++
                case "PT06":
                    countWallet_type.PT06AgeNULL++
                case "PT08":
                    countWallet_type.PT08AgeNULL++
                case "PT09":
                    countWallet_type.PT09AgeNULL++
                case "PT10":
                    countWallet_type.PT10AgeNULL++
                case "PT15":
                    countWallet_type.PT15AgeNULL++
                case "PT16":
                    countWallet_type.PT16AgeNULL++
                case "PT17":
                    countWallet_type.PT17AgeNULL++
                case "PT18":
                    countWallet_type.PT18AgeNULL++
                case "PT19":
                    countWallet_type.PT19AgeNULL++
                }
            }
            // CUstomer Segment
            switch customerAll.Customer_segment {
            case "02_PRECIOUSPLUS":
                countOccupation.PRECIOUSPLUS++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS02++
                case "PT02":
                    countWallet_type.PT02CS02++
                case "PT03":
                    countWallet_type.PT03CS02++
                case "PT05":
                    countWallet_type.PT05CS02++
                case "PT06":
                    countWallet_type.PT06CS02++
                case "PT08":
                    countWallet_type.PT08CS02++
                case "PT09":
                    countWallet_type.PT09CS02++
                case "PT10":
                    countWallet_type.PT10CS02++
                case "PT15":
                    countWallet_type.PT15CS02++
                case "PT16":
                    countWallet_type.PT16CS02++
                case "PT17":
                    countWallet_type.PT17CS02++
                case "PT18":
                    countWallet_type.PT18CS02++
                case "PT19":
                    countWallet_type.PT19CS02++
                }
            case "03_PRECIOUS":
                countOccupation.PRECIOUS++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS03++
                case "PT02":
                    countWallet_type.PT02CS03++
                case "PT03":
                    countWallet_type.PT03CS03++
                case "PT05":
                    countWallet_type.PT05CS03++
                case "PT06":
                    countWallet_type.PT06CS03++
                case "PT08":
                    countWallet_type.PT08CS03++
                case "PT09":
                    countWallet_type.PT09CS03++
                case "PT10":
                    countWallet_type.PT10CS03++
                case "PT15":
                    countWallet_type.PT15CS03++
                case "PT16":
                    countWallet_type.PT16CS03++
                case "PT17":
                    countWallet_type.PT17CS03++
                case "PT18":
                    countWallet_type.PT18CS03++
                case "PT19":
                    countWallet_type.PT19CS03++
                }
            case "04_PREWEALTH":
                countOccupation.PREWEALTH++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS04++
                case "PT02":
                    countWallet_type.PT02CS04++
                case "PT03":
                    countWallet_type.PT03CS04++
                case "PT05":
                    countWallet_type.PT05CS04++
                case "PT06":
                    countWallet_type.PT06CS04++
                case "PT08":
                    countWallet_type.PT08CS04++
                case "PT09":
                    countWallet_type.PT09CS04++
                case "PT10":
                    countWallet_type.PT10CS04++
                case "PT15":
                    countWallet_type.PT15CS04++
                case "PT16":
                    countWallet_type.PT16CS04++
                case "PT17":
                    countWallet_type.PT17CS04++
                case "PT18":
                    countWallet_type.PT18CS04++
                case "PT19":
                    countWallet_type.PT19CS04++
                }
            case "05_AFFUIENTTOBE":
                countOccupation.AFFUIENTTOBE++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS05++
                case "PT02":
                    countWallet_type.PT02CS05++
                case "PT03":
                    countWallet_type.PT03CS05++
                case "PT05":
                    countWallet_type.PT05CS05++
                case "PT06":
                    countWallet_type.PT06CS05++
                case "PT08":
                    countWallet_type.PT08CS05++
                case "PT09":
                    countWallet_type.PT09CS05++
                case "PT10":
                    countWallet_type.PT10CS05++
                case "PT15":
                    countWallet_type.PT15CS05++
                case "PT16":
                    countWallet_type.PT16CS05++
                case "PT17":
                    countWallet_type.PT17CS05++
                case "PT18":
                    countWallet_type.PT18CS05++
                case "PT19":
                    countWallet_type.PT19CS05++
                }
            case "06_RETIREPLANNER":
                countOccupation.RETIREPLANNER++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS06++
                case "PT02":
                    countWallet_type.PT02CS06++
                case "PT03":
                    countWallet_type.PT03CS06++
                case "PT05":
                    countWallet_type.PT05CS06++
                case "PT06":
                    countWallet_type.PT06CS06++
                case "PT08":
                    countWallet_type.PT08CS06++
                case "PT09":
                    countWallet_type.PT09CS06++
                case "PT10":
                    countWallet_type.PT10CS06++
                case "PT15":
                    countWallet_type.PT15CS06++
                case "PT16":
                    countWallet_type.PT16CS06++
                case "PT17":
                    countWallet_type.PT17CS06++
                case "PT18":
                    countWallet_type.PT18CS06++
                case "PT19":
                    countWallet_type.PT19CS06++
                }
            case "07_BUILDUPFORFEATURE":
                countOccupation.BUILDUPFORFEATURE++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS07++
                case "PT02":
                    countWallet_type.PT02CS07++
                case "PT03":
                    countWallet_type.PT03CS07++
                case "PT05":
                    countWallet_type.PT05CS07++
                case "PT06":
                    countWallet_type.PT06CS07++
                case "PT08":
                    countWallet_type.PT08CS07++
                case "PT09":
                    countWallet_type.PT09CS07++
                case "PT10":
                    countWallet_type.PT10CS07++
                case "PT15":
                    countWallet_type.PT15CS07++
                case "PT16":
                    countWallet_type.PT16CS07++
                case "PT17":
                    countWallet_type.PT17CS07++
                case "PT18":
                    countWallet_type.PT18CS07++
                case "PT19":
                    countWallet_type.PT19CS07++
                }
            case "08_FAMILYFOCUS":
                countOccupation.FAMILYFOCUS++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS08++
                case "PT02":
                    countWallet_type.PT02CS08++
                case "PT03":
                    countWallet_type.PT03CS08++
                case "PT05":
                    countWallet_type.PT05CS08++
                case "PT06":
                    countWallet_type.PT06CS08++
                case "PT08":
                    countWallet_type.PT08CS08++
                case "PT09":
                    countWallet_type.PT09CS08++
                case "PT10":
                    countWallet_type.PT10CS08++
                case "PT15":
                    countWallet_type.PT15CS08++
                case "PT16":
                    countWallet_type.PT16CS08++
                case "PT17":
                    countWallet_type.PT17CS08++
                case "PT18":
                    countWallet_type.PT18CS08++
                case "PT19":
                    countWallet_type.PT19CS08++
                }
            case "09_EARLYINCAREER":
                countOccupation.EARLYINCAREER++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS09++
                case "PT02":
                    countWallet_type.PT02CS09++
                case "PT03":
                    countWallet_type.PT03CS09++
                case "PT05":
                    countWallet_type.PT05CS09++
                case "PT06":
                    countWallet_type.PT06CS09++
                case "PT08":
                    countWallet_type.PT08CS09++
                case "PT09":
                    countWallet_type.PT09CS09++
                case "PT10":
                    countWallet_type.PT10CS09++
                case "PT15":
                    countWallet_type.PT15CS09++
                case "PT16":
                    countWallet_type.PT16CS09++
                case "PT17":
                    countWallet_type.PT17CS09++
                case "PT18":
                    countWallet_type.PT18CS09++
                case "PT19":
                    countWallet_type.PT19CS09++
                }
            case "10_LOWERMASS":
                countOccupation.LOWERMASS++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS10++
                case "PT02":
                    countWallet_type.PT02CS10++
                case "PT03":
                    countWallet_type.PT03CS10++
                case "PT05":
                    countWallet_type.PT05CS10++
                case "PT06":
                    countWallet_type.PT06CS10++
                case "PT08":
                    countWallet_type.PT08CS10++
                case "PT09":
                    countWallet_type.PT09CS10++
                case "PT10":
                    countWallet_type.PT10CS10++
                case "PT15":
                    countWallet_type.PT15CS10++
                case "PT16":
                    countWallet_type.PT16CS10++
                case "PT17":
                    countWallet_type.PT17CS10++
                case "PT18":
                    countWallet_type.PT18CS10++
                case "PT19":
                    countWallet_type.PT19CS10++
                }
            case "11_STUDENT":
                countOccupation.STUDENT++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS11++
                case "PT02":
                    countWallet_type.PT02CS11++
                case "PT03":
                    countWallet_type.PT03CS11++
                case "PT05":
                    countWallet_type.PT05CS11++
                case "PT06":
                    countWallet_type.PT06CS11++
                case "PT08":
                    countWallet_type.PT08CS11++
                case "PT09":
                    countWallet_type.PT09CS11++
                case "PT10":
                    countWallet_type.PT10CS11++
                case "PT15":
                    countWallet_type.PT15CS11++
                case "PT16":
                    countWallet_type.PT16CS11++
                case "PT17":
                    countWallet_type.PT17CS11++
                case "PT18":
                    countWallet_type.PT18CS11++
                case "PT19":
                    countWallet_type.PT19CS11++
                }
            case "12_RETIREHIGHWEALTH":
                countOccupation.RETIREHIGHWEALTH++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS12++
                case "PT02":
                    countWallet_type.PT02CS12++
                case "PT03":
                    countWallet_type.PT03CS12++
                case "PT05":
                    countWallet_type.PT05CS12++
                case "PT06":
                    countWallet_type.PT06CS12++
                case "PT08":
                    countWallet_type.PT08CS12++
                case "PT09":
                    countWallet_type.PT09CS12++
                case "PT10":
                    countWallet_type.PT10CS12++
                case "PT15":
                    countWallet_type.PT15CS12++
                case "PT16":
                    countWallet_type.PT16CS12++
                case "PT17":
                    countWallet_type.PT17CS12++
                case "PT18":
                    countWallet_type.PT18CS12++
                case "PT19":
                    countWallet_type.PT19CS12++
                }
            case "13_RETIREMEDIUMWEALTH":
                countOccupation.RETIREMEDIUMWEALTH++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS13++
                case "PT02":
                    countWallet_type.PT02CS13++
                case "PT03":
                    countWallet_type.PT03CS13++
                case "PT05":
                    countWallet_type.PT05CS13++
                case "PT06":
                    countWallet_type.PT06CS13++
                case "PT08":
                    countWallet_type.PT08CS13++
                case "PT09":
                    countWallet_type.PT09CS13++
                case "PT10":
                    countWallet_type.PT10CS13++
                case "PT15":
                    countWallet_type.PT15CS13++
                case "PT16":
                    countWallet_type.PT16CS13++
                case "PT17":
                    countWallet_type.PT17CS13++
                case "PT18":
                    countWallet_type.PT18CS13++
                case "PT19":
                    countWallet_type.PT19CS13++
                }
            case "14_RETIRELOWWEALTH":
                countOccupation.RETIRELOWWEALTH++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS14++
                case "PT02":
                    countWallet_type.PT02CS14++
                case "PT03":
                    countWallet_type.PT03CS14++
                case "PT05":
                    countWallet_type.PT05CS14++
                case "PT06":
                    countWallet_type.PT06CS14++
                case "PT08":
                    countWallet_type.PT08CS14++
                case "PT09":
                    countWallet_type.PT09CS14++
                case "PT10":
                    countWallet_type.PT10CS14++
                case "PT15":
                    countWallet_type.PT15CS14++
                case "PT16":
                    countWallet_type.PT16CS14++
                case "PT17":
                    countWallet_type.PT17CS14++
                case "PT18":
                    countWallet_type.PT18CS14++
                case "PT19":
                    countWallet_type.PT19CS14++
                }
            case "18_NEWCUST3MTH":
                countOccupation.NEWCUST3MTH++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS18++
                case "PT02":
                    countWallet_type.PT02CS18++
                case "PT03":
                    countWallet_type.PT03CS18++
                case "PT05":
                    countWallet_type.PT05CS18++
                case "PT06":
                    countWallet_type.PT06CS18++
                case "PT08":
                    countWallet_type.PT08CS18++
                case "PT09":
                    countWallet_type.PT09CS18++
                case "PT10":
                    countWallet_type.PT10CS18++
                case "PT15":
                    countWallet_type.PT15CS18++
                case "PT16":
                    countWallet_type.PT16CS18++
                case "PT17":
                    countWallet_type.PT17CS18++
                case "PT18":
                    countWallet_type.PT18CS18++
                case "PT19":
                    countWallet_type.PT19CS18++
                }
            case "99_OTH":
                countOccupation.OTH++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CS99++
                case "PT02":
                    countWallet_type.PT02CS99++
                case "PT03":
                    countWallet_type.PT03CS99++
                case "PT05":
                    countWallet_type.PT05CS99++
                case "PT06":
                    countWallet_type.PT06CS99++
                case "PT08":
                    countWallet_type.PT08CS99++
                case "PT09":
                    countWallet_type.PT09CS99++
                case "PT10":
                    countWallet_type.PT10CS99++
                case "PT15":
                    countWallet_type.PT15CS99++
                case "PT16":
                    countWallet_type.PT16CS99++
                case "PT17":
                    countWallet_type.PT17CS99++
                case "PT18":
                    countWallet_type.PT18CS99++
                case "PT19":
                    countWallet_type.PT19CS99++
                }    
            default:
                countOccupation.NULL_Custmer_segment++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01CSNULL++
                case "PT02":
                    countWallet_type.PT02CSNULL++
                case "PT03":
                    countWallet_type.PT03CSNULL++
                case "PT05":
                    countWallet_type.PT05CSNULL++
                case "PT06":
                    countWallet_type.PT06CSNULL++
                case "PT08":
                    countWallet_type.PT08CSNULL++
                case "PT09":
                    countWallet_type.PT09CSNULL++
                case "PT10":
                    countWallet_type.PT10CSNULL++
                case "PT15":
                    countWallet_type.PT15CSNULL++
                case "PT16":
                    countWallet_type.PT16CSNULL++
                case "PT17":
                    countWallet_type.PT17CSNULL++
                case "PT18":
                    countWallet_type.PT18CSNULL++
                case "PT19":
                    countWallet_type.PT19CSNULL++
                }
            }
            // Occupation
            switch customerAll.Occupation {
            case "Gov & State Enterprise":
                countOccupation.Gov_State_Enterprise++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01GOV++
                case "PT02":
                    countWallet_type.PT02GOV++
                case "PT03":
                    countWallet_type.PT03GOV++
                case "PT05":
                    countWallet_type.PT05GOV++
                case "PT06":
                    countWallet_type.PT06GOV++
                case "PT08":
                    countWallet_type.PT08GOV++
                case "PT09":
                    countWallet_type.PT09GOV++
                case "PT10":
                    countWallet_type.PT10GOV++
                case "PT15":
                    countWallet_type.PT15GOV++
                case "PT16":
                    countWallet_type.PT16GOV++
                case "PT17":
                    countWallet_type.PT17GOV++
                case "PT18":
                    countWallet_type.PT18GOV++
                case "PT19":
                    countWallet_type.PT19GOV++
                }
            case "Mass-Unidentify":
                countOccupation.Mass_Unidentify++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Mass++
                case "PT02":
                    countWallet_type.PT02Mass++
                case "PT03":
                    countWallet_type.PT03Mass++
                case "PT05":
                    countWallet_type.PT05Mass++
                case "PT06":
                    countWallet_type.PT06Mass++
                case "PT08":
                    countWallet_type.PT08Mass++
                case "PT09":
                    countWallet_type.PT09Mass++
                case "PT10":
                    countWallet_type.PT10Mass++
                case "PT15":
                    countWallet_type.PT15Mass++
                case "PT16":
                    countWallet_type.PT16Mass++
                case "PT17":
                    countWallet_type.PT17Mass++
                case "PT18":
                    countWallet_type.PT18Mass++
                case "PT19":
                    countWallet_type.PT19Mass++
                }
            case "N/A":
                countOccupation.NA++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01NA++
                case "PT02":
                    countWallet_type.PT02NA++
                case "PT03":
                    countWallet_type.PT03NA++
                case "PT05":
                    countWallet_type.PT05NA++
                case "PT06":
                    countWallet_type.PT06NA++
                case "PT08":
                    countWallet_type.PT08NA++
                case "PT09":
                    countWallet_type.PT09NA++
                case "PT10":
                    countWallet_type.PT10NA++
                case "PT15":
                    countWallet_type.PT15NA++
                case "PT16":
                    countWallet_type.PT16NA++
                case "PT17":
                    countWallet_type.PT17NA++
                case "PT18":
                    countWallet_type.PT18NA++
                case "PT19":
                    countWallet_type.PT19NA++
                }
            case "Salary":
                countOccupation.Salary++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Salary++
                case "PT02":
                    countWallet_type.PT02Salary++
                case "PT03":
                    countWallet_type.PT03Salary++
                case "PT05":
                    countWallet_type.PT05Salary++
                case "PT06":
                    countWallet_type.PT06Salary++
                case "PT08":
                    countWallet_type.PT08Salary++
                case "PT09":
                    countWallet_type.PT09Salary++
                case "PT10":
                    countWallet_type.PT10Salary++
                case "PT15":
                    countWallet_type.PT15Salary++
                case "PT16":
                    countWallet_type.PT16Salary++
                case "PT17":
                    countWallet_type.PT17Salary++
                case "PT18":
                    countWallet_type.PT18Salary++
                case "PT19":
                    countWallet_type.PT19Salary++
                }
            case "Self Employ (sSME)":
                countOccupation.Self_Employ++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01SE++
                case "PT02":
                    countWallet_type.PT02SE++
                case "PT03":
                    countWallet_type.PT03SE++
                case "PT05":
                    countWallet_type.PT05SE++
                case "PT06":
                    countWallet_type.PT06SE++
                case "PT08":
                    countWallet_type.PT08SE++
                case "PT09":
                    countWallet_type.PT09SE++
                case "PT10":
                    countWallet_type.PT10SE++
                case "PT15":
                    countWallet_type.PT15SE++
                case "PT16":
                    countWallet_type.PT16SE++
                case "PT17":
                    countWallet_type.PT17SE++
                case "PT18":
                    countWallet_type.PT18SE++
                case "PT19":
                    countWallet_type.PT19SE++
                }
            case "Student":
                countOccupation.Student++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01STD++
                case "PT02":
                    countWallet_type.PT02STD++
                case "PT03":
                    countWallet_type.PT03STD++
                case "PT05":
                    countWallet_type.PT05STD++
                case "PT06":
                    countWallet_type.PT06STD++
                case "PT08":
                    countWallet_type.PT08STD++
                case "PT09":
                    countWallet_type.PT09STD++
                case "PT10":
                    countWallet_type.PT10STD++
                case "PT15":
                    countWallet_type.PT15STD++
                case "PT16":
                    countWallet_type.PT16STD++
                case "PT17":
                    countWallet_type.PT17STD++
                case "PT18":
                    countWallet_type.PT18STD++
                case "PT19":
                    countWallet_type.PT19STD++
                }
            case "Wealth":
                countOccupation.Wealth++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Wealth++
                case "PT02":
                    countWallet_type.PT02Wealth++
                case "PT03":
                    countWallet_type.PT03Wealth++
                case "PT05":
                    countWallet_type.PT05Wealth++
                case "PT06":
                    countWallet_type.PT06Wealth++
                case "PT08":
                    countWallet_type.PT08Wealth++
                case "PT09":
                    countWallet_type.PT09Wealth++
                case "PT10":
                    countWallet_type.PT10Wealth++
                case "PT15":
                    countWallet_type.PT15Wealth++
                case "PT16":
                    countWallet_type.PT16Wealth++
                case "PT17":
                    countWallet_type.PT17Wealth++
                case "PT18":
                    countWallet_type.PT18Wealth++
                case "PT19":
                    countWallet_type.PT19Wealth++
                }
            case "Welfare":
                countOccupation.Welfare++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01Welfare++
                case "PT02":
                    countWallet_type.PT02Welfare++
                case "PT03":
                    countWallet_type.PT03Welfare++
                case "PT05":
                    countWallet_type.PT05Welfare++
                case "PT06":
                    countWallet_type.PT06Welfare++
                case "PT08":
                    countWallet_type.PT08Welfare++
                case "PT09":
                    countWallet_type.PT09Welfare++
                case "PT10":
                    countWallet_type.PT10Welfare++
                case "PT15":
                    countWallet_type.PT15Welfare++
                case "PT16":
                    countWallet_type.PT16Welfare++
                case "PT17":
                    countWallet_type.PT17Welfare++
                case "PT18":
                    countWallet_type.PT18Welfare++
                case "PT19":
                    countWallet_type.PT19Welfare++
                }
            default:
                countOccupation.NULL_Occupation++
                switch inv.Wallet_type {
                case "PT01":
                    countWallet_type.PT01OCNULL++
                case "PT02":
                    countWallet_type.PT02OCNULL++
                case "PT03":
                    countWallet_type.PT03OCNULL++
                case "PT05":
                    countWallet_type.PT05OCNULL++
                case "PT06":
                    countWallet_type.PT06OCNULL++
                case "PT08":
                    countWallet_type.PT08OCNULL++
                case "PT09":
                    countWallet_type.PT09OCNULL++
                case "PT10":
                    countWallet_type.PT10OCNULL++
                case "PT15":
                    countWallet_type.PT15OCNULL++
                case "PT16":
                    countWallet_type.PT16OCNULL++
                case "PT17":
                    countWallet_type.PT17OCNULL++
                case "PT18":
                    countWallet_type.PT18OCNULL++
                case "PT19":
                    countWallet_type.PT19OCNULL++
                }
            }
		}  else {
			// ถ้าไม่มี Customer_code ตรงกันเลย ให้เพิ่มค่า NULL_Usage_segment
			countOccupation.NULL_Usage_segment++
            countOccupation.NULL_Gender++
            countOccupation.NULL_Custmer_segment++

            switch inv.Wallet_type {
            case "PT01":
                countWallet_type.PT01NULL++
                countWallet_type.PT01GNULL++
                countWallet_type.PT01AgeNULL++
                countWallet_type.PT01CSNULL++
                countWallet_type.PT01OCNULL++
            case "PT02":
                countWallet_type.PT02NULL++
                countWallet_type.PT02GNULL++
                countWallet_type.PT02AgeNULL++
                countWallet_type.PT02CSNULL++
                countWallet_type.PT02OCNULL++
            case "PT03":
                countWallet_type.PT03NULL++
                countWallet_type.PT03GNULL++
                countWallet_type.PT03AgeNULL++
                countWallet_type.PT03CSNULL++
                countWallet_type.PT03OCNULL++
            case "PT05":
                countWallet_type.PT05NULL++
                countWallet_type.PT05GNULL++
                countWallet_type.PT05AgeNULL++
                countWallet_type.PT05CSNULL++
                countWallet_type.PT05OCNULL++
            case "PT06":
                countWallet_type.PT06NULL++
                countWallet_type.PT06GNULL++
                countWallet_type.PT06AgeNULL++
                countWallet_type.PT06CSNULL++
                countWallet_type.PT06OCNULL++
            case "PT08":
                countWallet_type.PT08NULL++
                countWallet_type.PT08GNULL++
                countWallet_type.PT08AgeNULL++
                countWallet_type.PT08CSNULL++
                countWallet_type.PT08OCNULL++
            case "PT09":
                countWallet_type.PT09NULL++
                countWallet_type.PT09GNULL++
                countWallet_type.PT09AgeNULL++
                countWallet_type.PT09CSNULL++
                countWallet_type.PT09OCNULL++
            case "PT10":
                countWallet_type.PT10NULL++
                countWallet_type.PT10GNULL++
                countWallet_type.PT10AgeNULL++
                countWallet_type.PT10CSNULL++
                countWallet_type.PT10OCNULL++
            case "PT15":
                countWallet_type.PT15NULL++
                countWallet_type.PT15GNULL++
                countWallet_type.PT15AgeNULL++
                countWallet_type.PT15CSNULL++
                countWallet_type.PT15OCNULL++
            case "PT16":
                countWallet_type.PT16NULL++
                countWallet_type.PT16GNULL++
                countWallet_type.PT16AgeNULL++
                countWallet_type.PT16CSNULL++
                countWallet_type.PT16OCNULL++
            case "PT17":
                countWallet_type.PT17NULL++
                countWallet_type.PT17GNULL++
                countWallet_type.PT17AgeNULL++
                countWallet_type.PT17CSNULL++
                countWallet_type.PT17OCNULL++
            case "PT18":
                countWallet_type.PT18NULL++
                countWallet_type.PT18GNULL++
                countWallet_type.PT18AgeNULL++
                countWallet_type.PT18CSNULL++
                countWallet_type.PT18OCNULL++
            case "PT19":
                countWallet_type.PT19NULL++
                countWallet_type.PT19GNULL++
                countWallet_type.PT19AgeNULL++
                countWallet_type.PT19CSNULL++
                countWallet_type.PT19OCNULL++
            }
		}

		countOccupation.Cusotmer_total++
	}

	return countOccupation,countWallet_type
}
