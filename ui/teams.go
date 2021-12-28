package ui

import (
	"fmt"
	santorini "santorini/pkg"
	"santorini/pkg/color"

	"github.com/gen64/go-tui"
)

type TeamWidget struct {
	pane  *tui.TUIPane
	board *santorini.Board
	bots  []santorini.TurnSelector
	teams []*team
}

func NewTeamWidget(bots []santorini.TurnSelector, board *santorini.Board, pane *tui.TUIPane) *TeamWidget {
	t := &TeamWidget{
		pane:  pane,
		board: board,
		teams: make([]*team, len(bots)),
		bots:  bots,
	}
	pane.SetStyle(tui.NewTUIPaneStyleFrame())
	pane.SetMinHeight(1)
	pane.SetOnDraw(t.update)
	pane.SetOnIterate(t.update)

	// Initialize the team names
	for i, bot := range bots {
		t.teams[i] = &team{
			name:    fmt.Sprintf("Team %d. %s", i+1, bot.Name()),
			workers: make([]santorini.Tile, 2),
		}
	}
	return t
}

type team struct {
	name    string
	workers []santorini.Tile
}

func (t *TeamWidget) update(p *tui.TUIPane) int {

	// Update the locations of all the workers
	tiles := t.board.GetTiles()
	for _, tile := range tiles {
		if !tile.IsOccupied() {
			continue
		}
		team := t.teams[tile.GetTeam()-1]
		team.workers[tile.GetWorker()-1] = tile
	}

	y := 0
	for i, team := range t.teams {
		isturn := ""
		if len(t.board.Moves)%len(t.teams) == i {
			isturn = " *"
		}
		writeLine(1, y, p, team.name+isturn)
		y++
		for j, worker := range team.workers {
			p.Write(3, y, fmt.Sprintf("%sWorker %d%s (%d, %d)", color.GetWorkerColor(i+1, j+1), j+1, color.Reset, worker.GetX(), worker.GetY()), false)
			y++
		}
		y++
	}
	return 1
}

func (t *TeamWidget) Iterate() {
	t.pane.Iterate()
}
