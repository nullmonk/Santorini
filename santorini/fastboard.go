package santorini

import (
	"fmt"
	"math"
	"strings"
)

// FastBoard uses a single array as the storage mechanism
type FastBoard struct {
	// Last three bits are the tile height, rest is team number
	board     []uint8
	width     uint8
	height    uint8
	turnCount uint // number of turns that have been taken on this board (used to calculate which teams turn it is)

}

// Create a new Fastboard with the default layout
func NewBoard(options ...func(*FastBoard) *FastBoard) *FastBoard {
	f := &FastBoard{
		board:  make([]uint8, 25),
		width:  5,
		height: 5,
	}

	if len(options) == 0 {
		options = append(options, Default2Player)
	}

	var r *FastBoard
	for _, o := range options {
		r = o(f)
	}

	return r
}

func NewBoardFromHash(hash string) (*FastBoard, error) {
	width := hash[0] - 65 + 3
	size := len(hash) - 1
	b := &FastBoard{
		width:  width,
		height: uint8(size / int(width)),
		board:  make([]uint8, size),
	}
	for i, v := range hash[1:] {
		b.board[i] = uint8(v) - 65
	}
	return b, nil
}

func (f *FastBoard) Dimensions() (uint8, uint8) {
	return f.width, f.height
}

func (f *FastBoard) Teams() []uint8 {
	teams := make([]uint8, 0, 2)
	for _, i := range f.board {
		teams = append(teams, i>>3)
	}
	return teams
}

/*
The boardhash is an ascii string that contains state of the game, including what is on each tile
and the size of the board
*/
func (f *FastBoard) GameHash() string {
	res := new(strings.Builder)
	// Save the board (we only need the w as height = size/w)
	w, _ := f.Dimensions()
	if w > 29 {
		panic("cannot hash a board larger than 29x29")
	}
	if w < 3 {
		panic("cannot hash a board smaller than 3x3")
	}

	// Stash the width, we can calc height by len of hash and width
	res.WriteByte(65 + w - 3)
	// Save the board data
	for _, i := range f.board {
		res.WriteByte(i + 65)
	}
	return res.String()
}

/**** Game Flow Functions ****/ //
func (f *FastBoard) Clone() *FastBoard {
	b := &FastBoard{
		board:     make([]uint8, len(f.board)),
		width:     f.width,
		height:    f.height,
		turnCount: f.turnCount,
	}
	copy(b.board, f.board)
	return b
}

// Undo a move
func (f *FastBoard) UndoTurn(t *Turn) error {
	// Not many other checks I can do to validate that the game hasnt changed...
	/*
		if actual.MoveTo.team != t.Worker.team {
			return fmt.Errorf("cannot undo, worker is not in expected location")
		}
	*/

	// Undo build at the Build position
	curHeight := f.board[(f.width*t.Build.y)+t.Build.x] & 0x7
	if curHeight < 1 {
		return fmt.Errorf("cannot undo, no building at %c%d", rune(t.Build.x+65), t.Build.y)
	}
	f.board[(f.width*t.Build.y)+t.Build.x] = curHeight - 1

	// Set the worker back to the original position
	curHeight = f.board[(f.width*t.Worker.y)+t.Worker.x] & 0x7
	f.board[(f.width*t.Worker.y)+t.Worker.x] = (t.Worker.team << 3) | curHeight

	// Set the MoveTo to team 0
	f.board[(f.width*t.MoveTo.y)+t.MoveTo.x] = t.MoveTo.height

	// Validate none of the blocks has changed from the expected
	want := t.Worker.team<<3 | t.Worker.height
	have := f.board[(f.width*t.Worker.y)+t.Worker.x]
	if have != want {
		return fmt.Errorf("worker block has changed: expected %v, found %v", want, have)
	}
	f.turnCount--
	return nil
}
func (f *FastBoard) PlayTurn(t *Turn) (victory bool, err error) {
	// Get the turn values from our board bc we dont trust user input
	turn := Turn{
		Worker: f.GetTile(t.Worker.x, t.Worker.y),
		MoveTo: f.GetTile(t.MoveTo.x, t.MoveTo.y),
		Build:  f.GetTile(t.Build.x, t.Build.y),
	}
	// Make sure the worker is good
	if turn.Worker.team == 0 {
		return false, fmt.Errorf("no worker at %d,%d", t.Worker.x, t.Worker.y)
	}

	// TODO: Do we force workers to stop checks? i.e. dont let them make a move that would cause them to lose?
	// If so, we need to implement that logic here

	// Make sure the move is valid
	if err := turn.Worker.CanMoveTo(turn.MoveTo); err != nil {
		return false, err
	}

	// cant build on an occupied spot, UNLESS that spot is where the worker currently is
	if turn.Build.IsOccupied() && !turn.Build.SameLocation(t.Worker) {
		return false, fmt.Errorf("the worker cannot build on the given block")
	}
	if turn.Build.SameLocation(turn.MoveTo) {
		return false, fmt.Errorf("the worker cannot build on the given block")
	}

	// Clear the current worker position
	if err := f.setTile(0, turn.Worker.height, turn.Worker.x, turn.Worker.y); err != nil {
		return false, fmt.Errorf("error clearing worker at (%d,%d) : %s", turn.Worker.x, turn.Worker.y, err)
	}

	// Move the worker to the new position
	if err := f.setTile(turn.Worker.team, turn.MoveTo.height, turn.MoveTo.x, turn.MoveTo.y); err != nil {
		return false, fmt.Errorf("error moving worker at (%d,%d) : %s", turn.MoveTo.x, turn.MoveTo.y, err)
	}

	// if we have won, just end it here
	if turn.MoveTo.height == 3 {
		return true, nil
	}

	// TODO: Calculate more complex win scenarios here (e.g. Checkmate, etc. so the game is instantly over)

	// Build at the new position
	if err := f.setTile(0, turn.Build.height+1, turn.Build.x, turn.Build.y); err != nil {
		return false, fmt.Errorf("error building at (%d,%d) : %s", turn.Build.x, turn.Build.y, err)
	}

	for _, b := range f.board {
		height := b & 0x7
		team := b >> 3
		if height == 3 && team != 0 {
			return true, nil
		}
	}

	f.turnCount++
	return false, nil
}

func (f *FastBoard) GetTile(x, y uint8) Tile {
	if x >= f.width {
		panic(fmt.Errorf("invalid x"))
	}
	if y >= f.height {
		panic(fmt.Errorf("invalid y"))
	}
	index := f.board[(f.width*y)+x]
	return Tile{
		team:   index >> 3,
		height: index & 0x7,
		x:      x,
		y:      y,
	}
}

func (f *FastBoard) setTile(team, height, x, y uint8) error {
	if height > 4 {
		return fmt.Errorf("invalid height chosen")
	}
	if height == 4 && team > 0 {
		return fmt.Errorf("cannot set cap with a worker present")
	}
	if x >= f.width {
		return fmt.Errorf("bad x")
	}
	if y >= f.height {
		return fmt.Errorf("bad y")
	}
	current := f.board[(f.width*y)+x]
	if current>>3 > 0 && team > 0 {
		return fmt.Errorf("position is occupied")
	}
	f.board[(f.width*y)+x] = (team << 3) | height
	return nil
}

/**** Bot Functions ****/ //
func (f *FastBoard) GetWorkers(team uint8) []Tile {
	workers := make([]Tile, 0, 2)
	// find the workers
	for i, tile := range f.board {
		y := uint8(i) / f.width
		x := uint8(i) - (f.width * y)
		// Not the team we are looking for
		if tile>>3 != team {
			continue
		}
		workers = append(workers, Tile{
			team,
			tile & 0x7, // height
			x,
			y,
		})
	}
	return workers
}

func (f *FastBoard) ValidTurns(team uint8) []*Turn {
	turns := make([]*Turn, 0, 8*2)
	workers := f.GetWorkers(team)
	// Get all the valid moves for each worker
	for _, w := range workers {
		moves := f.GetSurroundingTiles(w.x, w.y)
		for _, move := range moves {
			if err := w.CanMoveTo(move); err != nil {
				continue
			}
			// for all the moves, find the possible builds
			builds := f.GetSurroundingTiles(move.x, move.y)
			for _, build := range builds {
				// cant build on an occupied spot, UNLESS that spot is where the worker currently is
				if build.IsOccupied() && !build.SameLocation(w) {
					continue
				}
				turns = append(turns, &Turn{
					Worker: w,
					MoveTo: move,
					Build:  build,
				})
			}
		}
	}
	return turns
}

func (f *FastBoard) GetSurroundingTiles(x, y uint8) (tiles []Tile) {
	// List all surrounding tiles
	tiles = make([]Tile, 0, 8)
	if y > 0 {
		// North
		index := f.board[(f.width*(y-1))+x]
		tiles = append(tiles, Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x,
			y:      y - 1,
		})
	}
	if y < f.height-1 {
		// South
		index := f.board[(f.width*(y+1))+x]
		tiles = append(tiles, Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x,
			y:      y + 1,
		})
	}
	if x > 0 {
		// West
		index := f.board[(f.width*(y))+x-1]
		tiles = append(tiles, Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x - 1,
			y:      y,
		})
	}
	if x < f.width-1 {
		// East
		index := f.board[(f.width*(y))+x+1]
		tiles = append(tiles, Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x + 1,
			y:      y,
		})
	}
	if y > 0 && x < f.width-1 {
		// NorthEast
		index := f.board[(f.width*(y-1))+x+1]
		tiles = append(tiles, Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x + 1,
			y:      y - 1,
		})
	}
	if y > 0 && x > 0 {
		// NorthWest
		index := f.board[(f.width*(y-1))+x-1]
		tiles = append(tiles, Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x - 1,
			y:      y - 1,
		})
	}
	if y < f.height-1 && x < f.width-1 {
		// SouthEast
		index := f.board[(f.width*(y+1))+x+1]
		tiles = append(tiles, Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x + 1,
			y:      y + 1,
		})
	}
	if y < f.height-1 && x > 0 {
		// SouthEast
		index := f.board[(f.width*(y+1))+x-1]
		tiles = append(tiles, Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x - 1,
			y:      y + 1,
		})
	}
	return
}

func getDistance(t1, t2 Tile) float64 {
	dx := math.Pow(float64(t2.GetX())-float64(t1.GetX()), 2)
	dy := math.Pow(float64(t2.GetY())-float64(t1.GetY()), 2)
	return math.Sqrt(dx + dy)
}
