package santorini

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Simulation struct {
	Number int
	Board  Board
	Teams  []TurnSelector
	Logger *GameLog
	Victor uint8
	Turns  []string
	round  int
}

func NewSimulator(number int, logger *logrus.Logger, bots ...BotInitializer) *Simulation {
	b := NewBoard(Default2Player)
	lgr := NewGameLog(b, os.Args[1:len(bots)+1]...)
	var teams []TurnSelector

	switch len(bots) {
	case 3:
		teams = []TurnSelector{
			bots[0](1, b, lgr),
			bots[1](2, b, lgr),
			bots[2](3, b, lgr),
		}
	case 2:
		teams = []TurnSelector{
			bots[0](1, b, lgr),
			bots[1](2, b, lgr),
		}
	case 1:
		teams = []TurnSelector{
			bots[0](1, b, lgr),
		}
	}
	return &Simulation{
		Number: number,
		Board:  b,
		Teams:  teams,
		Logger: lgr,
		Turns:  make([]string, 0, 50),
	}
}

// doRound returns true when a team wins, false otherwise
func (sim *Simulation) doRound() (victor uint8) {
	sim.round += 1
	// Loop vars here so they can be used by panic
	var bot TurnSelector
	var i int
	defer func() {
		if err := recover(); err != nil {
			b, _ := json.Marshal(sim.Board)
			fmt.Fprintln(os.Stderr, bot.Name(), "caused a panic:", err)
			DumpBoard(sim.Board, os.Stdout)
			fmt.Println(string(b))
			os.Exit(1)
		}
	}()
	for i, bot = range sim.Teams {
		turn := bot.SelectTurn(sim.Board.Clone())
		if turn == nil {
			sim.Logger.Comment("Sim", "Team %d (%s) has no moves", i+1, bot.Name())
			return uint8((i+1)%len(sim.Teams)) + 1
		}

		isOver, err := sim.Board.PlayTurn(turn)
		sim.Logger.LogTurn(turn, bot.Name())
		sim.Turns = append(sim.Turns, turn.String())
		if isOver {
			return uint8(i) + 1
		}
		if err != nil {
			panic(err)
		}
	}

	return 0
}

// Run a game until it's completion
func (sim *Simulation) Run() {
	var victor uint8
	for victor == 0 {
		victor = sim.doRound()
	}
	sim.Victor = victor
	for i, bot := range sim.Teams {
		if int(victor)-1 == i {
			bot.GameOver(true)
		} else {
			bot.GameOver(false)
		}
	}
	sim.Logger.Comment("", "Simulation %d Completed, Team %d (%s) won after %d rounds", sim.Number, victor, sim.Teams[victor-1].Name(), sim.round)
}

func DumpBoardMini(b Board) string {
	width, height := b.Dimensions()
	line := ""
	for y := uint8(0); y < height; y++ {
		for x := uint8(0); x < width; x++ {
			tile := b.GetTile(x, y)
			tileIcon := TileIcon(tile)
			line += tileIcon
		}
		line += "\n"
	}
	return strings.TrimSpace(line)
}

func DumpBoard(b Board, p io.Writer) int {
	width, height := b.Dimensions()
	n := uint8(2)
	lines := make([]string, height*n)
	for y := uint8(0); y < height; y++ {
		lines[n*y] = ""
		lines[n*y+1] = ""
		for x := uint8(0); x < width; x++ {
			tile := b.GetTile(x, y)
			tileIcon := TileIcon(tile)
			lines[n*y] += tileIcon + tileIcon + tileIcon + "  "
			lines[n*y+1] += tileIcon + tileIcon + tileIcon + "  "
		}
	}

	for _, line := range lines {
		fmt.Fprintln(p, strings.TrimSpace(line))
	}
	return 1
}
