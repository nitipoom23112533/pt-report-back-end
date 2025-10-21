package survey

import(
	"pt-report-backend/db"
)

type SurveyRepo struct{

}

func NewSurveyRepo() *SurveyRepo {
	return &SurveyRepo{}
}
func (sr *SurveyRepo) GetSurvey() ([]Survey, error) {
	query := `SELECT survey_id,survey_name,is_active FROM adhoc_survey ORDER BY is_active DESC`

	var s []Survey
	err := db.DB.Select(&s,query)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func(sr *SurveyRepo) UpdateSurvey(survey *Survey) error {
	tx,err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `UPDATE adhoc_survey SET survey_id = :survey_id, survey_name = :survey_name, is_active = :is_active WHERE id = 1`
	_,err = tx.NamedExec(query,survey)
	if err != nil {
		return err
	}

	return tx.Commit()
}
