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
	Workers      []santorini.Tile
	EnemyWorkers []santorini.Tile
	Team         int

	randomBot     RandomSelector
	logger        log.Logger
	turns         []santorini.Turn // Turns for the round, by worker
	turnsByWorker map[int][]santorini.Turn
}

func NewBasicBot(team int, board *santorini.Board) *BasicBot {
	// Figure out where my workers are, and figure out where the enemy workers are
	ai := &BasicBot{
		Board:        board,
		Workers:      make([]santorini.Tile, 0, 2),
		EnemyWorkers: make([]santorini.Tile, 0, 2),
		Team:         team,

		logger:        *log.Default(),
		turnsByWorker: make(map[int][]santorini.Turn),
	}
	return ai
}

// Update the board status
func (bb *BasicBot) update() {
	bb.turns = bb.Board.GetValidTurns(bb.Team)

	bb.Workers = make([]santorini.Tile, 0, 2)
	// Figure out where my workers are, and where the enemy workers are
	for _, tile := range bb.Board.Tiles {
		if tile.IsOccupied() {
			if tile.GetWorker() == bb.Team {
				bb.Workers = append(bb.Workers, tile)
			} else {
				bb.EnemyWorkers = append(bb.EnemyWorkers, tile)
			}
		}
	}

	// sort the turns by the worker
	for i, _ := range bb.Workers {
		bb.turnsByWorker[i] = make([]santorini.Turn, 0, 10)
	}

	for _, turn := range bb.turns {
		bb.turnsByWorker[turn.Worker] = append(bb.turnsByWorker[turn.Worker], turn)
	}
}

func (bb *BasicBot) SelectTurn() *santorini.Turn {
	bb.update()
	if winningMoves := GetWinningMoves(bb.turns); len(winningMoves) > 0 {
		bb.logger.Print("Detected a winning move. Executing it")
		return &winningMoves[0]
	}

	// If we need to defend, do it
	if t := bb.defend(); t != nil {
		return t
	}

	// if a worker is almost trapped, get them out
	if t := bb.escapeTraps(); t != nil {
		return t
	}
	// Get the status of the workers
	stats := bb.getWorkerStatus()
	diff := stats[0] - stats[1]
	var workerToMove int
	if diff > 5 {
		// W1 is exceeding w2 by too much, balance by moving w2
		workerToMove = 1
	} else if diff < -5 {
		// w2 is exceeding, improve w1
		workerToMove = 2
	} else if diff > 0 {
		// w1 is doing good, work with it
		workerToMove = 2
	} else if diff < 0 {
		workerToMove = 1
	} else {
		workerToMove = 1
	}

	// If we cant move the worker that we need to, use the other
	if bb.turnsByWorker[workerToMove] == nil {
		// chose a move for the worker
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(bb.turns)-1)))
		if err != nil {
			panic(err)
		}
		return &bb.turns[n.Int64()]
	}

	bb.logger.Printf("Worker Stats: %v. Chose worker %v. (%v turns)", stats, workerToMove, len(bb.turnsByWorker[workerToMove]))
	// chose a move for the worker
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(bb.turnsByWorker[workerToMove])-1)))
	if err != nil {
		panic(err)
	}
	return &bb.turnsByWorker[workerToMove][n.Int64()]
}

// Check if a winning move exists in any of the possible moves
func GetWinningMoves(turns []santorini.Turn) []santorini.Turn {
	res := make([]santorini.Turn, 0, 1)
	for _, t := range turns {
		// Is winning turn
		if t.MoveTo.GetHeight() == 3 {
			res = append(res, t)
		}
	}
	return res
}

// defend tries to stop the enemy
func (bb *BasicBot) defend() *santorini.Turn {
	// See if the enemy can win, if they can, then try to block them
	defendMoves := make([]santorini.Turn, 0, 10) // Moves that we can make to defend ourselves

	enemyWinningMoves := GetWinningMoves(bb.Board.GetValidTurns(bb.EnemyWorkers[0].GetTeam()))

	// Try to block the enemy winning moves
	for _, et := range enemyWinningMoves {
		for _, myturn := range bb.turns {
			//if I can build where the enemy will go, then do it
			if myturn.Build.GetX() == et.MoveTo.GetX() && myturn.Build.GetY() == et.MoveTo.GetY() {
				defendMoves = append(defendMoves, myturn)
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
		return &defendMoves[0]
	}

	// Other defense goes here?
	return nil
}

// Worker Status is a metric of how good of a position a worker is in,
// We want to keep the workers close to the same stat (e.g. dont let one fall behind)
// But still allow the other to excel
func (bb *BasicBot) getWorkerStatus() (statuses []int) {
	statuses = make([]int, len(bb.Workers))

	// The bot with more moves is in a better position
	if len(bb.turnsByWorker[0]) > len(bb.turnsByWorker[1]) {
		statuses[0] += 1
	} else if len(bb.turnsByWorker[0]) < len(bb.turnsByWorker[1]) {
		statuses[1] += 1
	}

	for i, workerTile := range bb.Workers {
		// height is good in general
		statuses[i] += workerTile.GetHeight()

		// edges are a negative for defensibility, we want to be able float around and defend
		if workerTile.GetX() == 0 || workerTile.GetX() == bb.Board.Size-1 || workerTile.GetY() == 0 || workerTile.GetY() == bb.Board.Size {
			// We are on an edge
		} else {
			statuses[i] += 1
		}

		// surrounding blocks of the >= height are good, get a bonus for that
		for _, tile := range bb.Board.GetSurroundingTiles(workerTile.GetX(), workerTile.GetY()) {
			if tile.GetHeight() > 0 && tile.GetHeight() >= workerTile.GetHeight() {
				statuses[i] += 1
			}
		}
	}
	return statuses
}

// If a worker is close to being trapped, have it escape
func (bb *BasicBot) escapeTraps() *santorini.Turn {
	for _, tile := range bb.Workers {
		if len(bb.Board.GetMoveableTiles(tile)) == 1 {
			bb.logger.Printf("Worker %d is trapped, escaping", tile.GetWorker())
			return &bb.turnsByWorker[tile.GetWorker()][0]
		} else if len(bb.Board.GetMoveableTiles(tile)) == 0 {
			bb.logger.Printf("Worker %d is trapped!! %v %v", tile.GetWorker(), bb.Board.GetMoveableTiles(tile), len(bb.turnsByWorker[tile.GetWorker()]))
		}
	}
	return nil
}

// Perform a build for the worker
func (bb *BasicBot) build(worker int) *santorini.Turn {
	// First of all see if we have any turns for this worker
	return nil
}
