package player

import (
	"fmt"
	"strings"
	"sync"
)

// Describes a repository for the engine
type IRepository interface {
	// Fetches the available players in the system
	GetPlayers() []Model
	// Adds a new player to the system
	AddPlayer(Model) error
	// Update by ID
	UpdateByID(ID string, update Model) (Model, error)
	// Finds a user by ID
	FindByID(ID string) (Model, error)
	// Add given amounts of points to a player
	AddPoint(ID string, points int) (Model, error)
}

// Engine repository implementation
type Repository struct {
	mux     sync.RWMutex
	players []Model
}

// Interface check
var _ IRepository = (*Repository)(nil)

func (r *Repository) AddPlayer(user Model) error {
	if r == nil {
		return fmt.Errorf("Repository needs to be initialised before usage")
	}
	if r.players == nil {
		r.mux.Lock()
		r.players = []Model{user}
		r.mux.Unlock()
		return nil
	}
	r.mux.RLock()
	for _, p := range r.players {
		if p.ID == user.ID || strings.ToLower(p.Name) == strings.ToLower(user.Name) {
			r.mux.RUnlock()
			return fmt.Errorf("Cannot use the same id (expected, got) (%s, %s) or username (%s, %s)", user.ID, p.ID, user.Name, p.Name)
		}
	}
	r.mux.RUnlock()
	r.mux.Lock()
	r.players = append(r.players, user)
	r.mux.Unlock()
	return nil
}

func (r *Repository) GetPlayers() []Model {
	if r == nil {
		return []Model{}
	}
	defer r.mux.RUnlock()
	r.mux.RLock()
	return r.players
}

func (r *Repository) UpdateByID(ID string, update Model) (Model, error) {
	if r == nil {
		return Model{}, fmt.Errorf("Cannot update without the repository being initialised")
	}
	r.mux.RLock()
	for i, user := range r.players {
		if user.ID == ID {
			r.mux.RUnlock()
			r.mux.Lock()
			r.players[i] = update
			r.mux.Unlock()
			return update, nil
		}
	}
	r.mux.RUnlock()
	return Model{}, fmt.Errorf("%s not found in players", ID)
}

func (r *Repository) FindByID(ID string) (Model, error) {
	if r == nil {
		return Model{}, fmt.Errorf("Repository is not initialised")
	}
	r.mux.RLock()
	defer r.mux.RUnlock()
	for _, p := range r.players {
		if p.ID == ID {
			return p, nil
		}
	}
	return Model{}, fmt.Errorf("%s user not found", ID)
}

func (r *Repository) AddPoint(ID string, points int) (Model, error) {
	if r == nil {
		return Model{}, fmt.Errorf("Repository is not initialised")
	}
	r.mux.Lock()
	defer r.mux.Unlock()
	for i, p := range r.players {
		if p.ID == ID {
			p.Score = p.Score + points
			r.players[i].Score = p.Score
			return p, nil
		}
	}
	return Model{}, fmt.Errorf("%s player not found", ID)
}

// Constructor to create a new engine repository
func NewRepository() *Repository {
	return &Repository{
		mux: sync.RWMutex{},
	}
}
