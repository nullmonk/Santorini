package santorini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTileIsOccupied(t *testing.T) {
	unoccupied := Tile{}
	occupied := Tile{team: 1, worker: 1}

	assert.False(t, unoccupied.IsOccupied())
	assert.True(t, occupied.IsOccupied())
	assert.True(t, occupied.IsOccupiedBy(1, 1))
}
func TestTileIsCapped(t *testing.T) {
	uncapped := Tile{}
	capped := Tile{height: 4}

	assert.False(t, uncapped.IsCapped())
	assert.True(t, capped.IsCapped())
}
