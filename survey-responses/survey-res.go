package surveyresponses
import(
)


type SurveyResService struct {
	SurveyResRepo *SurveyResRepo
}

func NewSurveyResService() *SurveyResService {
	return &SurveyResService{}
}

func (s *SurveyResService) Responses(cCode string) error {
	return s.SurveyResRepo.ResponsesQuery(cCode)
}
