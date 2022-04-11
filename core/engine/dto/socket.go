package dto

// Event type for a player joining the game
const SocketEventTypeError = "error"
const SocketEventTypeJoin = "join"
const SocketEventTypeAttack = "attack"
const SocketEventTypeAttackResult = "attack_result"
const SocketEventTypeAttackChallenge = "attack_challenge"
const SocketEventTypeAttackSolution = "attack_solution"

// Describes a generic event over websockets
type SocketEvent struct {
	Type string `json:"type"`
}

// Describes a player joined event
type JoinEvent struct {
	SocketEvent
	ID string `json:"id"`
}

// Happens when a new attack is launched
type AttackEvent struct {
	SocketEvent
	// Describes the ID of the challenge
	TargetID string `json:"targetId"`
}

// Happens when an attack is for a valid target ID
// The challenge creator generates hints to resolve
type AttackChallengeEvent struct {
	SocketEvent
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
	TargetID string `json:"targetId"`
	Success  bool   `json:"success"`
}

// Sent when an error happens during an action
type ErrorEvent struct {
	SocketEvent
	Message string `json:"message"`
}
