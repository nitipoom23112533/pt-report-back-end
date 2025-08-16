package transaction

import (
	"pt-report-backend/invitation"
	"log"
    "pt-report-backend/db"
    "sync"
    "time"

)

type Service struct{
    Cache      []Transaction
    CacheMutex sync.RWMutex

}

func NewService() *Service {
	return &Service{}
}


func (s *Service)GetAllTransaction(startDate, endDate string)([]invitation.Invitation,error)  {
    log.Println("GetAllTransaction called")
    query := `SELECT Customer_code, Wallet_type
                            FROM pt_transaction FORCE INDEX (idx_t_date_wallet_type)
                            WHERE T_date BETWEEN ? AND ?;`
    var transaction []invitation.Invitation
    err := db.DB.Select(&transaction,query, startDate, endDate)
    if err != nil {
        log.Printf("Error fetching invitation: %v", err)
		return nil, err
    }
    log.Println("Transaction fetched successfully")
    return transaction,nil
}

func (s *Service) PreloadTransactionCache(startDate, endDate string) ([]Transaction, error) {

    query := `SELECT Customer_code, Wallet_type,T_date
                            FROM pt_transaction 
                            WHERE T_date BETWEEN ? AND ?;`
	var transactions []Transaction
	err := db.DB.Select(&transactions, query, startDate, endDate)
	if err != nil {
		log.Printf("Error fetching invitation: %v", err)
		return nil, err
	}

	s.CacheMutex.Lock()
	s.Cache = transactions
	s.CacheMutex.Unlock()

	log.Println("Transactions cache refreshed")
	return transactions, nil
}

func (s *Service) GetCachedTransactions(startDate string, endDate string) ([]Transaction,error) {
    log.Println("GetCachedTransactions called")

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

    var filtered []Transaction
    for _, tra := range s.Cache {
            invDate, err := time.Parse(datetimeLayout, tra.T_date) // สมมติเก็บวันที่ใน field DateStr เป็น string
            if err != nil {
                log.Println(err)
            }
        
            if (invDate.Equal(start) || invDate.After(start)) && (invDate.Equal(end) || invDate.Before(end)) {
           
                filtered = append(filtered, tra)
                
            }
        }
    
    return filtered, nil
}
