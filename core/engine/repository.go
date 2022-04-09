package engine

import (
	"fmt"

	"github.com/riltech/centurion/core/player"
)

// Describes a repository for the engine
type IRepository interface {
	// Fetches the available players in the system
	GetPlayers() []player.Model
	// Adds a new player to the system
	AddPlayer(player.Model) error
}

// Engine repository implementation
type Repository struct {
	players []player.Model
}

// Interface check
var _ IRepository = (*Repository)(nil)

func (r *Repository) AddPlayer(user player.Model) error {
	if r == nil {
		return fmt.Errorf("Repository needs to be initialised before usage")
	}
	if r.players == nil {
		r.players = []player.Model{user}
		return nil
	}
	r.players = append(r.players, user)
	return nil
}

func (r Repository) GetPlayers() []player.Model {
	return r.players
}

// Constructor to create a new engine repository
func NewRepository() *Repository {
	return &Repository{}
}
