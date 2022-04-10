package challenge

// Describes a player service interface
type IService interface {
	// For fetching available challenges
	GetChallenges() []Model
}

type Service struct {
	repository IRepository
}

// Interface check
var _ IService = (*Service)(nil)

func (s Service) GetChallenges() []Model {
	return s.repository.GetChallenges()
}

func NewService(repository IRepository) IService {
	return &Service{repository}
}
