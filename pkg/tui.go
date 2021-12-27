package santorini

import (
	"fmt"
	"math"
	"os"
	"strings"

	"santorini/pkg/color"

	"github.com/gen64/go-tui"
	"github.com/sirupsen/logrus"
)

func (board Board) String() string {
	rows := make([]string, board.Size)
	for y := 0; y < board.Size; y++ {
		columns := make([]string, board.Size)
		for x := 0; x < board.Size; x++ {
			tile := board.GetTile(x, y)
			display := fmt.Sprintf("%s%d%s", color.GetWorkerColor(tile.team, tile.worker), tile.height, color.Reset)
			columns = append(columns, display)
		}
		rows[y] = strings.Join(columns, " ")
	}

	return strings.Join(rows, "\n")
}

func (board Board) Draw(pane tui.TUIPane) string {
	rows := make([]string, board.Size)
	for y := 0; y < board.Size; y++ {
		columns := make([]string, board.Size)
		for x := 0; x < board.Size; x++ {
			tile := board.GetTile(x, y)
			display := fmt.Sprintf("%s%d%s", color.GetWorkerColor(tile.team, tile.worker), tile.height, color.Reset)
			columns = append(columns, display)
		}
		rows[y] = strings.Join(columns, " ")
	}

	return strings.Join(rows, "\n")
}

type Game struct {
	Board       *Board
	turnCounter int // the turn in the game it is Round = turnCounter/len(Teams)
	Logs        *logBuffer
	Teams       []TurnSelector

	boardPane *tui.TUIPane
	turnPane  *tui.TUIPane
	logPane   *tui.TUIPane
	teamPane  *tui.TUIPane

	lastPrompt string
	promptBuf  string
	cursorX    int
	promptPane *tui.TUIPane
	t          *tui.TUI
}

func NewGame(bots ...BotInitializer) *Game {
	// Make the panes
	g := &Game{
		t:     tui.NewTUI("", "", ""),
		Board: defaultPosition(len(bots)),
		Teams: make([]TurnSelector, 0, len(bots)),
	}
	for i, bot := range bots {
		g.Teams = append(g.Teams, bot(i+1, g.Board, logrus.StandardLogger()))
	}
	g.boardPane, g.logPane = g.t.GetPane().SplitVertically(-(5*g.Board.Size + 4), tui.UNIT_CHAR)
	g.logPane, g.teamPane = g.logPane.SplitVertically((5*g.Board.Size + 4), tui.UNIT_CHAR)
	g.logPane, g.promptPane = g.logPane.SplitHorizontally(3, tui.UNIT_CHAR)
	g.promptPane.SetMinHeight(1)
	g.promptPane.SetOnDraw(func(p *tui.TUIPane) int {
		p.Write(1, 1, "Press Enter to Start Game", false)
		return 1
	})
	g.teamPane.SetStyle(tui.NewTUIPaneStyleFrame())
	g.teamPane.SetOnDraw(g.updateTeams)
	g.teamPane.SetOnIterate(g.updateTeams)
	g.boardPane, g.turnPane = g.boardPane.SplitHorizontally(-(g.Board.Size*3 + 3), tui.UNIT_CHAR)
	g.boardPane.SetStyle(tui.NewTUIPaneStyleFrame())
	g.boardPane.SetOnDraw(g.updateBoard)

	/*g.logPane.SetStyle(&tui.TUIPaneStyle{
		NE: "─", NW: "─", SE: "─", SW: "─", E: " ", W: " ", N: "─", S: "─",
	})
	*/
	g.logPane.SetStyle(tui.NewTUIPaneStyleFrame())
	//
	g.turnPane.SetStyle(tui.NewTUIPaneStyleFrame())
	g.turnPane.SetMinHeight(1)
	g.turnPane.SetOnDraw(g.updateTurnStatus)
	g.turnPane.SetOnIterate(g.updateTurnStatus)

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
				g.Logs.Printf(g.lastPrompt)
			} else {
				g.Step()
			}
		} else {
			g.promptBuf += string(b)
		}
		writeLine(1, 1, g.promptPane, g.promptBuf)
	})
	g.Logs = newLogBuffer(g.logPane)
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
	if turn.IsVictory() {
		g.Logs.Printf("%s moves %sWorker %d%s to %d,%d",
			bot.Name(),
			color.GetWorkerColor(turn.Team, turn.Worker),
			turn.Worker,
			color.Reset,
			turn.MoveTo.GetX(),
			turn.MoveTo.GetY(),
		)
	} else {
		g.Logs.Printf("%s moves %sWorker %d%s to %d,%d and builds %d,%d",
			bot.Name(),
			color.GetWorkerColor(turn.Team, turn.Worker),
			turn.Worker,
			color.Reset,
			turn.MoveTo.GetX(),
			turn.MoveTo.GetY(),
			turn.Build.GetX(),
			turn.Build.GetY(),
		)
	}
	if g.Board.PlayTurn(*turn) {
		g.Logs.Printf("Game Over. %s wins in %d turns", bot.Name(), g.turnCounter/len(g.Teams))
		g.Logs.Printf("Type 'exit' to exit the game")
	}
	g.boardPane.Draw()
	g.Prompt("Press Enter to Continue")
}
func (g *Game) Run() {
	g.t.Run(os.Stdout, os.Stderr)
}

func (g *Game) updateTurnStatus(p *tui.TUIPane) int {
	msg := "Press ↵ to continue"
	if g.turnCounter == 0 {
		msg = "Press ↵ to start game"
	} else if g.Board.IsOver {
		msg = "Type 'exit' to quit"
	}
	writeLine(1, 0, p, msg)
	return 0
}

func (g *Game) Prompt(prompt string) string {
	writeLine(1, 1, g.promptPane, fmt.Sprintf("%s: ", prompt))
	return ""
}

func (g *Game) updateTeams(p *tui.TUIPane) int {
	y := 0
	for i, team := range g.Teams {
		teamName := fmt.Sprintf("Team %d. %s", i+1, team.Name())
		if g.turnCounter%len(g.Teams) == i {
			teamName += " *"
		}
		writeLine(1, y, p, teamName)
		y++
		for j, worker := range g.Board.GetWorkerTiles(i + 1) {
			p.Write(3, y, fmt.Sprintf("%sWorker %d%s (%d, %d)", color.GetWorkerColor(i+1, j+1), j+1, color.Reset, worker.GetX(), worker.GetY()), false)
			y++
		}
		y++
	}
	return 1
}

func (g *Game) updateBoard(p *tui.TUIPane) int {
	// Print the actual board
	for y := 0; y < g.Board.Size; y++ {
		for x := 0; x < g.Board.Size; x++ {
			tile := g.Board.GetTile(x, y)
			tileIcon := fmt.Sprint(tile.height)
			if tile.height == 4 {
				tileIcon = "^"
			}
			tileIcon = fmt.Sprintf("%s%v%s", color.GetWorkerColor(tile.team, tile.worker), tileIcon, color.Reset)
			p.Write(5*x+2, 3*y+1, tileIcon+tileIcon+tileIcon, false)
			p.Write(5*x+2, 3*y+2, tileIcon+tileIcon+tileIcon, false)
		}
	}
	return 1
}

type logBuffer struct {
	items []string
	pane  *tui.TUIPane
}

func newLogBuffer(pane *tui.TUIPane) *logBuffer {
	lb := &logBuffer{
		items: make([]string, pane.GetHeight()*4),
		pane:  pane,
	}
	pane.SetOnDraw(lb.up)
	pane.SetOnIterate(lb.up)
	return lb
}
func (l *logBuffer) Printf(format string, a ...interface{}) {
	ln := strings.TrimSpace(fmt.Sprintf(format, a...))

	// Word wrap the line so it fits in the gamebox
	lns := wordWrap(ln, l.pane.GetWidth()-2, "↵")
	if len(ln) > 0 {
		l.items = append(l.items, lns...)
		l.pane.Iterate()
	}
}

// Update the logs
func (l *logBuffer) up(p *tui.TUIPane) int {
	if l.items == nil {
		return 0
	}
	min := int(math.Max(float64(len(l.items)-p.GetHeight()+2), 0))
	for i := 0; i < p.GetHeight()-2; i++ {
		var logLn string
		if min+i < len(l.items) {
			logLn = l.items[min+i]
		}
		writeLine(0, i, p, logLn)
	}
	/*
		for i, line := range l.items[min:] {
			p.Write(0, i, line, false)
		}*/
	return 1
}

func writeLine(x, y int, p *tui.TUIPane, line string) {
	style := p.GetStyle()
	borderWidth := 0
	if style != nil && len(strings.TrimSpace(style.E+style.W)) > 0 {
		borderWidth = 1
	}
	maxWidth := p.GetWidth() - x - borderWidth - borderWidth - x
	if len(line) > maxWidth {
		p.Write(x, y, line[len(line)-maxWidth:], false)
	} else {
		p.Write(x, y, line+strings.Repeat(" ", maxWidth-len(line)), false)
	}
}

func wordWrap(text string, lineWidth int, prefix ...string) []string {
	rows := make([]string, 0, 2)
	pref := ""
	if len(prefix) > 0 {
		pref = prefix[0]
		lineWidth -= len(pref)
	}
	size := len(text)
	for size > lineWidth {
		rows = append(rows, strings.TrimSpace(text[:lineWidth])+pref)
		text = strings.TrimSpace(text[lineWidth:])
		size = len(text)
	}
	rows = append(rows, text)
	return rows
}
