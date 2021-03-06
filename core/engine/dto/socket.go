package dto

// Event type for a player joining the game
const SocketEventTypeError = "error"
const SocketEventTypeJoin = "join"
const SocketEventTypeAttack = "attack"
const SocketEventTypeAttackResult = "attack_result"
const SocketEventTypeAttackChallenge = "attack_challenge"
const SocketEventTypeAttackSolution = "attack_solution"
const SocketEventTypeAttackerFailedToAttack = "attacker_failed_to_attack"
const SocketEventTypeDefenderFailedToDefend = "defender_failed_to_defend"
const SocketEventTypeDefendActionRequest = "defend_action_request"
const SocketEventTypeDefendAction = "defend_action"
const SocketEventTypeSolutionEvaluationRequest = "solution_evaluation_request"
const SocketEventTypeSolutionEvaluation = "solution_evaluation"

// Describes a generic event over websockets
type SocketEvent struct {
	Type string `json:"type"`
}

// Describes a player joined event
type JoinEvent struct {
	SocketEvent
	// ID which the player uses from the registration
	ID string `json:"id"`
}

// Happens when a new attack is launched
type AttackEvent struct {
	SocketEvent
	// ID of the challenge
	TargetID string `json:"targetId"`
}

// Happens when an attack is for a valid target ID
// The challenge creator generates hints to resolve
type AttackChallengeEvent struct {
	SocketEvent
	// ID of the challenge
	TargetID string `json:"targetId"`
	// Hints for the challenge
	Hints []interface{} `json:"hints"`
}

// Happens when an attacker already acquired hints
// and now sending in the solution for given hints
type AttackSolutionEvent struct {
	SocketEvent
	// ID of the given challenge
	TargetID string `json:"targetId"`
	// Hints received in the previous event
	Hints []interface{} `json:"hints"`
	// Solutions for the hints
	Solutions []interface{} `json:"solutions"`
}

// Sent back when an attack was evaluated
type AttackResultEvent struct {
	SocketEvent
	// ID of the challenge
	TargetID string `json:"targetId"`
	// Result of the attack
	Success bool `json:"success"`
	// Optional message to pass on from the defender
	Message string `json:"message"`
}

// Sent when an error happens during an action
type ErrorEvent struct {
	SocketEvent
	Message string `json:"message"`
}

// Happens when a defender is not online to provide
// hints for a challenge
type DefenderFailedToDefendEvent struct {
	SocketEvent
	// ID of the challenge
	TargetID string `json:"targetId"`
}

// Happens when an attacker cannot provide solutions for hints
// or goes offline during the attack flow
type AttackerFailedToAttackEvent struct {
	SocketEvent
	// ID of the challenge
	TargetID string `json:"targetId"`
	// ID of the combat
	CombatID string `json:"combatId"`
}

// Happens when an attacker attacks a challenge module
// installed by a defender and the module owner
// has to defend
type DefendActionRequestEvent struct {
	SocketEvent
	// ID of the challenge
	TargetID string `json:"targetId"`
	// Combat ID is the ID of the combat in which
	// the action is requested
	// NOTE: This needs to be echo'd back in the action
	CombatID string `json:"combatId"`
}

// Happens when a defender sends out hints
// for the attacker
type DefendActionEvent struct {
	SocketEvent
	// Hints generated for the challenge
	Hints []interface{} `json:"hints"`
	// ID of the given combat
	CombatID string `json:"combatID"`
}

// Happens when the attacker hands in a solution
// and the defender is requested to evaluate the solution
type SolutionEvaluationRequestEvent struct {
	SocketEvent
	// ID of the challenge
	TargetID string `json:"targetId"`
	// ID of the given combat
	CombatID string `json:"combatId"`
	// Solution array
	Solutions []interface{} `json:"solutions"`
	// Hints used for generating the solution
	Hints []interface{} `json:"hints"`
}

// Happens when the defender is done with the
// evaluation of a given solution
type SolutionEvaluationEvent struct {
	SocketEvent
	// ID of the challenge
	TargetID string `json:"targetId"`
	// ID of the given combat
	CombatID string `json:"combatId"`
	// Indicates if the evaluation was successful or not
	Success bool `json:"success"`
	// Optional message to pass on the challenger
	Message string `json:"message"`
}
