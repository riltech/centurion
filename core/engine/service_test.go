package engine

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/player"
	"github.com/stretchr/testify/assert"
)

func TestServiceRegistration(t *testing.T) {
	eventBus := bus.NewBus()
	repo := player.NewRepository()
	service := NewService(eventBus, repo)

	go service.Start()

	information := bus.RegistrationEvent{
		Name: "John",
		ID:   uuid.NewString(),
		Team: "defender",
	}
	eventBus.Send(&bus.BusEvent{
		Type:        bus.EventTypeRegistration,
		Information: information,
	})
	<-time.After(500 * time.Millisecond)
	players := repo.GetPlayers()
	assert.True(t, len(players) == 1)
	playa := players[0]
	assert.Equal(t, information.ID, playa.ID)
	assert.Equal(t, information.Name, playa.Name)
	assert.Equal(t, information.Team, playa.Team)
}
