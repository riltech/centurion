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
	// Returns if the solution is the 5th or 10th or 15th (etc..)
	// NOTE: the given challengeId has to be pre-solved in combat
	IsFifthUniqueSolution(attackerID string, challengeID string) bool
	// Calculates how many percentage of challenges the attackers were able to solve
	GetOverallAttackerSuccessPrecent(numberOfUniqueChallenges int) int
	// Returns the number of how many attackers completed a given challenge
	GetNumberOfUniqueCompletionsPerChallenges() map[string]uint
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
	if len(archive) == 0 {
		return 0
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

func (s Service) IsFifthUniqueSolution(attackerID string, challengeID string) bool {
	archive := s.repository.GetArchive()
	uniqueCompleted := map[string]int{}
	for _, item := range archive {
		if item.AttackerID == attackerID {
			uniqueCompleted[item.ChallengeID] = 1
		}
	}
	toStartWith := 0
	if _, ok := uniqueCompleted[challengeID]; !ok {
		toStartWith++
	}
	for range uniqueCompleted {
		toStartWith++
	}
	return toStartWith%5 == 0
}

func (s Service) GetOverallAttackerSuccessPrecent(numberOfUniqueChallenges int) int {
	uniques := map[string]uint8{}
	for _, a := range s.repository.GetArchive() {
		if a.CombatState == CombatStateAttackSucceeded {
			uniques[a.ChallengeID] = 1
		}
	}
	if numberOfUniqueChallenges == 0 {
		return 100
	}
	return int((float32(len(uniques)) / float32(numberOfUniqueChallenges)) * 100)
}

func (s Service) GetNumberOfUniqueCompletionsPerChallenges() map[string]uint {
	challengeCompletion := map[string]uint{}
	archive := s.repository.GetArchive()
	for _, a := range archive {
		if a.CombatState != CombatStateAttackSucceeded {
			continue
		}
		if _, ok := challengeCompletion[a.ID]; !ok {
			challengeCompletion[a.ID] = 0
		}
		challengeCompletion[a.ID] = challengeCompletion[a.ID] + 1
	}
	return challengeCompletion
}

func NewService(repository IRepository) IService {
	return &Service{repository}
}
