package transaction

type Transaction struct {
	Customer_code string `db:"Customer_code"`
	Wallet_type string `db:"Wallet_type"`
}