package ui

import (
	"os"
	santorini "santorini/pkg"

	"github.com/gen64/go-tui"
	"github.com/sirupsen/logrus"
)

type Game struct {
	Board       *santorini.Board
	turnCounter int // the turn in the game it is Round = turnCounter/len(Teams)
	Teams       []santorini.TurnSelector

	widgets struct {
		Board  *BoardWidget  // Displays the board
		Prompt *PromptWidget // Displays prompts for the user
		Teams  *TeamWidget   // lists the teams and workers
		Logs   *LogWidget
		Input  *InputWidget // Input that the user can type into
	}

	t *tui.TUI
}

func NewGame(bots ...santorini.BotInitializer) *Game {
	// Make the panes
	g := &Game{
		t:     tui.NewTUI("", "", ""),
		Board: santorini.DefaultPosition(len(bots)),
		Teams: make([]santorini.TurnSelector, 0, len(bots)),
	}
	for i, bot := range bots {
		g.Teams = append(g.Teams, bot(i+1, g.Board, logrus.StandardLogger()))
	}
	var boardPane, promptPane, teamPane, logPane, inputPane *tui.TUIPane

	boardPane, logPane = g.t.GetPane().SplitVertically(-(5*g.Board.Size + 4), tui.UNIT_CHAR)
	logPane, teamPane = logPane.SplitVertically((5*g.Board.Size + 4), tui.UNIT_CHAR)
	logPane, inputPane = logPane.SplitHorizontally(3, tui.UNIT_CHAR)
	boardPane, promptPane = boardPane.SplitHorizontally(-(g.Board.Size*3 + 3), tui.UNIT_CHAR)

	g.widgets.Teams = NewTeamWidget(g.Teams, g.Board, teamPane)
	g.widgets.Board = NewBoardWidget(g.Board, boardPane)
	g.widgets.Prompt = NewPromptWidget(promptPane)
	g.widgets.Logs = NewLogWidget(logPane)
	g.widgets.Input = NewInputWidget(inputPane)
	logger = g.widgets.Logs
	g.widgets.Prompt.Set("Press ↵ to start game")

	g.t.SetOnKeyPress(func(t *tui.TUI, b []byte) {
		if g.widgets.Input.onKeyPress(t, b) {
			g.Step()
		}
	})
	return g
}

// Perform the next step in the game
func (g *Game) Step() {
	// Figure out whose turn it is next
	if g.Board.IsOver {
		if g.widgets.Input.Value() == "exit" {
			os.Exit(0)
		}
		return
	}

	botNum := g.turnCounter % len(g.Teams)
	g.turnCounter += 1
	bot := g.Teams[botNum]
	turn := bot.SelectTurn()
	g.widgets.Logs.LogTurn(bot, *turn)
	if g.Board.PlayTurn(*turn) {
		g.widgets.Logs.Printf("Game Over. %s wins in %d turns", bot.Name(), g.turnCounter/len(g.Teams))
		g.widgets.Prompt.Set("Type 'exit' to quit")
	} else {
		if g.turnCounter == 0 {
			g.widgets.Prompt.Set("Press ↵ to start game")
		} else {
			g.widgets.Prompt.Set("Press ↵ to continue")
		}
	}
	g.Refresh()
}

func (g *Game) Refresh() {
	g.widgets.Board.Iterate()
	g.widgets.Prompt.Iterate()
	g.widgets.Teams.Iterate()
}

func (g *Game) Run() {
	g.t.Run(os.Stdout, os.Stderr)
}
