package santorini

import "fmt"

// BlankFastBoard sets no players, it is useful to override the default
func BlankBoard(b *FastBoard) *FastBoard {
	return b
}

func Default2Player(b *FastBoard) *FastBoard {
	b.setTile(1, 0, 2, 1)
	b.setTile(1, 0, 2, 3)
	b.setTile(2, 0, 1, 2)
	b.setTile(2, 0, 3, 2)
	return b
}

func CustomSize(w, h uint8) func(*FastBoard) *FastBoard {
	if w > 29 || h > 29 {
		panic(fmt.Errorf("Cannot make a board greater than 29x29: Got %dx%d", w, h))
	}
	return func(b *FastBoard) *FastBoard {
		return &FastBoard{
			width:  w,
			height: h,
			board:  make([]uint8, int(w)*int(h)),
		}
	}
}
