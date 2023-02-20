package board

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Two blocks from the middle
var OuterRing = [][]uint8{
	{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0},
	{0, 1} /*                   */, {4, 1},
	{0, 2} /*                   */, {4, 2},
	{0, 3} /*                   */, {4, 3},
	{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4},
}

// 1 block from the middle
var InnerRing = [][]uint8{
	{1, 1}, {2, 1}, {3, 1},
	{1, 2} /*   */, {3, 2},
	{1, 3}, {2, 3}, {3, 3},
}

func TestTileIsOccupied(t *testing.T) {
	unoccupied := Tile{}
	occupied := Tile{team: 1}

	assert.True(t, unoccupied.GetTeam() == 0)
	assert.True(t, occupied.GetTeam() == 1)
}
func TestTileIsCapped(t *testing.T) {
	uncapped := Tile{}
	capped := Tile{height: 4}

	assert.False(t, uncapped.IsCapped())
	assert.True(t, capped.IsCapped())
}

func TestTileMovement(t *testing.T) {
	// Worker normal: Move down
	assert.Equal(t, nil, Tile{1, 0, 2, 1}.CanMoveTo(Tile{0, 0, 2, 2}))

	// Worker cannot move where another worker is
	assert.EqualError(t,
		Tile{1, 0, 2, 1}.CanMoveTo(Tile{2, 0, 2, 2}),
		"the worker cannot move to an occupied block",
		"the worker cannot move to where another worker is",
	)
	assert.EqualError(t,
		Tile{1, 0, 2, 1}.CanMoveTo(Tile{0, 2, 2, 2}),
		"the worker cannot jump 2 blocks",
	)
	assert.Equal(t, nil,
		Tile{1, 3, 2, 1}.CanMoveTo(Tile{0, 0, 2, 2}),
		"the worker CAN jump down",
	)

	assert.EqualError(t,
		Tile{1, 3, 2, 1}.CanMoveTo(Tile{0, 4, 2, 2}),
		"the worker cannot move to an occupied block",
		"the worker cannot move to a capped block",
	)
	assert.EqualError(t,
		Tile{1, 0, 0, 0}.CanMoveTo(Tile{0, 1, 0, 0}),
		"the worker cannot move to the given block",
		"the worker cannot move to itself",
	)

	src := Tile{1, 1, 2, 2}
	for _, dst := range OuterRing {
		assert.EqualError(t,
			src.CanMoveTo(Tile{0, 1, dst[0], dst[1]}),
			"the worker cannot move to the given block",
			"Out of range",
		)
	}
	for _, dst := range InnerRing {
		assert.Equal(t, nil,
			src.CanMoveTo(Tile{0, 1, dst[0], dst[1]}),
			"The worker should be able to go here",
		)
	}
}
func TestTileJson(t *testing.T) {
	tile := Tile{
		2,
		2,
		4,
		1,
	}
	b, err := json.Marshal(&tile)
	if err != nil {
		t.Error(err)
	}
	tile2 := Tile{}
	if err = json.Unmarshal(b, &tile2); err != nil {
		t.Error(err)
	}
	assert.True(t, tile.GetHeight() == tile2.GetHeight())
	assert.True(t, tile.GetTeam() == tile2.GetTeam())
	assert.True(t, tile.GetX() == tile2.GetX())
	assert.True(t, tile.GetY() == tile2.GetY())
}
