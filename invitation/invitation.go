package invitation

import(
	"time"
)
type Invitation struct {
	EDR_id      string `db:"EDR_id"`
	Wallet_type string `db:"Wallet_type"`
	IN_date     string `db:"IN_date"`
	T_date      string `db:"T_date"`
}

type Duration struct {
	Start_date time.Time `db:"start_date"`
	End_date   time.Time `db:"end_date"`
}
