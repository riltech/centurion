package combat

const (
	// FLOW

	// Starting state for the combat
	CombatStateAttackInitiated = "attack_initiated"
	// After attack initiation on defense modules
	// provided by defenders, the defender is requested
	// to provide hints for the attacker
	CombatStateDefenseRequested = "defense_requested"
	// Happens after hints are received from the defender
	// and sent back to the attacker
	CombatStateAttackerChallenged = "attacker_challenged"
	// Happens after the attacker had received hints and provided a solution for a challenge
	CombatStateSolutionProvided = "solution_provided"
	// Happens when a defender has to validate a given solution for a given challenge
	CombatStateSolutionEvaluationRequested = "solution_validation_requested"
	// Happens when the defender evaluated a given solution
	CombatStateSolutionEvaluated = "solution_validated"

	// ERRORS (Also final states)

	// Happens when the defense fails to complete
	// the challenge flow for some reason (e.g. offline, invalid flow handling, missing data)
	CombatStateDefenseFailed = "defense_failed"
	// Happens when the attacker fails to complete
	// the challenge flow for some reason (e.g. offline, invalid flow handling, missing data)
	CombatStateAttackFailed = "attack_failed"

	// END STATES

	// Happens when the defender succeeds in defending their module
	CombatStateDefenseSucceeded = "defense_succeeded"
	// Happens when the attacker successfully penetrates a defense module
	CombatStateAttackSucceeded = "attack_succeeded"
)

// Collection of the available CombatStates
var CombatStateCollection = []string{
	CombatStateAttackFailed,
	CombatStateAttackInitiated,
	CombatStateAttackSucceeded,
	CombatStateAttackerChallenged,
	CombatStateDefenseFailed,
	CombatStateDefenseRequested,
	CombatStateDefenseSucceeded,
	CombatStateSolutionEvaluated,
	CombatStateSolutionEvaluationRequested,
	CombatStateSolutionProvided,
}

// Final combat states are states which after the state should not
// update anymore
func IsFinalCombatState(state string) bool {
	found := false
	for _, combatState := range CombatStateCollection {
		if combatState == state {
			found = true
			break
		}
	}
	return found &&
		(state == CombatStateAttackFailed ||
			state == CombatStateDefenseFailed ||
			state == CombatStateAttackSucceeded ||
			state == CombatStateDefenseSucceeded)
}
