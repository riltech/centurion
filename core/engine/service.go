package engine

import (
	"fmt"

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
}

// Engine service implementation
type Service struct {
	// Layer dependencies
	bus              bus.IBus
	playerRepository player.IRepository

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
	return s.playerRepository.AddPlayer(player.Model{
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
		body, err := e.DecodeRegistrationEvent()
		if err != nil {
			logrus.Error(err)
			return
		}
		if err := s.registration(body); err != nil {
			logrus.Error(err)
		}
		return
	default:
		logrus.Error(fmt.Errorf("Unknown event %s", e.Type))
	}
}

// Interface check
var _ IService = (*Service)(nil)

// Constructor for the engine service
func NewService(bus bus.IBus, repo player.IRepository) IService {
	return &Service{
		bus,
		repo,
		make(chan uint8, 1),
	}
}
