package urlmanager

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) CheckURL(url string) int {
	return 0

}
