package santorini

type Tile uint8

type Board struct {
	tiles [][]Tile
}

func (b *Board) NewGame() {
	// BOARD_SIZE * BOARD_SIZE
	b.tiles = make([][]Tile, 5) // 5x5 board
	for i := range b.tiles {
		b.tiles[i] = make([]Tile, 5)
	}
}

func (b *Board) GetTile(x, y int) (t Tile, err error) {
	return
}
