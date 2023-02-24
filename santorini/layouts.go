package santorini

import "fmt"

// BlankBoard sets no players, it is useful to override the default
func BlankBoard(b Board) Board {
	return b
}

func Default2Player(b Board) Board {
	b.setTile(1, 0, 2, 1)
	b.setTile(1, 0, 2, 3)
	b.setTile(2, 0, 1, 2)
	b.setTile(2, 0, 3, 2)
	return b
}

func CustomSize(w, h uint8) func(Board) Board {
	if w > 29 || h > 29 {
		panic(fmt.Errorf("Cannot make a board greater than 29x29: Got %dx%d", w, h))
	}
	return func(b Board) Board {
		return &FastBoard{
			width:  w,
			height: h,
			board:  make([]uint8, int(w)*int(h)),
		}
	}
}
