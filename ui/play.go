package ui

import (
	"fmt"
	santorini "santorini/pkg"

	"github.com/sirupsen/logrus"
)

func NewPlayerInitializer(g *Game) santorini.BotInitializer {
	return func(team int, board *santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
		return &Player{
			game: g,
			team: team,
			name: fmt.Sprintf("Team %d", team),
		}
	}
}

type Player struct {
	game *Game
	name string
	team int
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) IsDeterministic() bool {
	return false
}

func (p *Player) SelectTurn() *santorini.Turn {
	return nil
}
