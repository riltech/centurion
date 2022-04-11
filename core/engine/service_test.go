package engine

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialization(t *testing.T) {
	const input = `{"id": "25", "name": "asd"}`

	type base struct {
		Name string `json:"name"`
	}
	type extension struct {
		base
		ID string `json:"id"`
	}
	var b base
	var e extension
	assert.Nil(t, json.Unmarshal([]byte(input), &b))
	assert.Nil(t, json.Unmarshal([]byte(input), &e))
	assert.Equal(t, b.Name, "asd")
	assert.Equal(t, e.Name, "asd")
	assert.Equal(t, e.ID, "25")
}
