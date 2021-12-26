package bots

import (
	"crypto/rand"
	"log"
	"math/big"
	santorini "santorini/pkg"
)

/* BasicBot is a bot that will perform the following actions:
 *
 * 1. If the bot can win, do it
 * 2. If the enemy can win, and we can block it, then do that
 * 3. Random
 */
type BasicBot struct {
	Board        *santorini.Board
	Workers      []*santorini.Worker
	EnemyWorkers []*santorini.Worker

	randomBot RandomSelector
	logger    log.Logger
	turns     [][]santorini.Turn // Turns for the round
}

func NewSantoriniAI(team int, board *santorini.Board) *BasicBot {
	// Figure out where my workers are, and figure out where the enemy workers are
	ai := &BasicBot{
		Board:        board,
		Workers:      make([]*santorini.Worker, 0, 2),
		EnemyWorkers: make([]*santorini.Worker, 0, 2),

		randomBot: RandomSelector{
			Board:   board,
			Workers: make([]santorini.Worker, 0, 2),
		},
		turns: make([][]santorini.Turn, 0, 2),
	}

	for _, row := range board.Tiles {
		for _, tile := range row {
			if tile.Worker != nil {
				if tile.Worker.Team == team {
					ai.Workers = append(ai.Workers, tile.Worker)
					ai.randomBot.Workers = append(ai.randomBot.Workers, *tile.Worker)
				} else {
					ai.EnemyWorkers = append(ai.EnemyWorkers, tile.Worker)
				}
			}
		}
	}

	return ai
}

func (bb *BasicBot) SelectTurn() santorini.Turn {
	for i, worker := range bb.Workers {
		bb.turns[i] = worker.GetValidTurns(*bb.Board)
		// if we can win, do it
		if winningMoves := GetWinningMoves(bb.turns[i]); len(winningMoves) > 0 {
			bb.logger.Print("Detected a winning move. Executing it")
			return winningMoves[0]
		}
	}

	// If we need to defend, do it
	if t, ok := bb.defend(); ok {
		return t
	}

	// Get the status of the workers
	stats := bb.getWorkerStatus()
	diff := stats[0] - stats[1]
	var workerToMove int
	if diff > 2 {
		// W1 is exceeding w2 by too much, balance by moving w2
		workerToMove = 1
	} else if diff < -2 {
		// w2 is exceeding, improve w1
		workerToMove = 0
	} else if diff > 0 {
		// w1 is doing good, work with it
		workerToMove = 0
	} else if diff < 0 {
		workerToMove = 1
	}

	// If we cant move the worker that we need to, use the other
	if bb.turns[workerToMove] == nil {
		if workerToMove == 0 {
			workerToMove = 1
		} else {
			workerToMove = 0
		}
	}

	// chose a move for the worker
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(bb.turns[workerToMove])-1)))
	if err != nil {
		panic(err)
	}
	return bb.turns[workerToMove][n.Int64()]
}

// Check if a winning move exists in any of the possible moves
func GetWinningMoves(turns []santorini.Turn) []santorini.Turn {
	res := make([]santorini.Turn, 0, 1)
	for _, t := range turns {
		// Is winning turn
		if t.MoveTo.Height == 3 {
			res = append(res, t)
		}
	}
	return res
}

// defend tries to stop the enemy
func (bb *BasicBot) defend() (turn santorini.Turn, ok bool) {
	// See if the enemy can win, if they can, then try to block them
	enemyWinningMoves := make([]santorini.Turn, 0, 10)
	defendMoves := make([]santorini.Turn, 0, 10) // Moves that we can make to defend ourselves

	for _, enemy := range bb.EnemyWorkers {
		enemyturns := enemy.GetValidTurns(*bb.Board)
		enemyWinningMoves = append(enemyWinningMoves, GetWinningMoves(enemyturns)...)
	}

	// Try to block the enemy winning moves
	for _, et := range enemyWinningMoves {
		for _, workerturns := range bb.turns {
			for _, turn := range workerturns {
				//if I can build where the enemy will go, then do it
				if turn.Build.GetX() == et.MoveTo.GetX() && turn.Build.GetY() == et.MoveTo.GetY() {
					defendMoves = append(defendMoves, turn)
				}
			}
		}
	}

	if len(enemyWinningMoves) > len(defendMoves) {
		bb.logger.Println("Enemy has more winning moves than I cannot block")
	}
	// if we need to defend ourselves, do it
	// TODO: order the defend moves based on how good the move is
	if len(defendMoves) > 0 {
		bb.logger.Print("Capping enemy for defense")
		return defendMoves[0], true
	}

	// Other defense goes here?
	return
}

// Worker Status is a metric of how good of a position a worker is in,
// We want to keep the workers close to the same stat (e.g. dont let one fall behind)
// But still allow the other to excel
func (bb *BasicBot) getWorkerStatus() (statuses []int) {
	statuses = make([]int, len(bb.Workers))

	// The bot with more moves is in a better position
	if len(bb.turns[0]) > len(bb.turns[1]) {
		statuses[0] += 1
	} else if len(bb.turns[0]) < len(bb.turns[1]) {
		statuses[1] += 1
	}

	for i, w := range bb.Workers {
		workerTile := bb.Board.GetTile(w.X, w.Y)
		// height is good in general
		statuses[i] += workerTile.Height

		// edges are a negative for defensibility, we want to be able float around and defend
		if w.X == 0 || w.X == bb.Board.Size-1 || w.Y == 0 || w.Y == bb.Board.Size {
			// We are on an edge
		} else {
			statuses[i] += 1
		}

		// surrounding blocks of the >= height are good, get a bonus for that
		for _, tile := range bb.Board.GetSurroundingTiles(w.X, w.Y) {
			if tile.Height >= workerTile.Height {
				statuses[i] += 1
			}
		}
	}
	return statuses
}
