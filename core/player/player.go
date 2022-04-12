package player

// Describes attacker team type
const TeamTypeAttacker = "attacker"

// Describes defender team type
const TeamTypeDefender = "defender"

// Describes a player in the system
// who is a user already registered to a given team
// and (possibly) actively participating in the game
type Model struct {
	// ID of the player
	ID string
	// Name of the player
	Name string
	// Team of the player. Either "attacker" or "defender"
	Team string
	// Current score of the individual player
	Score int
	// Online indicator
	Online bool
}
