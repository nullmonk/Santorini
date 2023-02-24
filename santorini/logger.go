package santorini

import (
	"fmt"
	"io"
	"strings"
)

type GameLog struct {
	b     Board            // The game we are logging
	teams []string         // idx = teamNum - 1, str = team name
	fil   *strings.Builder // file for the entire log
	buf   []string         // buffer to write game logs too
}

func NewGameLog(b Board, team ...string) *GameLog {
	g := &GameLog{
		b:     b,
		teams: team,
		fil:   new(strings.Builder),
		buf:   make([]string, 0, 2),
	}
	g.fil.WriteString("BOARD: ")
	g.fil.WriteString(b.Hash())
	g.fil.WriteRune('\n')
	g.Comment("TEAMS", strings.Join(team, " vs. "))
	g.writeComments()
	return g
}

func (g *GameLog) RegisterTeam(name string) {
	g.teams = append(g.teams, name)
}

func (g *GameLog) Comment(who, format string, a ...interface{}) {
	if who != "" {
		who += ": "
	}
	g.buf = append(g.buf, who+fmt.Sprintf(format, a...))
}

func (g *GameLog) LogTurn(t *Turn, who ...string) {
	g.fil.WriteString("TURN: ")
	g.fil.WriteString(t.String())
	if len(who) > 0 {
		g.fil.WriteString(" # " + who[0])
	}
	g.fil.WriteRune('\n')
	g.writeComments()
}

func (g *GameLog) writeComments() {
	for _, l := range g.buf {
		g.fil.WriteString("# ")
		g.fil.WriteString(l)
		g.fil.WriteRune('\n')
	}
	g.buf = g.buf[:0]
}

func (g *GameLog) Dump(f io.StringWriter) {
	g.writeComments()
	f.WriteString(g.fil.String())
}
