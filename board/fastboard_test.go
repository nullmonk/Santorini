package board

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randUint8(max int64) uint8 {
	b, _ := rand.Int(rand.Reader, big.NewInt(max))
	return uint8(b.Int64())
}

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

func TestSetGet(t *testing.T) {
	f := NewFastBoard()
	for i := 0; i < 1000; i++ {
		height := randUint8(3)
		team := randUint8(3)
		x := randUint8(4)
		y := randUint8(4)
		assert.Nil(t, f.setTile(team, height, x, y), "error setTile(%d, %d, %d, %d)", team, height, x, y)
		tile := f.GetTile(x, y)
		assert.Equal(t, height, tile.height)
		assert.Equal(t, team, tile.team)
		assert.Equal(t, x, tile.x)
		assert.Equal(t, y, tile.y)
		assert.Nil(t, f.setTile(0, 0, x, y), "error setTile(0, 0, %d, %d)", x, y)
	}
}

func testGetTile(t *testing.T) {
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

func TestBuilding(t *testing.T) {
	f := NewFastBoard()
	assert.Nil(t, f.setTile(1, 0, 2, 1), "placing worker on empty board")
	tile := f.GetTile(2, 1)
	//fmt.Printf("here %d %08s\n", tile.team, strconv.FormatUint(uint64(f.board[(f.width*0)+0]), 2))
	assert.Equal(t, uint8(1), tile.team, "Team should be set")
	candidates := f.ValidTurns(1)
	assert.Equal(t, 55, len(candidates), "did not return all the expected turns")
	for _, turn := range candidates {
		f.board = make([]uint8, 25)
		assert.Nil(t, f.setTile(1, 0, 2, 1), "placing worker on empty board")
		_, err := f.PlayTurn(turn)
		assert.Nil(t, err, "failed to take turn %s", turn)
		t2 := f.GetTile(turn.Build.x, turn.Build.y)
		assert.Equal(t, turn.Build.height+1, t2.height, "Height should increase")
		assert.Equal(t, uint8(0), t2.team, "Build should not have a team")
		t2 = f.GetTile(turn.MoveTo.x, turn.MoveTo.y)
		assert.Equal(t, turn.MoveTo.height, t2.height, "worker height should be the same")
		assert.Equal(t, turn.Worker.team, t2.team, "MoveTo should have a worker")
	}
}
func TestUndoSimple(t *testing.T) {
	f := NewFastBoard()
	assert.Nil(t, f.setTile(1, 0, 1, 2), "placing worker on empty board")
	candidates := f.ValidTurns(1)
	assert.Equal(t, 55, len(candidates), "did not return all the expected turns")
	for i, turn := range candidates {
		assert.Equal(t, uint8(1), turn.Worker.team)
		origHash := f.Hash()
		_, err := f.PlayTurn(turn)
		err2 := f.UndoTurn(turn)
		assert.Equal(t, origHash, f.Hash(), "[Test %d] board hashes do not match: %s != %s", i, f.Hash(), origHash)
		assert.Nil(t, err, "failed to take turn %s", turn)
		assert.Nil(t, err2, "failed to undo turn %s", turn)
		if err != nil || err2 != nil {
			break
		}
	}
}
func TestUndo(t *testing.T) {
	f := NewFastBoard()
	count := 10
	// Keep track of each turn that we make
	turns := make([]*Turn, 0, count)
	assert.Nil(t, f.setTile(1, 0, 2, 2), "placing worker on empty board")
	startHsh := f.Hash()
	for i := 0; i < count; i++ {
		// Get the possible turns for the player
		candidates := f.ValidTurns(1)
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates)-1)))
		if err != nil {
			t.Fatal(err)
		}
		turn := candidates[n.Int64()]
		// play that turn
		over, err := f.PlayTurn(turn)
		if err != nil {
			t.Fatal(err)
		}
		turns = append(turns, turn)
		if over {
			break
		}
	}
	// Undo all the turns
	for i := len(turns) - 1; i >= 0; i-- {
		err := f.UndoTurn(turns[i])
		assert.Nil(t, err, "Failed to undo turn %s", turns[i])

	}
	assert.Equal(t, startHsh, f.Hash())
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
