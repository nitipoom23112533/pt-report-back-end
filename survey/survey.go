package survey

type SurveyService struct {
	SurveyRepo *SurveyRepo


}
type Survey struct {
	SurveyId string `db:"survey_id" json:"surveyId"`
	SurveyName string `db:"survey_name" json:"surveyName"`
	IsActive bool `db:"is_active" json:"isActive"`

}

func NewSurveyService() *SurveyService{

	return &SurveyService{}
}

func(ss *SurveyService)GetSurvey() ([]Survey,error) {
	return ss.SurveyRepo.GetSurvey()
}

func(ss *SurveyService)UpdateSurvey(survey *Survey) error {
	return ss.SurveyRepo.UpdateSurvey(survey)

}