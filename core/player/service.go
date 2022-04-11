package player

import (
	"strings"
)

// Describes a player service interface
type IService interface {
	// Adds a new player
	AddPlayer(Model) error
	// Sets a player's online status
	SetPlayerOnlineStatus(ID string, online bool) (Model, error)
	// Finds a player by ID
	FindByID(ID string) (Model, error)
	// Checks if a given player is already registered or not
	IsPlayerExist(*Model) bool
}

type Service struct {
	repository IRepository
}

// Interface check
var _ IService = (*Service)(nil)

func (s Service) AddPlayer(p Model) error {
	return s.repository.AddPlayer(p)
}

func (s Service) IsPlayerExist(p *Model) bool {
	if p == nil {
		return false
	}
	for _, storedPlayer := range s.repository.GetPlayers() {
		if storedPlayer.ID == p.ID {
			return true
		}
		if strings.ToLower(storedPlayer.Name) == strings.ToLower(p.Name) {
			return true
		}
	}
	return false
}

func (s Service) SetPlayerOnlineStatus(ID string, online bool) (Model, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return Model{}, err
	}
	user.Online = online
	return s.repository.UpdateByID(ID, user)
}

func (s Service) FindByID(ID string) (Model, error) {
	return s.repository.FindByID(ID)
}

func NewService(repository IRepository) IService {
	return &Service{repository}
}
