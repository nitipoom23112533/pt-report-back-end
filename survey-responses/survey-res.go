package surveyresponses
import(
	"log"
)


type SurveyResService struct {
	SurveyResRepo *SurveyResRepo
}

func NewSurveyResService() *SurveyResService {
	return &SurveyResService{}
}

func (s *SurveyResService) Responses(cCode string) error {
	log.Println("Start Import Responses")
	return s.SurveyResRepo.ResponsesQuery(cCode)
}
