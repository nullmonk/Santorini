package ui

import (
	"fmt"
	"math"
	"santorini/santorini"
	"strings"

	"github.com/gen64/go-tui"
)

var logger *LogWidget

type LogWidget struct {
	logs []string
	pane *tui.TUIPane
}

func NewLogWidget(pane *tui.TUIPane) *LogWidget {
	lb := &LogWidget{
		logs: make([]string, pane.GetHeight()*4),
		pane: pane,
	}
	pane.SetStyle(tui.NewTUIPaneStyleFrame())
	pane.SetOnDraw(lb.update)
	pane.SetOnIterate(lb.update)
	return lb
}
func (l *LogWidget) Printf(format string, a ...interface{}) {
	ln := strings.TrimSpace(fmt.Sprintf(format, a...))
	// Word wrap the line so it fits in the gamebox
	lns := wordWrap(ln, l.pane.GetWidth()-2, "↵")
	if len(ln) > 0 {
		l.logs = append(l.logs, lns...)
		l.pane.Iterate()
	}
}

// Log a turn that was taken by a bot
func (l *LogWidget) LogTurn(bot santorini.TurnSelector, turn santorini.Turn) {
	msg := fmt.Sprintf("%s moves %sWorker %d%s to %d,%d",
		bot.Name(),
		GetWorkerColor(int(turn.Worker.GetTeam()), 1),
		turn.Worker,
		Reset,
		turn.MoveTo.GetX(),
		turn.MoveTo.GetY(),
	)
	if !turn.IsWinningMove() {
		msg += fmt.Sprintf(" and builds %d,%d", turn.Build.GetX(), turn.Build.GetY())
	}
	l.Printf(msg)
}

// Update the logs
func (l *LogWidget) update(p *tui.TUIPane) int {
	if l.logs == nil {
		return 0
	}
	min := int(math.Max(float64(len(l.logs)-p.GetHeight()+2), 0))
	for i := 0; i < p.GetHeight()-2; i++ {
		var logLn string
		if min+i < len(l.logs) {
			logLn = l.logs[min+i]
		}
		writeLine(0, i, p, logLn)
	}
	return 1
}
