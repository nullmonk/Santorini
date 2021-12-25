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

func TestIsOccupied(t *testing.T) {
	workerOne := Worker{Team: 1, Number: 1}
	workerTwo := Worker{Team: 1, Number: 2}
	workerThree := Worker{Team: 2, Number: 1}
	workerFour := Worker{Team: 2, Number: 2}

	unoccupied := Tile{}
	assert.False(t, unoccupied.IsOccupied())
	assert.False(t, unoccupied.IsOccupiedBy(workerOne))
	assert.False(t, unoccupied.IsOccupiedBy(workerTwo))
	assert.False(t, unoccupied.IsOccupiedBy(workerThree))
	assert.False(t, unoccupied.IsOccupiedBy(workerFour))

	occupied := Tile{Worker: &workerOne}
	assert.True(t, occupied.IsOccupied())
	assert.True(t, occupied.IsOccupiedBy(workerOne))
	assert.False(t, unoccupied.IsOccupiedBy(workerTwo))
	assert.False(t, unoccupied.IsOccupiedBy(workerThree))
	assert.False(t, unoccupied.IsOccupiedBy(workerFour))
}
