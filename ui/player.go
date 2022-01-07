package ui

import (
	"fmt"
	santorini "santorini/pkg"
	"santorini/pkg/color"
	"strconv"

	"github.com/gen64/go-tui"
	"github.com/sirupsen/logrus"
)

func NewPlayerInitializer(g *Game) santorini.BotInitializer {
	return func(team int, board *santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
		return &Player{
			game: g,
			team: team,
		}
	}
}

type Player struct {
	game *Game
	name string
	team int

	awaitAnswers []interface{}
	hijacked     bool

	turnStage    int // 0 - select worker, 1 select move, 2 select turn, 3 turn complete
	selectedTurn santorini.Turn
}

func (p *Player) SetName(name string) {
	p.name = name
}

func (p *Player) Name() string {
	return fmt.Sprintf("Player %d", p.team)
}

func (p *Player) IsDeterministic() bool {
	return false
}

func (p *Player) SetChoices(prompt string, options map[string]interface{}) {
	if len(options) == 1 {
		return
	}

	p.game.widgets.Prompt.Set(prompt)
	p.game.widgets.Logs.Printf(prompt)
	p.awaitAnswers = make([]interface{}, 0, len(options))
	i := 0
	for name, choice := range options {
		p.game.widgets.Logs.Printf("\t%d. %s\n", i, name)
		i++
		p.awaitAnswers = append(p.awaitAnswers, choice)
	}
}

func (p *Player) GetChoice() interface{} {
	selection := p.game.widgets.Input.Value()
	if i, err := strconv.ParseInt(selection, 10, 64); err == nil {
		if int(i) < len(p.awaitAnswers) {
			answer := p.awaitAnswers[int(i)]
			p.awaitAnswers = nil
			return answer
		}
	}
	return nil
}

// Take over the game control from the main game. Used to prompt the player for moves
func (p *Player) Hijack() {
	if p.hijacked {
		return
	}
	// backup the input function so we can pass control back to the game
	p.game.t.SetOnKeyPress(func(t *tui.TUI, b []byte) {
		if p.game.widgets.Input.onKeyPress(t, b) {
			p.resume() // Continue the turn selection for the player
		}
	})
}

// Has the player selected a turn yet?
func (p *Player) isFinished() bool {
	return p.turnStage > 2
}

func (p *Player) getMove(chosen interface{}) {
	if chosen == nil {
		options := make(map[string]interface{})
		worker := p.game.Board.GetWorkerTile(p.team, p.selectedTurn.Worker)
		for _, tile := range p.game.Board.GetMoveableTiles(worker) {
			heightDiff := tile.GetHeight() - worker.GetHeight()
			label := ""
			if heightDiff > 0 {
				label = "up 1 tile"
			} else if heightDiff == -1 {
				label = "down 1 tile"
			} else if heightDiff < -1 {
				label = "down 2 tiles"
			}

			if tile.GetHeight() == 3 {
				label = "Winning Move!"
			}
			name := fmt.Sprintf("%s (%d,%d)", GetTileDir(worker, tile), tile.GetX(), tile.GetY())
			if label != "" {
				name += " - " + label
			}
			options[name] = p.game.Board.GetTile(tile.GetX(), tile.GetY())
		}
		p.SetChoices("Choose a move", options)
	} else {
		p.selectedTurn.MoveTo = chosen.(santorini.Tile)
		p.turnStage++
	}
}

func (p *Player) getWorker(chosen interface{}) {
	// If we arent passed a tile, then we need to ask for it
	if chosen == nil {
		workers := p.game.Board.GetWorkerTiles(p.team)
		options := make(map[string]interface{})

		for _, worker := range workers {
			movesAvailable := p.game.Board.GetMoveableTiles(worker)
			if len(movesAvailable) > 0 {
				name := fmt.Sprintf("%sWorker %d (%d, %d)%s",
					color.GetWorkerColor(p.team, worker.GetWorker()),
					worker.GetWorker(),
					worker.GetX(),
					worker.GetY(),
					color.Reset)
				options[name] = worker.GetWorker()
			}
		}

		p.SetChoices("Choose worker", options)
	} else {
		// We are passed a worker, set the value
		p.selectedTurn.Worker = chosen.(int)
		p.turnStage++
	}
}

func (p *Player) getBuild(chosen interface{}) {
	if chosen == nil {
		options := make(map[string]interface{})
		for _, tile := range p.game.Board.GetBuildableTiles(p.team, p.selectedTurn.Worker, p.selectedTurn.MoveTo) {
			name := fmt.Sprintf("%s (%d,%d)", GetTileDir(p.selectedTurn.MoveTo, tile), tile.GetX(), tile.GetY())
			label := ""
			if tile.GetHeight() == 3 {
				label = "Cap tile"
			}
			if label != "" {
				name += " - " + label
			}
			options[name] = tile
		}

		p.SetChoices("Choose a tile to build on", options)
	} else {
		p.selectedTurn.Team = p.team
		p.selectedTurn.Build = chosen.(santorini.Tile)
		p.turnStage++
	}
}
func (p *Player) resume() {
	var chosen interface{} // whatever option was selected by the player
	// See if we have input that we are awaiting
	if p.awaitAnswers != nil {
		if chosen = p.GetChoice(); chosen == nil {
			return
		}
		p.awaitAnswers = nil
	}
	// if we dont have a worker chosen, select a worker
	if p.turnStage == 0 {
		p.getWorker(chosen)
		chosen = nil
	}

	if p.turnStage == 1 {
		p.getMove(chosen)
		chosen = nil
		if p.selectedTurn.MoveTo.GetHeight() == 3 {
			p.turnStage++
			p.game.Step()
			return
		}
	}

	// If we won right here, dont even prompt for a build
	if p.turnStage == 2 {
		p.getBuild(chosen)
		p.game.Step()
	}

	p.game.Refresh()
}

func (p *Player) SelectTurn() *santorini.Turn {
	// backup the input function so we can pass control back to the game object
	p.game.t.SetOnKeyPress(func(t *tui.TUI, b []byte) {
		if p.game.widgets.Input.onKeyPress(t, b) {
			p.game.Step() // Continue the turn selection for the player
		}
	})
	p.turnStage = 0
	p.awaitAnswers = nil
	p.hijacked = false
	return &p.selectedTurn
}

func GetTileDir(src, dst santorini.Tile) string {
	dx := dst.GetX() - src.GetX()
	dy := dst.GetY() - src.GetY()
	if dx < 0 && dy == 0 {
		return "←" // left
	}
	if dx == 0 && dy < 0 {
		return "↑" // up
	}
	if dx > 0 && dy == 0 {
		return "→" // right
	}
	if dx == 0 && dy > 0 {
		return "↓" // down
	}
	if dx < 0 && dy < 0 {
		return "↖" // up, left
	}
	if dx > 0 && dy < 0 {
		return "↗" // up, right
	}
	if dx > 0 && dy > 0 {
		return "↘" // down, right
	}
	if dx < 0 && dy > 0 {
		return "↙" // down, left
	}
	return ""
}
