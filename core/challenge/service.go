package challenge

import (
	"fmt"
)

// Describes a player service interface
type IService interface {
	// Used for adding the default defender modules in the beginning of the game
	AddDefaultModules() error
	// Adds a new challenge to the system
	AddChallenge(Model) error
	// For fetching available challenges
	GetChallenges() []Model
	// Finds a given challenge by ID
	FindByID(ID string) (Model, error)
	// Generates hint for a default challenge
	GenerateHintForDefault(Model) ([]interface{}, error)
	// Validates a given solution for a default module
	IsValidSolutionToDefaultModule(m Model, hints []interface{}, solutions []interface{}) (bool, error)
	// Returns if this is the first module of the defender
	IsFirstModule(Model) bool
}

type Service struct {
	repository IRepository
}

// Interface check
var _ IService = (*Service)(nil)

func NewService(repository IRepository) IService {
	return &Service{repository}
}

func (s Service) GetChallenges() []Model {
	return s.repository.GetChallenges()
}

func (s Service) FindByID(ID string) (m Model, e error) {
	for _, challenge := range s.repository.GetChallenges() {
		if challenge.ID == ID {
			m = challenge
			return
		}
	}
	e = fmt.Errorf("%s challenge is not found", ID)
	return
}

func (s Service) AddChallenge(m Model) error {
	return s.repository.AddChallenge(m)
}

func (s Service) AddDefaultModules() error {
	for _, challenge := range getDefaultChallenges() {
		if err := s.repository.AddChallenge(challenge); err != nil {
			return err
		}
	}
	return nil
}

func (s Service) GenerateHintForDefault(m Model) ([]interface{}, error) {
	return getHintsForDefaultModule(m)
}

func (s Service) IsValidSolutionToDefaultModule(m Model, hints []interface{}, solutions []interface{}) (bool, error) {
	return isValidDefaultModuleSolution(m, hints, solutions)
}

func (s Service) IsFirstModule(m Model) bool {
	numberOfModules := 0
	for _, c := range s.repository.GetChallenges() {
		if c.CreatorID == m.CreatorID {
			numberOfModules++
		}
	}
	return numberOfModules == 0
}
