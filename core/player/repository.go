package player

import (
	"fmt"
	"strings"
)

// Describes a repository for the engine
type IRepository interface {
	// Fetches the available players in the system
	GetPlayers() []Model
	// Adds a new player to the system
	AddPlayer(Model) error
}

// Engine repository implementation
type Repository struct {
	players []Model
}

// Interface check
var _ IRepository = (*Repository)(nil)

func (r *Repository) AddPlayer(user Model) error {
	if r == nil {
		return fmt.Errorf("Repository needs to be initialised before usage")
	}
	if r.players == nil {
		r.players = []Model{user}
		return nil
	}
	for _, p := range r.players {
		if p.ID == user.ID || strings.ToLower(p.Name) == strings.ToLower(user.Name) {
			return fmt.Errorf("Cannot use the same id (expected, got) (%s, %s) or username (%s, %s)", user.ID, p.ID, user.Name, p.Name)
		}
	}
	r.players = append(r.players, user)
	return nil
}

func (r Repository) GetPlayers() []Model {
	return r.players
}

// Constructor to create a new engine repository
func NewRepository() *Repository {
	return &Repository{}
}
