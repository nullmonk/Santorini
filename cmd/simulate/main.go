package main

import (
	"encoding/json"
	"fmt"
	"log"
	"santorini/bots"
	santorini "santorini/pkg"
)

const MaxSimulations = 10000

type TurnSelector interface {
	SelectTurn() *santorini.Turn
	Name() string
}

type Simulation struct {
	Number int
	Board  *santorini.Board
	Team1  TurnSelector
	Team2  TurnSelector

	round int
}

// doRound returns true when a team wins, false otherwise
func (sim *Simulation) doRound() bool {
	sim.round += 1

	// Team 1 Select
	turn1 := sim.Team1.SelectTurn()
	if turn1 == nil {
		fmt.Printf("Team 2 (%s) Wins! Team 1 has no remaining moves\n", sim.Team2.Name())
		return true
	}

	// Team 1 Play
	if sim.Board.PlayTurn(*turn1) {
		fmt.Printf("Team 1 (%s) Wins!\n", sim.Team1.Name())
		return true
	}

	// Team 2 Select
	turn2 := sim.Team2.SelectTurn()
	if turn2 == nil {
		fmt.Printf("Team 1 (%s) Wins! Team 2 has no remaining moves\n", sim.Team1.Name())
		return true
	}

	// Team 2 Play
	if sim.Board.PlayTurn(*turn2) {
		fmt.Printf("Team 2 (%s) Wins!\n", sim.Team2.Name())
		return true
	}

	// The game continues...
	return false
}

// Run a game until it's completion
func (sim *Simulation) Run() {
	for !sim.doRound() {
		//log.Printf("Completed Round %d", sim.round)
	}

	log.Printf("Simulation %d Completed, Team %d won after %d rounds", sim.Number, sim.Board.Victor, sim.round)
}

func defaultPosition() *santorini.Board {
	board := santorini.NewBoard()

	// Select Worker Tiles
	workerTileA1 := board.GetTile(2, 1)
	workerTileA2 := board.GetTile(2, 3)
	workerTileB1 := board.GetTile(1, 2)
	workerTileB2 := board.GetTile(3, 2)

	// Place Workers
	board.PlaceWorker(1, 1, workerTileA1)
	board.PlaceWorker(1, 2, workerTileA2)
	board.PlaceWorker(2, 1, workerTileB1)
	board.PlaceWorker(2, 2, workerTileB2)
	return board
}

func main() {
	team1Wins := 0
	team2Wins := 0

	for i := 0; i < MaxSimulations; i++ {
		board := defaultPosition()
		// Initialize Simulation
		sim := &Simulation{
			Number: i,
			Team1:  bots.NewRandomBot(1, board),
			Team2:  bots.NewKyleBot(2, board),
			Board:  board,
		}
		sim.Run()

		board = defaultPosition()
		// Initialize Simulation
		sim = &Simulation{
			Number: i,
			Team1:  bots.NewRandomBot(1, board),
			Team2:  bots.NewKyleBot(2, board),
			Board:  board,
		}
		sim.Run()
		//fmt.Printf("Final Board:\n%s\n", board)
		if board.Victor == 1 {
			logData, _ := json.Marshal(board.Moves)
			fmt.Printf("LOSS DATA: %s\n", string(logData))
			team1Wins += 1
		} else {
			team2Wins += 1
		}
	}

	log.Printf("End of %d Simulations, Team 1 won %d and Team 2 won %d", MaxSimulations, team1Wins, team2Wins)
}
