package scoreboard

import "github.com/riltech/centurion/core/player"

// Describes a scoreboard service interface
type IService interface {
	// Returns the score board for each team
	GetBoards() (attacker Model, defender Model)
	// Adds a given point to a team and to a player
	AddPoint(playerID string, point int) error
}

// Service implementation
type Service struct {
	repository    IRepository
	playerService player.IService
}

// Interface check
var _ IService = (*Service)(nil)

// Constructor for the scoreboard service
func NewService(repository IRepository, playerService player.IService) IService {
	return &Service{
		repository:    repository,
		playerService: playerService,
	}
}

func (s Service) GetBoards() (Model, Model) {
	return s.repository.GetBoards()
}

func (s Service) AddPoint(playerID string, points int) error {
	p, err := s.playerService.AddPoint(playerID, points)
	if err != nil {
		return err
	}
	s.repository.AddPoint(p.Team, points)
	return nil
}
