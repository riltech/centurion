package challenge

// Describes a default module challenge which is part of the starting
// collection when the game launches
const ChallengeTypeDefault = "default"

// Describes any other challenges created by players
const ChallengeTypePlayerCreated = "player_created"

// Describes a challenge
type Model struct {
	ID          string
	CreatorID   string
	Name        string
	Description string
	Type        string
	Example     Example
}

// Example of a challenge
type Example struct {
	Hints     []interface{}
	Solutions []interface{}
}
