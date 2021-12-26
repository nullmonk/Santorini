package main

import (
	"bufio"
	"fmt"
	"os"
	santorini "santorini/pkg"
	"santorini/pkg/color"
	"strconv"
	"strings"
)

type TurnSelector interface {
	SelectTurn() santorini.Turn
}

func playTurn(board *santorini.Board, selector TurnSelector) {

}

func main() {
	// Initialize a new board
	board := santorini.NewBoard()

	// Initialize team workers
	workerA1 := &santorini.Worker{Team: 1, Number: 1}
	workerA2 := &santorini.Worker{Team: 1, Number: 2}
	workerB1 := &santorini.Worker{Team: 2, Number: 1}
	workerB2 := &santorini.Worker{Team: 2, Number: 2}

	// Place workers
	board.PlaceWorker(workerA1, 2, 1)
	board.PlaceWorker(workerA2, 2, 3)
	board.PlaceWorker(workerB1, 1, 2)
	board.PlaceWorker(workerB2, 3, 2)

	// Initialize RNG Team 1
	team1 := RandomSelector{
		Board:   board,
		Workers: []santorini.Worker{*workerA1, *workerA2},
	}

	// Initialize Team 2
	team2 := RandomSelector{
		Board:   board,
		Workers: []santorini.Worker{*workerB1, *workerB2},
	}

	// REPL
	reader := bufio.NewReader(os.Stdin)
	for round := 0; round < 1000; round++ {
		fmt.Printf("\nStarting Round %d\n\n", round+1)
		// Print the board
		fmt.Println(board)

		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if text == "exit" {
			return
		}
		if strings.Contains(text, ",") {
			parts := strings.Split(text, ",")
			x, err := strconv.ParseInt(parts[0], 10, 32)
			if err != nil {
				panic(err)
			}
			y, err := strconv.ParseUint(parts[1], 10, 8)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Information for Tile %d,%d:\n", x, y)
			tile := board.GetTile(uint8(x), uint8(y))
			if tile.Worker != nil {
				fmt.Printf("\tWorker %d: Team %d ", tile.Worker.Number, tile.Worker.Team)
			}
		}

		// Team 1 Move
		turn1 := team1.SelectTurn()
		if turn1 == nil {
			fmt.Printf("Team 2 Wins! Team 1 has no remaining moves")
			break
		}

		board.PlayTurn(*turn1)
		fmt.Printf("Team 1 moves %sWorker %d%s to %d,%d and builds %d,%d\n",
			color.GetWorkerColor(turn1.Worker.Team, turn1.Worker.Number),
			turn1.Worker.Number,
			color.Reset,
			turn1.MoveTo.GetX(),
			turn1.MoveTo.GetY(),
			turn1.Build.GetX(),
			turn1.Build.GetX(),
		)
		fmt.Printf("\n%s\n\n", board)

		// Team 2 Move
		turn2 := team2.SelectTurn()
		if turn2 == nil {
			fmt.Printf("Team 1 Wins! Team 2 has no remaining moves")
			break
		}
		board.PlayTurn(*turn2)
		fmt.Printf("Team 2 moves %sWorker %d%s to %d,%d and builds %d,%d\n",
			color.GetWorkerColor(turn2.Worker.Team, turn2.Worker.Number),
			turn2.Worker.Number,
			color.Reset,
			turn2.MoveTo.GetX(),
			turn2.MoveTo.GetY(),
			turn2.Build.GetX(),
			turn2.Build.GetX(),
		)
		fmt.Printf("\n%s\n\n", board)
	}
}
