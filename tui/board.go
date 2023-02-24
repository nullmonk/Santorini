package ui

import (
	"fmt"
	"santorini/santorini"

	"github.com/gen64/go-tui"
)

// BoardWidget displays the board and updates it as needed
type BoardWidget struct {
	board santorini.Board
	pane  *tui.TUIPane
}

func NewBoardWidget(board santorini.Board, pane *tui.TUIPane) *BoardWidget {
	b := &BoardWidget{
		pane:  pane,
		board: board,
	}
	pane.SetStyle(tui.NewTUIPaneStyleFrame())
	pane.SetMinHeight(1)
	pane.SetOnDraw(b.update)
	pane.SetOnIterate(b.update)
	return b
}

func (b *BoardWidget) update(p *tui.TUIPane) int {
	w, h := b.board.Dimensions()
	for y := uint8(0); y < h; y++ {
		for x := uint8(0); x < w; x++ {
			tile := b.board.GetTile(x, y)
			tileIcon := fmt.Sprint(tile.GetHeight())
			if tile.GetHeight() == 4 {
				tileIcon = "^"
			}
			tileIcon = fmt.Sprintf("%s%v%s", GetWorkerColor(int(tile.GetTeam()), 1), tileIcon, Reset)
			p.Write(5*int(x)+2, 3*int(y)+1, tileIcon+tileIcon+tileIcon, false)
			p.Write(5*int(x)+2, 3*int(y)+2, tileIcon+tileIcon+tileIcon, false)
		}
	}
	return 1
}

func (b *BoardWidget) Iterate() {
	b.pane.Iterate()
}
