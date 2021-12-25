package santorini

import "fmt"

const BOARD_SIZE = 5

type Tile struct {
	Height int // 0 no building, 4 capped
	Worker interface{}
}

func (t Tile) IsOccupied() bool {
	if t.Height > 3 {
		return true
	}
	return t.Worker != nil
}

type Board struct {
	tiles [][]Tile
}

func NewBoard() *Board {
	b := new(Board)
	b.tiles = make([][]Tile, BOARD_SIZE)
	for i := range b.tiles {
		b.tiles[i] = make([]Tile, BOARD_SIZE)
	}
	return b
}

func (b Board) GetTile(x, y int) (t Tile) {
	if x >= BOARD_SIZE || x < 0 {
		panic(fmt.Errorf("invalid x"))
	}
	if y >= BOARD_SIZE || y < 0 {
		panic(fmt.Errorf("invalid x"))
	}
	return b.tiles[x][y]
}

func (b *Board) SetTile(x, y int, tile Tile) (err error) {
	return nil
}
