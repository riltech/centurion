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
	// Returns true of the attacker has already completed the given challenge before
	// during the game session
	IsAttackerCompletedBefore(attackerID string, challengeID string) bool
	// Returns the percentages of failed defense events
	// values are between 0-100
	GetDefenseFailPercent() int
	// Returns the percentages of successful attack events
	// values are between 0-100
	GetAttackerSuccessPercent() int
}

// Service implementation
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

func (s Service) IsAttackerCompletedBefore(attackerID, challengeID string) bool {
	finishedCombats := s.repository.GetArchive()
	for _, c := range finishedCombats {
		if c.AttackerID == attackerID && c.ChallengeID == challengeID && c.CombatState == CombatStateAttackSucceeded {
			return true
		}
	}
	return false
}

func (s Service) GetDefenseFailPercent() int {
	failedEvents := []Model{}
	archive := s.repository.GetArchive()
	for _, c := range archive {
		if c.CombatState == CombatStateDefenseFailed {
			failedEvents = append(failedEvents, c)
		}
	}
	return int(float32(len(failedEvents)) / float32(len(archive)) * 100)
}

func (s Service) GetAttackerSuccessPercent() int {
	successEvents := []Model{}
	validEventsOverall := []Model{}
	archive := s.repository.GetArchive()
	for _, c := range archive {
		if c.CombatState == CombatStateAttackSucceeded {
			successEvents = append(successEvents, c)
			validEventsOverall = append(validEventsOverall, c)
			continue
		}
		if c.CombatState == CombatStateDefenseSucceeded {
			validEventsOverall = append(validEventsOverall, c)
			continue
		}
		if c.CombatState == CombatStateAttackFailed {
			validEventsOverall = append(validEventsOverall, c)
			continue
		}
	}
	if len(validEventsOverall) == 0 {
		return 100
	}
	return int(float32(len(successEvents)) / float32(len(validEventsOverall)) * 100)
}

func NewService(repository IRepository) IService {
	return &Service{repository}
}
