package preload

import(
	"pt-report-backend/invitation"
	"pt-report-backend/transaction"
	"log"
	"time"
)


type PreloadService struct{
	InvitationService *invitation.Service
	TransactionService *transaction.Service

}



func (p *PreloadService)RunPreload(){
	// invitationService := invitation.NewService()
	// transactionService := transaction.NewService()
	startDate := "2025-01-01"
	endDate := time.Now().Format("2006-01-02")
	_,err := p.TransactionService.PreloadTransactionCache(startDate,endDate)
	
	if err != nil {
		log.Printf("Preload transactions failed: %v", err)
	}

	// _, err = invitationService.PreloadInvitationsCache(db.DB, startDate, endDate, "invitationDate", "1")
	_, err = p.InvitationService.PreloadInvitationsCache(startDate, endDate)

	if err != nil {
		log.Printf("Preload invitations failed: %v", err)
	}
	// test,err := invitationService.GetCachedInvitations(startDate,endDate,"invitationDate","1")
	// if err != nil {
	// 	log.Printf("data: %v", err)
	// } 
	// log.Println(test)
	_, err = p.InvitationService.PreloadCustomers(startDate, endDate)

	if err != nil {
		log.Printf("Preload customers failed: %v", err)
	}
	// test,err := invitationService.GetCachedCustomers(startDate,endDate,"0")
	// if err != nil {
	// 	log.Printf("data: %v", err)
	// } 
	log.Println("Preload Success")
}