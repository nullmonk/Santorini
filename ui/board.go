package ui

import (
	"fmt"
	santorini "santorini/pkg"
	"santorini/pkg/color"

	"github.com/gen64/go-tui"
)

// BoardWidget displays the board and updates it as needed
type BoardWidget struct {
	board *santorini.Board
	pane  *tui.TUIPane
}

func NewBoardWidget(board *santorini.Board, pane *tui.TUIPane) *BoardWidget {
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
	for y := 0; y < b.board.Size; y++ {
		for x := 0; x < b.board.Size; x++ {
			tile := b.board.GetTile(x, y)
			tileIcon := fmt.Sprint(tile.GetHeight())
			if tile.GetHeight() == 4 {
				tileIcon = "^"
			}
			tileIcon = fmt.Sprintf("%s%v%s", color.GetWorkerColor(tile.GetTeam(), tile.GetWorker()), tileIcon, color.Reset)
			p.Write(5*x+2, 3*y+1, tileIcon+tileIcon+tileIcon, false)
			p.Write(5*x+2, 3*y+2, tileIcon+tileIcon+tileIcon, false)
		}
	}
	return 1
}

func (b *BoardWidget) Iterate() {
	b.pane.Iterate()
}
