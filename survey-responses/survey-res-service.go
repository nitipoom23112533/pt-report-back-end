package surveyresponses

import(
	"pt-report-backend/db"
	"log"
)

type SurveyResRepo struct {
}

func NewSurveyResRepo() *SurveyResRepo {
	return &SurveyResRepo{}
}


func (sr *SurveyResRepo) ResponsesQuery(cCode string) (error) {
    tx, err := db.DB.Beginx()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    query := `INSERT INTO survey_responses (customer_code)
                    VALUES (:customer_code)
                    ON DUPLICATE KEY UPDATE
                    response_date = CURRENT_TIMESTAMP,
                    updated_at = CURRENT_TIMESTAMP;`

    CusCode := map[string]interface{}{
			"customer_code": cCode,
		}
    _, err = tx.NamedExec(query, CusCode)
		if err != nil {
            log.Println(err.Error())
			return err
		}

    return tx.Commit()
}