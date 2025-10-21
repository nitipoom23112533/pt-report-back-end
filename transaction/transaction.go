package transaction

type Transaction struct {
	// Customer_code 	string `db:"Customer_code"`
	EDR_id      	string `db:"Customer_code"`
	Wallet_type 	string `db:"Wallet_type"`
	T_date      	string `db:"T_date"`
}