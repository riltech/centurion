package challenge

import (
	"fmt"
	"strings"
	"sync"
)

// Describes a repository for the engine
type IRepository interface {
	// Fetches the available challenges in the system
	GetChallenges() []Model
	// Adds a new challenge
	AddChallenge(Model) error
}

// Challenge repository implementation
type Repository struct {
	mux        sync.RWMutex
	challenges []Model
}

// Interface check
var _ IRepository = (*Repository)(nil)

func (r *Repository) AddChallenge(challenge Model) error {
	if r == nil {
		return fmt.Errorf("Repository needs to be initialised before usage")
	}
	if r.challenges == nil {
		r.mux.Lock()
		r.challenges = []Model{challenge}
		r.mux.Unlock()
		return nil
	}
	for _, c := range r.challenges {
		if c.ID == challenge.ID || strings.ToLower(c.Name) == strings.ToLower(challenge.Name) {
			return fmt.Errorf("Cannot use the same id (expected, got) (%s, %s) or challenge title (%s, %s)", challenge.ID, c.ID, challenge.Name, c.Name)
		}
	}
	r.mux.Lock()
	r.challenges = append(r.challenges, challenge)
	r.mux.Unlock()
	return nil
}

func (r *Repository) GetChallenges() []Model {
	defer r.mux.RUnlock()
	r.mux.RLock()
	return r.challenges
}

// Constructor to create a new engine repository
func NewRepository() *Repository {
	return &Repository{}
}
