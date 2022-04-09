package engine

import (
	"fmt"
	"strings"

	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/player"
	"github.com/sirupsen/logrus"
)

// Describes the engine service interface
type IService interface {
	// Starts synchronous processing of the service towards the event bus
	Start()
	// Stops synchronous processing of the service towards the event bus
	Stop()
	// Returns if a given player already exists or not
	IsPlayerExist(*player.Model) bool
}

// Engine service implementation
type Service struct {
	// Layer dependencies
	bus        bus.IBus
	repository IRepository

	// Internal dependencies
	stop chan uint8
}

func (s *Service) Start() {
	for {
		select {
		case value := <-s.bus.Listen(bus.EventTypeRegistration):
			if value == nil {
				panic(fmt.Errorf("Received empty event for %s", bus.EventTypeRegistration))
			}
			s.routeEvent(*value)
		case <-s.stop:
			close(s.stop)
			return
		}
	}
}

func (s Service) Stop() {
	s.stop <- 1
}

// Responsible for registrating players
func (s Service) registration(e *bus.RegistrationEvent) error {
	if e == nil {
		return fmt.Errorf("Cannot process nil for registration")
	}
	return s.repository.AddPlayer(player.Model{
		ID:   e.ID,
		Name: e.Name,
		Team: e.Team,
	})
}

// Responsible for routing events
// NOTICE: Reference is used here just to make sure
// that we avoid unwanted accidental mutation
func (s Service) routeEvent(e bus.BusEvent) {
	switch e.Type {
	case bus.EventTypeRegistration:
		body, err := e.DecodeRegistration()
		if err != nil {
			logrus.Error(err)
			return
		}
		if err := s.registration(body); err != nil {
			panic(err)
		}
		return
	default:
		panic(fmt.Errorf("Unknown event %s", e.Type))
	}
}

func (s Service) IsPlayerExist(p *player.Model) bool {
	if p == nil {
		return false
	}
	for _, storedPlayer := range s.repository.GetPlayers() {
		if storedPlayer.ID == p.ID {
			return true
		}
		if strings.ToLower(storedPlayer.Name) == strings.ToLower(p.Name) {
			return true
		}
	}
	return false
}

// Interface check
var _ IService = (*Service)(nil)

// Constructor for the engine service
func NewService(bus bus.IBus, repo IRepository) IService {
	return &Service{
		bus,
		repo,
		make(chan uint8, 1),
	}
}
