package santorini

// BlankBoard sets no players, it is useful to override the default
func BlankBoard(b Board) {}

func Default2Player(b Board) {
	b.setTile(1, 0, 2, 1)
	b.setTile(1, 0, 2, 3)
	b.setTile(2, 0, 1, 2)
	b.setTile(2, 0, 3, 2)
}
