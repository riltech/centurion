package combat

import "time"

// Describes a combat in the system
// which can happen between attackers and defenders
type Model struct {
	// ID of the combat
	ID string
	// ID of the challenge being solved
	ChallengeID string
	// ID of the attacker
	AttackerID string
	// ID of the defender
	DefenderID string
	// Current state of the combat
	CombatState string
	// Time of creation
	CreatedAt time.Time
	// Time of the last update on the model
	LastUpdateAt time.Time
}

// Returns if the state is in the final stage
// which indicates that it should be immutable
func (m Model) IsInFinalState() bool {
	return IsFinalCombatState(m.CombatState)
}
