package santorini

import "fmt"

type Tile struct {
	Height int // 0 no building, 4 capped
	Worker *Worker
}

func (t Tile) IsOccupied() bool {
	return t.Worker != nil
}

func (t Tile) IsCapped() bool {
	return t.Height > 3
}

type Board struct {
	Size  uint8
	Tiles [][]Tile
	Teams []Team

	Turn uint // Which teams turn it is
}

// NewBoard initializes a game with the default board size and two teams
func NewBoard() *Board {
	return NewBoardCustom(BoardSize, 2)
}

// NewBoardCustom initializes a game with the given board size and number of teams
func NewBoardCustom(size uint8, teams uint8) *Board {
	b := new(Board)
	b.Size = size
	b.Tiles = make([][]Tile, b.Size)

	for i := range b.Tiles {
		b.Tiles[i] = make([]Tile, b.Size)
	}
	return b
}

func (b Board) GetTile(x, y uint8) (t Tile) {
	if x >= BoardSize {
		panic(fmt.Errorf("invalid x"))
	}
	if y >= BoardSize {
		panic(fmt.Errorf("invalid x"))
	}
	return b.Tiles[x][y]
}

func (b Board) MoveWorker(w Worker, x, y uint8) {
	// TODO Check that the move is valid
	// TODO Log the move
	currentTile := b.GetTile(w.X, w.Y)
	currentTile.Worker = nil
	newTile := b.GetTile(x, y)
	newTile.Worker = &w
}

func (b Board) Build(w Worker, x, y uint8) {
	// TODO log the build
	// TODO validate the build
	tile := b.GetTile(x, y)
	tile.Height += 1
}
