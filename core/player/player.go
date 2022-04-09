package player

// Describes a player in the system
// who is a user already registered to a given team
// and (possibly) actively participating in the game
type Model struct {
	ID    string
	Name  string
	Team  string
	Score int
}
