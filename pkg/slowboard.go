package santorini

import (
	"fmt"
)

type DefaultBoard struct {
	Size  int
	Tiles []Tile
	Teams map[int]bool // true if the player is playing (e.g. not trapped)

	IsOver bool
	Victor int // Who won the game
	Moves  []Turn

	lastTeam int
}

// NewDefaultBoard initializes a game with the default board size and two teams
func NewDefaultBoard(options ...func(*DefaultBoard)) *DefaultBoard {
	// Default DefaultBoard
	/*
	 *  0 0 0 0 0
	 *  0 0 X 0 0
	 *  0 Y 0 Y 0
	 *  0 0 X 0 0
	 *  0 0 0 0 0
	 */
	board := &DefaultBoard{
		Size:  5,
		Teams: make(map[int]bool),
	}

	// Apply Options
	for _, opt := range options {
		opt(board)
	}

	// Build Tiles
	board.Tiles = make([]Tile, board.Size*board.Size)
	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			index := (board.Size * y) + x
			board.Tiles[index] = Tile{
				x: x,
				y: y,
			}
		}
	}

	return board
}

func (board DefaultBoard) GetTiles() (tiles []Tile) {
	tiles = make([]Tile, len(board.Tiles))
	copy(tiles, board.Tiles)
	return
}

func (board DefaultBoard) GetTile(x, y int) (t Tile) {
	if x >= board.Size {
		panic(fmt.Errorf("invalid x"))
	}
	if y >= board.Size {
		panic(fmt.Errorf("invalid y"))
	}
	index := (board.Size * y) + x
	return board.Tiles[index]
}

func (board *DefaultBoard) setTile(tile Tile) {
	if tile.x >= board.Size || tile.x < 0 {
		panic(fmt.Errorf("invalid x %d", tile.x))
	}
	if tile.y >= board.Size || tile.y < 0 {
		panic(fmt.Errorf("invalid y %d", tile.y))
	}

	index := (board.Size * tile.y) + tile.x
	board.Tiles[index] = tile
}

func (board DefaultBoard) GetSurroundingTiles(x, y int) (tiles []Tile) {
	// List all surrounding tiles
	type Position struct {
		X int
		Y int
	}
	candidates := []Position{
		{x, y + 1},     // North
		{x, y - 1},     // South
		{x + 1, y},     // East
		{x - 1, y},     // West
		{x + 1, y + 1}, // Northeast
		{x - 1, y + 1}, // Northwest
		{x + 1, y - 1}, // Southeast
		{x - 1, y - 1}, // Southwest
	}

	// Filter potential tiles
	for _, candidate := range candidates {
		if candidate.X >= board.Size || candidate.X < 0 {
			continue
		}
		if candidate.Y >= board.Size || candidate.Y < 0 {
			continue
		}

		// Otherwise, it is a valid tile
		tiles = append(tiles, board.GetTile(candidate.X, candidate.Y))
	}

	return
}

// GetMoveableTiles returns all tiles that may be moved to from the provided position.
func (board DefaultBoard) GetMoveableTiles(curTile Tile) (tiles []Tile) {
	candidates := board.GetSurroundingTiles(curTile.x, curTile.y)
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
		if candidate.height > curTile.height+1 {
			continue
		}

		// Otherwise, it is a valid move
		tiles = append(tiles, board.GetTile(candidate.x, candidate.y))
	}

	return
}

// GetBuildableTiles returns all tiles that may be built from the provided position.
func (board DefaultBoard) GetBuildableTiles(team, worker int, buildTile Tile) (tiles []Tile) {
	candidates := board.GetSurroundingTiles(buildTile.x, buildTile.y)

	// Filter invalid tiles
	for _, candidate := range candidates {
		// Occupied Constraints
		if candidate.IsOccupied() && !candidate.IsOccupiedBy(team, worker) {
			continue
		}

		// Capped Constraints
		if candidate.IsCapped() {
			continue
		}

		// Otherwise, it is a valid build
		tiles = append(tiles, board.GetTile(candidate.x, candidate.y))
	}

	return
}

// PlayTurn will update the board state with the results of the provided turn, or panic if the turn is illegal
func (board *DefaultBoard) PlayTurn(turn Turn) (gameover bool) {
	// Have workers been trapped
	teamsInPlay := 0
	playingTeam := 0
	for team, playing := range board.Teams {
		if playing {
			teamsInPlay++
			playingTeam = team
		}
	}
	if teamsInPlay < 2 && len(board.Teams) > 1 {
		board.IsOver = true
		board.Victor = playingTeam
		return true
	}
	if turn.Team == board.lastTeam {
		panic(fmt.Errorf("it is not team %d's turn", turn.Team))
	}
	board.lastTeam = turn.Team

	if turn.Team == 0 {
		panic(fmt.Errorf("must set team taking the turn: %+v", turn))
	}
	if turn.Worker == 0 {
		panic(fmt.Errorf("must set worker used for the turn: %+v", turn))
	}

	board.Moves = append(board.Moves, turn)

	// 1. Clear existing tile
	workerTile := board.GetWorkerTile(turn.Team, turn.Worker)
	workerTile.team = 0
	workerTile.worker = 0
	board.setTile(workerTile)

	// 2. Update destination tile
	dstTile := board.GetTile(turn.MoveTo.x, turn.MoveTo.y)
	dstTile.team = turn.Team
	dstTile.worker = turn.Worker
	board.setTile(dstTile)

	// 3. Check if the game has been won
	// Has someone capped?
	if dstTile.height == 3 {
		board.Victor = turn.Team
		board.IsOver = true
		return true
	}

	// 4. Build
	buildTile := board.GetTile(turn.Build.x, turn.Build.y)
	if buildTile.height > 3 {
		panic(fmt.Errorf("cannot build tile %+v", turn))
	}
	buildTile.height += 1
	board.setTile(buildTile)

	// The Game Continues...
	return false
}

// PlaceWorker on the board, should be called before any turns are made
func (board *DefaultBoard) PlaceWorker(team, worker, x, y int) {
	workerTile := board.GetTile(x, y)
	workerTile.team = team
	workerTile.worker = worker
	board.setTile(workerTile)
	board.Teams[team] = true
}
