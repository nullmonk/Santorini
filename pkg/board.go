package santorini

import "fmt"

const BOARD_SIZE = 5

type Tile uint8

type Board struct {
	tiles [][]Tile
}

func NewBoard() *Board {
	// BOARD_SIZE * BOARD_SIZE
	b := new(Board)
	b.tiles = make([][]Tile, BOARD_SIZE) // 5x5 board
	for i := range b.tiles {
		b.tiles[i] = make([]Tile, BOARD_SIZE)
	}
	return b
}

func (b Board) GetTile(x, y int) (t Tile, err error) {
	if x > BOARD_SIZE {
		return 0, fmt.Errorf("x cannot be > %v", BOARD_SIZE)
	}
	return
}
