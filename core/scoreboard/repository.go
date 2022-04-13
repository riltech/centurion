package scoreboard

import (
	"sync"

	"github.com/riltech/centurion/core/player"
)

// Describes a repository for the engine
type IRepository interface {
	// Returns the score board for each team
	GetBoards() (attacker Model, defender Model)
	// Adds a given point to a team
	// Use enums from player package to [team]
	AddPoint(team string, point int)
}

// Combat repository implementation
type Repository struct {
	mux      sync.RWMutex
	attacker Model
	defender Model
}

// Interface check
var _ IRepository = (*Repository)(nil)

func NewRepository() IRepository {
	return &Repository{
		mux: sync.RWMutex{},
		attacker: Model{
			Team:         player.TeamTypeAttacker,
			OverallScore: 0,
		},
		defender: Model{
			Team:         player.TeamTypeDefender,
			OverallScore: 0,
		},
	}
}

func (r *Repository) GetBoards() (Model, Model) {
	if r == nil {
		return Model{}, Model{}
	}
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.attacker, r.defender
}

func (r *Repository) AddPoint(team string, points int) {
	if r == nil {
		return
	}
	r.mux.Lock()
	defer r.mux.Unlock()
	if team == player.TeamTypeAttacker {
		r.attacker.OverallScore = r.attacker.OverallScore + points
		return
	}
	if team == player.TeamTypeDefender {
		r.defender.OverallScore = r.defender.OverallScore + points
		return
	}
}
