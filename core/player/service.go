package player

import (
	"strings"
)

// Describes a player service interface
type IService interface {
	// Adds a new player
	AddPlayer(Model) error
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

func NewService(repository IRepository) IService {
	return &Service{repository}
}
