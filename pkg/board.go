package santorini

import "fmt"

type Tile struct {
	Height int // 0 no building, 4 capped
	Worker *Worker

	x uint8
	y uint8
}

func (t Tile) IsOccupied() bool {
	return t.Worker != nil
}
func (t Tile) IsOccupiedBy(worker Worker) bool {
	if t.Worker == nil {
		return false
	}

	return t.Worker.Team == worker.Team && t.Worker.Number == worker.Number
}

func (t Tile) IsCapped() bool {
	return t.Height > 3
}

type Board struct {
	Size  uint8
	Tiles [][]*Tile

	IsOver bool
	Victor int // Who won the game
	Moves  []Turn
}

// NewBoard initializes a game with the default board size and two teams
func NewBoard(options ...func(*Board)) *Board {
	// Default Board
	/*
	 *  0 0 0 0 0
	 *  0 0 X 0 0
	 *  0 Y 0 Y 0
	 *  0 0 X 0 0
	 *  0 0 0 0 0
	 */
	board := &Board{
		Size: 5,
	}

	// Apply Options
	for _, opt := range options {
		opt(board)
	}

	// Build Tiles
	board.Tiles = make([][]*Tile, board.Size)
	for x := 0; x < int(board.Size); x++ {
		board.Tiles[x] = make([]*Tile, board.Size)

		for y := 0; y < int(board.Size); y++ {
			board.Tiles[x][y] = &Tile{
				x: uint8(x),
				y: uint8(y),
			}
		}
	}

	return board
}

func (board Board) GetTile(x, y uint8) (t *Tile) {
	if x >= BoardSize {
		panic(fmt.Errorf("invalid x"))
	}
	if y >= BoardSize {
		panic(fmt.Errorf("invalid y"))
	}
	return board.Tiles[x][y]
}

func (board Board) GetSurroundingTiles(x, y uint8) (tiles []Tile) {
	// List all surrounding tiles
	type Position struct {
		X uint8
		Y uint8
	}
	candidates := []Position{
		{x, y + 1},     // North
		{x, y - 1},     // South
		{x + 1, y},     // East
		{x - 1, y},     // West
		{x + 1, y + 1}, // Northeast
		{x - 1, y + 1}, // Northwest
		{x - 1, y + 1}, // Southeast
		{x - 1, y - 1}, // Southwest
	}

	// Filter potential tiles
	for _, candidate := range candidates {
		if candidate.X >= board.Size {
			continue
		}
		if candidate.Y >= board.Size {
			continue
		}

		// Otherwise, it is a valid tile
		tiles = append(tiles, *board.GetTile(x, y))
	}

	return
}

// GetMoveableTiles returns all tiles that may be moved to from the provided position.
func (board Board) GetMoveableTiles(x, y uint8) (tiles []Tile) {
	candidates := board.GetSurroundingTiles(x, y)
	// Filter invalid tiles
	for _, candidate := range candidates {
		// Occupied Constraints
		if candidate.IsOccupied() {
			continue
		}

		// Capped Constraints
		if candidate.IsCapped() {
			continue
		}

		// Height Constraints
		curTile := board.GetTile(x, y)
		if candidate.Height > curTile.Height+1 {
			continue
		}

		// Otherwise, it is a valid move
		tiles = append(tiles, candidate)
	}

	return
}

// GetBuildableTiles returns all tiles that may be built from the provided position.
func (board Board) GetBuildableTiles(x, y uint8, worker Worker) (tiles []Tile) {
	candidates := board.GetSurroundingTiles(x, y)

	// Filter invalid tiles
	for _, candidate := range candidates {
		// Occupied Constraints
		if candidate.IsOccupied() && !candidate.IsOccupiedBy(worker) {
			continue
		}

		// Capped Constraints
		if candidate.IsCapped() {
			continue
		}

		// Otherwise, it is a valid move
		tiles = append(tiles, candidate)
	}

	return
}

func (board *Board) PlayTurn(turn Turn) (gameover bool) {
	board.Moves = append(board.Moves, turn)

	// Step 1. Move
	srcTile := board.Tiles[turn.Worker.X][turn.Worker.Y]
	worker := srcTile.Worker

	dstTile := board.Tiles[turn.MoveTo.x][turn.MoveTo.y]

	// TODO Validate the move here

	// Update the worker location
	worker.X = turn.MoveTo.x
	worker.Y = turn.MoveTo.y

	srcTile.Worker = nil
	dstTile.Worker = worker

	// Check if the game has been won
	if dstTile.Height == 3 {
		board.Victor = worker.Team
		board.IsOver = true
		return true
	}

	// Step 2. Build
	dstTile = board.Tiles[turn.Build.x][turn.Build.y]

	// Validate the build here
	dstTile.Height += 1
	return false
}

func (board *Board) PlaceWorker(worker *Worker, x, y uint8) {
	board.Tiles[x][y].Worker = worker
}
