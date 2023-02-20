package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetTileErrors(t *testing.T) {
	f := NewFastBoard()
	assert.EqualError(t, f.setTile(3, 0, 1, 1), "invalid team chosen")
	assert.EqualError(t, f.setTile(1, 5, 1, 1), "invalid height chosen")
	assert.EqualError(t, f.setTile(1, 4, 1, 1), "cannot set cap with a worker present")
	assert.Equal(t, nil, f.setTile(1, 3, 2, 4)) // put a team 1 worker at 2, 4
	assert.EqualError(t, f.setTile(2, 3, 2, 4), "position is occupied")
	assert.EqualError(t, f.setTile(0, 3, 6, 4), "bad x")
	assert.EqualError(t, f.setTile(0, 3, 2, 7), "bad y")
}

func TestGetTile(t *testing.T) {
	f := NewFastBoard()
	assert.Equal(t, nil, f.setTile(1, 3, 2, 4)) // put a team 1 worker at 2, 4
	tile := f.GetTile(2, 4)
	assert.True(t, tile.GetHeight() == 3)
	assert.True(t, tile.GetTeam() == 1)
	assert.True(t, tile.GetX() == 2)
	assert.True(t, tile.GetY() == 4)
	assert.Panics(t, func() {
		f.GetTile(6, 1)
	}, "should be an invalid x")
	assert.Panics(t, func() {
		f.GetTile(1, 7)
	}, "should be an invalid y")
}

func TestGetMoves(t *testing.T) {
	f := NewFastBoard()
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			// put a work in the corner
			assert.Nil(t, f.setTile(1, 0, uint8(x*4), uint8(y*4)), "placing worker on empty board")
			assert.Equal(t, 18, len(f.ValidTurns(1)), "Corners have 18 valid moves (%d, %d)", x*4, y*4)
			assert.Nil(t, f.setTile(0, 0, uint8(x*4), uint8(y*4)), "removing worker")
		}
	}

	// Place a worker in a trap and see if it can move
	assert.Nil(t, f.setTile(1, 0, 2, 2), "placing worker on empty board")
	for _, c := range InnerRing {
		assert.Nil(t, f.setTile(0, 4, c[0], c[1]), "placing cap on empty board (%d,%d)", c[0], c[1])
	}
	assert.Equal(t, 0, len(f.ValidTurns(1)), "Trapper players cannot move")
	// reset
	f = NewFastBoard()

	// Place a worker in a trap with one way out see if it can move
	assert.Nil(t, f.setTile(1, 0, 2, 2), "placing worker on empty board")
	for _, c := range InnerRing[:len(InnerRing)-1] {
		assert.Nil(t, f.setTile(0, 2, c[0], c[1]), "placing double block empty board (%d,%d)", c[0], c[1])
	}
	turns := f.ValidTurns(1)
	assert.Equal(t, 8, len(turns), "Trapper player only has 8 moves")

	// reset
	f = NewFastBoard()

	// Place a worker next to another worker to make sure that limits his moves
	assert.Nil(t, f.setTile(1, 0, 0, 0), "placing worker on empty board")
	assert.Nil(t, f.setTile(2, 0, 1, 1), "placing worker on empty board")
	assert.Equal(t, 8, len(f.ValidTurns(1)), "Player should only have 8 moves")
}
