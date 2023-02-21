package santorini

import (
	"encoding/json"
	"fmt"
	"os"
	"santorini/board"

	"github.com/sirupsen/logrus"
)

type BotInitializer func(team int, board board.Board, logger *logrus.Logger) TurnSelector

type TurnSelector interface {
	// The name of the bot
	Name() string

	// Perform the next turn for the bot
	SelectTurn(b board.Board) *board.Turn

	// True if the bot will perform the same given the same inputs
	IsDeterministic() bool
}

type Simulation struct {
	Number int
	Board  board.Board
	Teams  []TurnSelector
	logger *logrus.Logger
	round  int
}

func NewSimulator(number int, logger *logrus.Logger, bots ...BotInitializer) *Simulation {
	b := DefaultPosition(len(bots))
	lgr := logger
	// Unless we are debugging, hide all bot logs except for fatal ones
	if logger.Level != logrus.DebugLevel {
		// Copy the formatter
		lgr = &logrus.Logger{
			Out:       logger.Out,
			Hooks:     logger.Hooks,
			Formatter: logger.Formatter,
			Level:     logrus.FatalLevel,
		}
	}

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
		logger: logger,
	}
}

// doRound returns true when a team wins, false otherwise
func (sim *Simulation) doRound() bool {
	sim.round += 1
	// Loop vars here so they can be used by panic
	var bot TurnSelector
	var i int
	defer func() {
		if err := recover(); err != nil {
			b, _ := json.Marshal(sim.Board)
			fmt.Fprintln(os.Stderr, bot.Name(), "caused a panic:", err)
			fmt.Println(string(b))
			os.Exit(1)
		}
	}()
	for i, bot = range sim.Teams {
		turn := bot.SelectTurn()
		if turn == nil {
			sim.logger.Debugf("Team %d (%s) has no moves", i+1, bot.Name())
			continue
		}

		isOver, err := sim.Board.PlayTurn(turn)
		if isOver {
			return true
		}
		if err != nil {
			panic(err)
		}
	}

	return false
}

// Run a game until it's completion
func (sim *Simulation) Run() {
	for !sim.doRound() {
		//log.Printf("Completed Round %d", sim.round)
	}

	sim.logger.Debugf("Simulation %d Completed, Team %d (%s) won after %d rounds", sim.Number, sim.Board.Victor, sim.Teams[sim.Board.Victor-1].Name(), sim.round)
}

// Default starting position for bots
func DefaultPosition(numTeams int) *Board {
	board := NewBoard()

	workers := make([][]int, 0, numTeams)

	if numTeams == 3 {
		workers = append(workers,
			[]int{1, 1, 0, 1},
			[]int{1, 2, 4, 1},
			[]int{2, 1, 2, 0},
			[]int{2, 2, 2, 4},
			[]int{3, 1, 0, 3},
			[]int{3, 2, 4, 3})
	} else {
		workers = append(workers,
			[]int{1, 1, 2, 1},
			[]int{1, 2, 2, 3})
		if numTeams == 2 {
			workers = append(workers,
				[]int{2, 1, 1, 2},
				[]int{2, 2, 3, 2},
			)
		}
	}

	for _, w := range workers {
		board.PlaceWorker(w[0], w[1], w[2], w[3])
	}

	return board
}
