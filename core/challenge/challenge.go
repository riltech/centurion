package challenge

import "time"

// Describes a default module challenge which is part of the starting
// collection when the game launches
const ChallengeTypeDefault = "default"

// Describes any other challenges created by players
const ChallengeTypePlayerCreated = "player_created"

// Describes a challenge
type Model struct {
	// ID of the challenge
	ID string
	// ID of the creator
	CreatorID string
	// Name of the challenge
	Name string
	// Description of the challenge
	Description string
	// Internal type (see consts for types)
	Type string
	// Creation time
	CreatedAt time.Time
	// Example resolution
	Example Example
}

// Example of a challenge
type Example struct {
	// Hints array
	Hints []interface{}
	// Solutions array
	Solutions []interface{}
}
