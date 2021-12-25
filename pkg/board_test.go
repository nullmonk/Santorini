package santorini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			tile := board.GetTile(uint8(x), uint8(y))
			assert.Equal(t, uint8(x), tile.x)
			assert.Equal(t, uint8(y), tile.y)
			assert.Equal(t, 0, tile.Height)
			assert.Nil(t, tile.Worker)
		}
	}
}
