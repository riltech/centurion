package engine

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/player"
	"github.com/stretchr/testify/assert"
)

func TestServiceIsPlayerExist(t *testing.T) {
	bus := bus.NewBus()
	repo := NewRepository()
	service := NewService(bus, repo)
	assert.False(t, service.IsPlayerExist(&player.Model{
		ID:    "xxx",
		Name:  "xxx",
		Team:  "xxxx",
		Score: 0,
	}))
	playa := player.Model{
		ID:    uuid.NewString(),
		Name:  "John",
		Team:  "attacker",
		Score: 0,
	}
	if err := repo.AddPlayer(playa); err != nil {
		assert.Nil(t, err)
	}
	assert.True(t, service.IsPlayerExist(&player.Model{
		ID:    "xxxx",
		Name:  "jOhn",
		Team:  "defender",
		Score: 0,
	}))
	assert.True(t, service.IsPlayerExist(&player.Model{
		ID:    "xxxx",
		Name:  "JOHN",
		Team:  "defender",
		Score: 0,
	}))
	assert.True(t, service.IsPlayerExist(&player.Model{
		ID:    playa.ID,
		Name:  "frankie",
		Team:  "defender",
		Score: 0,
	}))
}

func TestServiceRegistration(t *testing.T) {
	eventBus := bus.NewBus()
	repo := NewRepository()
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
