package board

import (
	"fmt"
	"math"
)

// fastBoard implements the board interface and uses a single array as the storage mechanism
type FastBoard struct {
	// Last three bits are the tile height, rest is team number
	board     []uint8
	teams     uint8 // Number of teams
	workers   uint8 // num workers per team
	width     uint8
	height    uint8
	turnCount uint // number of turns that have been taken on this board (used to calculate which teams turn it is)
}

func NewFastBoard(options ...func(Board)) *FastBoard {
	f := &FastBoard{
		board:   make([]uint8, 25),
		width:   5,
		height:  5,
		workers: 2,
		teams:   2,
	}

	return f
}

func (f *FastBoard) Clone() Board {
	return &FastBoard{
		board:     f.board,
		width:     f.width,
		height:    f.height,
		workers:   f.workers,
		teams:     f.teams,
		turnCount: f.turnCount,
	}
}

func (f *FastBoard) PlayTurn(t Turn) (victory bool, err error) {
	whosTurn := f.turnCount%uint(f.teams) + 1
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
	if turn.Worker.team != uint8(whosTurn) || turn.Worker.team != t.Worker.team {
		return false, fmt.Errorf("the worker does not match the team taking the turn")
	}

	// TODO: Do we force workers to stop checks? i.e. dont let them make a move that would cause them to lose?
	// If so, we need to implement that logic here

	// Make sure the move is valid
	if err := turn.Worker.CanMoveTo(turn.MoveTo); err != nil {
		return false, err
	}

	// cant build on an occupied spot, UNLESS that spot is where the worker currently is
	if turn.Build.IsOccupied() && !turn.Build.Same(t.Worker) {
		return false, fmt.Errorf("the worker cannot build on the given block")
	}
	if turn.Build.Same(turn.MoveTo) {
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
	if team > f.teams {
		return fmt.Errorf("invalid team chosen")
	}
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

func (f *FastBoard) ValidTurns(team uint8) []Turn {
	turns := make([]Turn, 0, 8*f.workers)
	workers := make([]Tile, 0, f.workers)
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

	// Get all the valid moves for each worker
	for _, w := range workers {
		moves := f.getSurroundingTiles(w.x, w.y)
		for _, move := range moves {
			if err := w.CanMoveTo(move); err != nil {
				continue
			}
			// for all the moves, find the possible builds
			builds := f.getSurroundingTiles(move.x, move.y)
			for _, build := range builds {
				// cant build on an occupied spot, UNLESS that spot is where the worker currently is
				if build.IsOccupied() && !build.Same(w) {
					continue
				}
				turns = append(turns, Turn{
					Worker: w,
					MoveTo: move,
					Build:  build,
				})
			}
		}
	}
	return turns
}

func (f *FastBoard) getSurroundingTiles(x, y uint8) (tiles []Tile) {
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
