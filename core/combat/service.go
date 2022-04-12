package combat

import "fmt"

// Describes a combat service interface
type IService interface {
	// Adds a new player
	AddCombat(Model) error
	// Finds a combat by ID
	FindByID(ID string) (Model, error)
	// Updates the CombatState of a given combat
	// Find available states in the package
	UpdateCombatState(ID string, state string) (Model, error)
	// Find by attacker id and challenge id
	FindByAttackerAndChallenge(attackerID string, challengeID string) (Model, error)
}

type Service struct {
	repository IRepository
}

// Interface check
var _ IService = (*Service)(nil)

func (s Service) AddCombat(p Model) error {
	return s.repository.AddCombat(p)
}

func (s Service) FindByID(ID string) (Model, error) {
	return s.repository.FindByID(ID)
}

func (s *Service) UpdateCombatState(ID string, state string) (Model, error) {
	if s == nil {
		return Model{}, fmt.Errorf("Cannot update state without a service instance")
	}
	m, err := s.FindByID(ID)
	if err != nil {
		return Model{}, err
	}
	found := false
	for _, validState := range CombatStateCollection {
		if state == validState {
			found = true
			break
		}
	}
	if !found {
		return Model{}, fmt.Errorf("%s is not a valid state", state)
	}
	return s.repository.UpdateCombatState(m.ID, state)
}

func (s Service) FindByAttackerAndChallenge(attackerID string, challengeID string) (Model, error) {
	combats := s.repository.GetCombats()
	for _, c := range combats {
		if c.AttackerID != attackerID {
			continue
		}
		if challengeID != c.ChallengeID {
			continue
		}
		return c, nil
	}
	return Model{}, fmt.Errorf("Not found")
}

func NewService(repository IRepository) IService {
	return &Service{repository}
}
