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
	// Add given amounts of points to a player
	AddPoint(ID string, points int) (Model, error)
	// Fetch a given team completely
	// use Team enums for the string
	GetTeam(team string) []Model
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

func (s Service) AddPoint(ID string, points int) (Model, error) {
	return s.repository.AddPoint(ID, points)
}

func (s Service) GetTeam(team string) []Model {
	players := s.repository.GetPlayers()
	toReturn := []Model{}
	for _, p := range players {
		if p.Team == team {
			toReturn = append(toReturn, p)
		}
	}
	return toReturn
}

func NewService(repository IRepository) IService {
	return &Service{repository}
}
