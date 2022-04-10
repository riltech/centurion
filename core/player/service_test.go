package player

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServiceIsPlayerExist(t *testing.T) {
	repo := NewRepository()
	service := NewService(repo)
	assert.False(t, service.IsPlayerExist(&Model{
		ID:    "xxx",
		Name:  "xxx",
		Team:  "xxxx",
		Score: 0,
	}))
	playa := Model{
		ID:    uuid.NewString(),
		Name:  "John",
		Team:  "attacker",
		Score: 0,
	}
	if err := repo.AddPlayer(playa); err != nil {
		assert.Nil(t, err)
	}
	assert.True(t, service.IsPlayerExist(&Model{
		ID:    "xxxx",
		Name:  "jOhn",
		Team:  "defender",
		Score: 0,
	}))
	assert.True(t, service.IsPlayerExist(&Model{
		ID:    "xxxx",
		Name:  "JOHN",
		Team:  "defender",
		Score: 0,
	}))
	assert.True(t, service.IsPlayerExist(&Model{
		ID:    playa.ID,
		Name:  "frankie",
		Team:  "defender",
		Score: 0,
	}))
}
