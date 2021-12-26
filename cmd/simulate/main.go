package main

import (
	"fmt"
	"log"
	"santorini/bots"
	santorini "santorini/pkg"
)

const MaxSimulations = 100

type TurnSelector interface {
	SelectTurn() *santorini.Turn
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
		fmt.Printf("Team 2 Wins! Team 1 has no remaining moves\n")
		return true
	}

	// Team 1 Play
	if sim.Board.PlayTurn(*turn1) {
		fmt.Printf("Team 1 Wins!\n")
		return true
	}

	// Team 2 Select
	turn2 := sim.Team2.SelectTurn()
	if turn2 == nil {
		fmt.Printf("Team 1 Wins! Team 2 has no remaining moves\n")
		return true
	}

	// Team 2 Play
	if sim.Board.PlayTurn(*turn2) {
		fmt.Printf("Team 2 Wins!\n")
		return true
	}

	// The game continues...
	return false
}

// Run a game until it's completion
func (sim *Simulation) Run() {
	for !sim.doRound() {
		log.Printf("Completed Round %d", sim.round)
	}

	log.Printf("Simulation %d Completed, Team %d won after %d rounds", sim.Number, sim.Board.Victor, sim.round)
}

func main() {
	team1Wins := 0
	team2Wins := 0

	for i := 0; i < MaxSimulations; i++ {
		// Initialize a new board
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

		team1 := bots.NewBasicBot(1, board)
		team2 := bots.NewKyleBot(2, board)
		/*
			team2 := bots.NewBasicBot(2, board)
			team1 := bots.KyleBot{
				Team:      1,
				EnemyTeam: 2,
				Board:     board,
			}
		*/

		// Initialize Simulation
		sim := &Simulation{
			Number: i,
			Team1:  team1,
			Team2:  team2,
			Board:  board,
		}
		sim.Run()

		fmt.Printf("Final Board:\n%s\n", board)
		if board.Victor == 1 {
			team1Wins += 1
		} else {
			team2Wins += 1
		}
	}

	log.Printf("End of %d Simulations, Team 1 won %d and Team 2 won %d", MaxSimulations, team1Wins, team2Wins)
}
