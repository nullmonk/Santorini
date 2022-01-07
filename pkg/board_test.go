package santorini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			tile := board.GetTile(x, y)
			assert.Equal(t, x, tile.x)
			assert.Equal(t, y, tile.y)
			assert.Equal(t, 0, tile.height)
			assert.Equal(t, 0, tile.team)
			assert.Equal(t, 0, tile.worker)
			assert.False(t, tile.IsCapped())
			assert.False(t, tile.IsOccupied())
		}
	}
}

func TestPlayTurn(t *testing.T) {
	board := NewBoard()

	board.PlaceWorker(1, 1, 1, 2)

	tile := board.GetTile(1, 2)
	assert.Equal(t, 1, tile.x)
	assert.Equal(t, 2, tile.y)
	assert.Equal(t, 1, tile.worker)
	assert.Equal(t, 1, tile.team)

	board.PlayTurn(Turn{
		Team:   1,
		Worker: 1,
		MoveTo: Tile{x: 2, y: 2},
		Build:  Tile{x: 1, y: 2},
	})

	oldTile := board.GetTile(1, 2)
	assert.Equal(t, 0, oldTile.team)
	assert.Equal(t, 0, oldTile.worker)
	assert.Equal(t, 1, oldTile.x)
	assert.Equal(t, 2, oldTile.y)
	assert.Equal(t, 1, oldTile.height)

	newTile := board.GetTile(2, 2)
	assert.Equal(t, 1, newTile.team)
	assert.Equal(t, 1, newTile.worker)
	assert.Equal(t, 2, newTile.x)
	assert.Equal(t, 2, newTile.y)
	assert.Equal(t, 0, newTile.height)
}
