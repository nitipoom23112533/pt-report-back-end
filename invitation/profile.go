package invitation

type Customer struct {
	CustomerCode string `db:"Customer_code"` // กำหนดชื่อฟิลด์ให้ตรงกับ DB
	Occupation string `db:"Occupation"`
	Customer_segment string `db:"Customer_segment"`
	Usage_segment string `db:"Usage_segment"`
	Age_range string `db:"Age_range"`
	Gender string `db:"Gender"`
	Customer_date string `db:"Customer_date"`
}
type CountOccupation struct {
	// //Occupation
	Gov_State_Enterprise  int `json:"Gov_State_Enterprise"`
	Mass_Unidentify  int `json:"Mass_Unidentify"`
	NA  int `json:"Na"`
	Salary  int `json:"Salary"`
	Self_Employ  int `json:"Self_Employ"`
	Student  int `json:"Student"`
	Wealth  int `json:"Wealth"`
	Welfare  int `json:"Welfare"`
	NULL_Occupation int `json:"NULL_Occupation"`

	//Custmer_segment
	PRECIOUSPLUS int `json:"Preciousplus"`
	PRECIOUS int `json:"Precious"`
	PREWEALTH int `json:"Prewealth"`
	AFFUIENTTOBE int `json:"Affuienttobe"`
	RETIREPLANNER int `json:"Retireplanner"`
	BUILDUPFORFEATURE int `json:"Buildupforfeature"`
	FAMILYFOCUS int `json:"Familyfocus"`
	EARLYINCAREER int `json:"Earlyincareer"`
	LOWERMASS int `json:"Lowermass"`
	STUDENT int `json:"Student2"`
	RETIREHIGHWEALTH int `json:"Retierehighwealth"`
	RETIREMEDIUMWEALTH int `json:"Retiremediumwealth"`
	RETIRELOWWEALTH int `json:"Retirelowealth"`
	NEWCUST3MTH int `json:"Newcust3mth"`
	OTH int `json:"Oth"`
	NULL_Custmer_segment int `json:"NULL_Custmer_segment"`

	// Life Stage Segment
	CareerStarterLower int `json:"CareerStarterLower"`
	CareerStarterMiddle int `json:"CareerStarterMiddle"`
	CareerStarterUpper int `json:"CareerStarterUpper"`
	ChildrenOrStudent int `json:"ChildrenOrStudent"`
	FutureBuilder int `json:"FutureBuilder"`
	LowerMass int `json:"LowerMass"`
	MassLower int `json:"MassLower"`
	MassMiddle int `json:"MassMiddle"`
	MassUpper int `json:"MassUpper"`
	LSSNA int `json:"LSSNA"`
	PreSenior int `json:"Presenior"`
	SeniorLower int `json:"SeniorLower"`
	SeniorUpper int `json:"SeniorUpper"`
	UniversityStudent int `json:"UniversityStudent"`
	WealthToBe int `json:"WealthToBe"`
	WealthPotentail int `json:"WealthPotentail"`
	LSSWealth int `json:"LSSWealth"`

	//Usage_segment
	Low                int `json:"Low"`
	Medium             int `json:"Medium"`
	High               int `json:"High"`
	Login_Only         int `json:"Login_Only"`
	Screen_View_Only   int `json:"Screen_View_Only"`
	Inactive           int `json:"Inactive"`
	New_User           int `json:"New_User"`
	NULL_Usage_segment int `json:"NULL_Usage_segment"`

	//Age_range
	Age_range_1 int `json:"Age_range_1"`
	Age_range_2 int `json:"Age_range_2"`
	Age_range_3 int `json:"Age_range_3"`
	Age_range_4 int `json:"Age_range_4"`
	Age_range_5 int `json:"Age_range_5"`
	Age_range_6 int `json:"Age_range_6"`
	Age_range_7 int `json:"Age_range_7"`
	Age_range_8 int `json:"Age_range_8"`
	NULL_Age_range int `json:"NULL_Age_range"`

	//Gender
	Gender_M int `json:"Gender_M"`
	Gender_F int `json:"Gender_F"`
	NULL_Gender int `json:"NULL_Gender"`

	//Cusotmer_total
	Cusotmer_total int `json:"Cusotmer_total"`
}
