package ui

import (
	"os"
	santorini "santorini/pkg"
	"strings"

	"github.com/gen64/go-tui"
	"github.com/sirupsen/logrus"
)

type Game struct {
	Board       *santorini.Board
	turnCounter int // the turn in the game it is Round = turnCounter/len(Teams)
	Teams       []santorini.TurnSelector

	widgets struct {
		Board     *BoardWidget  // Displays the board
		Prompt    *PromptWidget // Displays prompts for the user
		Teams     *TeamWidget   // lists the teams and workers
		Logs      *LogWidget
		inputPane *tui.TUIPane // Input that the user can type into
	}

	input     string // The unread input buffer
	lastInput string // Last input that was successfully entered

	lastPrompt string
	promptBuf  string
	t          *tui.TUI
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
	var boardPane, promptPane, teamPane, logPane *tui.TUIPane
	boardPane, logPane = g.t.GetPane().SplitVertically(-(5*g.Board.Size + 4), tui.UNIT_CHAR)

	// The logpane is the center of the screen containing the logs of the game and prompts for moves
	logPane, teamPane = logPane.SplitVertically((5*g.Board.Size + 4), tui.UNIT_CHAR)
	logPane, g.widgets.inputPane = logPane.SplitHorizontally(3, tui.UNIT_CHAR)

	// Prompt pane displays mini prompts to the user
	g.widgets.inputPane.SetMinHeight(1)
	g.widgets.inputPane.SetStyle(tui.NewTUIPaneStyleFrame())
	g.widgets.inputPane.SetOnDraw(func(p *tui.TUIPane) int {
		p.Write(1, 1, "Press Enter to Start Game", false)
		return 1
	})

	boardPane, promptPane = boardPane.SplitHorizontally(-(g.Board.Size*3 + 3), tui.UNIT_CHAR)

	/*g.widgets.logPane.SetStyle(&tui.TUIPaneStyle{
		NE: "─", NW: "─", SE: "─", SW: "─", E: " ", W: " ", N: "─", S: "─",
	})
	*/
	//

	g.widgets.Teams = NewTeamWidget(g.Teams, g.Board, teamPane)
	g.widgets.Board = NewBoardWidget(g.Board, boardPane)
	g.widgets.Prompt = NewPromptWidget(promptPane)
	g.widgets.Logs = NewLogWidget(logPane)
	logger = g.widgets.Logs
	g.widgets.Prompt.Set("Press ↵ to start game")

	g.t.SetOnKeyPress(func(t *tui.TUI, b []byte) {
		// Backspace
		if b[0] == 127 {
			if len(g.promptBuf) > 1 {
				g.promptBuf = g.promptBuf[:len(g.promptBuf)-1]
			} else {
				g.promptBuf = ""
			}
		} else if string(b) == "\n" {
			g.lastPrompt = g.promptBuf
			g.promptBuf = ""
			if strings.HasPrefix(g.lastPrompt, "print: ") {
				g.widgets.Logs.Printf(g.lastPrompt)
			} else {
				g.Step()
			}
		} else {
			g.promptBuf += string(b)
		}
		writeLine(1, 1, g.widgets.inputPane, g.promptBuf)
	})
	return g
}

// Perform the next step in the game
func (g *Game) Step() {
	// Figure out whose turn it is next
	if g.Board.IsOver {
		if g.lastPrompt == "exit" {
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
