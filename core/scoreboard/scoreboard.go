package scoreboard

// Describes a scoreboard in the system
type Model struct {
	// Either 'attacker' or 'defender'
	// Use Team enums from player package
	Team string
	// Overall score of the team
	OverallScore int
}
