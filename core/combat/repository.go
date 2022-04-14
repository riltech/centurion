package combat

import (
	"fmt"
	"sync"
	"time"

	"github.com/riltech/centurion/core/logger"
)

// Describes a repository for the combat package
type IRepository interface {
	// Fetches the available combats in the system
	GetCombats() []Model
	// Adds a new combat to the system
	AddCombat(Model) error
	// Finds a combat by ID
	FindByID(ID string) (Model, error)
	// Updates combat state
	UpdateCombatState(ID string, state string) (Model, error)
	// Returns the archive (aka finished events)
	GetArchive() []Model
}

// Combat repository implementation
type Repository struct {
	mux     sync.RWMutex
	combats []Model
	archive []Model
}

// Interface check
var _ IRepository = (*Repository)(nil)

func (r *Repository) AddCombat(combat Model) error {
	if r == nil {
		return fmt.Errorf("Repository needs to be initialised before usage")
	}
	r.mux.Lock()
	defer r.mux.Unlock()
	combat.CreatedAt = time.Now()
	if r.combats == nil {
		r.combats = []Model{combat}
		return nil
	}
	for _, c := range r.combats {
		if c.ID == combat.ID {
			return fmt.Errorf("Cannot use the same id (expected, got) (%s, %s)", combat.ID, c.ID)
		}
	}
	r.combats = append(r.combats, combat)
	return nil
}

func (r *Repository) GetCombats() []Model {
	defer r.mux.RUnlock()
	r.mux.RLock()
	return r.combats
}

func (r *Repository) FindByID(ID string) (Model, error) {
	if r == nil {
		return Model{}, fmt.Errorf("Repository is not initialised")
	}
	r.mux.RLock()
	defer r.mux.RUnlock()
	for _, c := range r.combats {
		if c.ID == ID {
			return c, nil
		}
	}
	return Model{}, fmt.Errorf("%s combat not found", ID)
}

// moves a combat from the active ones to the archive
// NOTE: This function is not thread safe
// Using this function requires the mux already being locked
func (r *Repository) archiveElement(index int) {
	if r == nil {
		return
	}
	if len(r.combats)-1 > index {
		logger.LogError(fmt.Errorf("Could not archieve %d combat because len is %d", index, len(r.combats)))
		return
	}
	r.archive = append(r.archive, r.combats[index])
	r.combats = append(r.combats[:index], r.combats[index+1:]...)
}

func (r *Repository) UpdateCombatState(ID string, state string) (Model, error) {
	if r == nil {
		return Model{}, fmt.Errorf("Repository is not initialised")
	}
	r.mux.Lock()
	defer r.mux.Unlock()
	for i, c := range r.combats {
		if c.ID == ID {
			c.CombatState = state
			r.combats[i].CombatState = state
			r.combats[i].LastUpdateAt = time.Now()
			if IsFinalCombatState(state) {
				r.archiveElement(i)
			}
			return c, nil
		}
	}
	return Model{}, fmt.Errorf("%s not found", ID)
}

func (r *Repository) GetArchive() []Model {
	if r == nil {
		return nil
	}
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.archive
}

// Constructor to create a new engine repository
func NewRepository() *Repository {
	return &Repository{
		mux: sync.RWMutex{},
	}
}
