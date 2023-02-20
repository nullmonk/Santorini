package bots

import (
	"fmt"
	santorini "santorini/pkg"
	"sort"

	"github.com/sirupsen/logrus"
)

/* BasicBot is a bot that will perform the following actions:
 *
 * 1. If the bot can win, do it
 * 2. If the enemy can win, and we can block it, then do that
 * 3. Random
 */
type BasicBot struct {
	Board        *santorini.Board
	Workers      map[int]santorini.Tile
	EnemyWorkers []santorini.Tile
	Team         int

	logger *logrus.Logger
	turns  []santorini.Turn // Turns for the round, by worker

	chosenWorker  int // the worker we recommend moving
	turnsByWorker map[int][]santorini.Turn
}

func (bb *BasicBot) Name() string {
	return "BasicBot"
}

func (bb BasicBot) IsDeterministic() bool {
	return true
}

func (bb *BasicBot) log(fmtstr string, args ...interface{}) {
	bb.logger.Debug("BasicBot: ", fmt.Sprintf(fmtstr, args...))
}

func NewBasicBot(team int, board *santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
	// Figure out where my workers are, and figure out where the enemy workers are
	ai := &BasicBot{
		Board:        board,
		Workers:      make(map[int]santorini.Tile, 2),
		EnemyWorkers: make([]santorini.Tile, 0, 2),
		Team:         team,

		logger:        logger,
		turnsByWorker: make(map[int][]santorini.Turn),
	}
	return ai
}

// Update the board status
func (bb *BasicBot) update() {
	bb.turns = bb.Board.GetValidTurns(bb.Team)

	bb.Workers = make(map[int]santorini.Tile, 2)
	bb.EnemyWorkers = make([]santorini.Tile, 0, 2)
	// Figure out where my workers are, and where the enemy workers are
	for _, tile := range bb.Board.Tiles {
		if tile.IsOccupied() {
			if tile.GetTeam() == bb.Team {
				bb.Workers[tile.GetWorker()] = tile
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
		bb.log("Detected a winning move. Executing it")
		return &winningMoves[0]
	}

	// If we need to defend, do it
	if t := bb.defend(); t != nil {
		return t
	}

	// if a worker is almost trapped, get them out
	if t := bb.escapeTraps(); t != nil {
		bb.chosenWorker = t.Worker
		bb.rankMove(*t)
	}

	bb.sortMoves()
	/* Debug, print top ten moves and weights
	for x, i := range bb.turns {
		if x > len(bb.turns)-10 {
			fmt.Printf("%d %+v\n", bb.rankMove(i), i)
		}
	}
	*/
	return &bb.turns[len(bb.turns)-1] // use the last move (Highest ranked)
}

func (bb *BasicBot) rankMove(turn santorini.Turn) int {
	rank := 0

	worker := bb.Workers[turn.Worker]
	// if the worker is moving up/down, add/remove points (going up good)
	if diff := turn.MoveTo.GetHeight() - worker.GetHeight(); diff > 0 {
		rank += 50 // move up
	} else if diff == -2 {
		rank -= 100
	} else if diff == -1 {
		rank -= 20
	}

	// dislike corners and edges
	rank -= (8 - len(bb.Board.GetSurroundingTiles(turn.Build.GetX(), turn.Build.GetY()))) * 5
	// dont like moving to corner
	if len(bb.Board.GetSurroundingTiles(turn.MoveTo.GetX(), turn.MoveTo.GetY())) == 3 {
		rank -= 20
	}

	// if the move will limit us in the future, subtract a point
	if len(bb.Board.GetMoveableTiles(turn.MoveTo)) < 2 {
		rank -= 10
	}
	if len(bb.Board.GetBuildableTiles(bb.Team, -1, turn.MoveTo)) < 2 {
		rank -= 10
	}

	// Dont build 2 up (unless capping, which is already handled)
	if turn.Build.GetHeight() > turn.MoveTo.GetHeight() {
		rank -= 30
	} else if turn.Build.GetHeight()+1 == 3 {
		// If the build is increasing the height to 3, super rank it
		rank += 30
	} else if turn.Build.GetHeight()+1 > turn.MoveTo.GetHeight() {
		// Building up next to ourselves is good (as oposed to starting on the ground)
		rank += 20
	} else if turn.Build.GetHeight() > 0 {
		rank += 30
	}

	surroundingBuild := bb.Board.GetSurroundingTiles(turn.Build.GetX(), turn.Build.GetY())

	for _, tile := range surroundingBuild {
		if tile.IsOccupied() && tile.GetTeam() != bb.Team {
			if turn.Build.GetHeight() == 2 {
				rank -= 111111111
			}
			rank -= 10
		}
		if tile.GetHeight() > 0 {
			rank += 3
		}
	}

	// Build on blocks that are touching other blocks laterally
	for _, tile := range bb.Board.GetSurroundingTiles(turn.MoveTo.GetX(), turn.MoveTo.GetY()) {
		// Try not to move next to my buddy
		if tile.GetTeam() == bb.Team {
			rank -= 30
		}
	}
	if turn.Build.GetHeight() == 2 && turn.MoveTo.GetHeight() == 2 {
		rank += 10
	}

	// use the recommended worker
	if turn.Worker == bb.chosenWorker {
		rank += 100000
	}
	return rank
}

func (bb *BasicBot) sortMoves() {
	sort.Slice(bb.turns, func(i, j int) bool {
		return bb.rankMove(bb.turns[i]) < bb.rankMove(bb.turns[j])
	})
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
		bb.log("Enemy has more winning moves than I can block")
	}
	// if we need to defend ourselves, do it
	// TODO: order the defend moves based on how good the move is
	if len(defendMoves) > 0 {
		bb.log("Capping enemy for defense")
		sort.Slice(defendMoves, func(i, j int) bool {
			return bb.rankMove(defendMoves[i]) < bb.rankMove(defendMoves[j])
		})
		return &defendMoves[len(defendMoves)-1] // use the last move (Highest ranked)
	}

	// Other defense goes here?
	return nil
}

// If a worker is close to being trapped, have it escape
func (bb *BasicBot) escapeTraps() *santorini.Turn {
	for _, tile := range bb.Workers {
		if len(bb.Board.GetMoveableTiles(tile)) == 1 {
			bb.log("Worker %d is trapped, escaping", tile.GetWorker())
			if len(bb.turnsByWorker[tile.GetWorker()]) > 0 {
				return &bb.turnsByWorker[tile.GetWorker()][0]
			}
		} else if len(bb.Board.GetMoveableTiles(tile)) == 0 {
			bb.log("Worker %d is trapped!! %v %v", tile.GetWorker(), bb.Board.GetMoveableTiles(tile), len(bb.turnsByWorker[tile.GetWorker()]))
		}
	}
	return nil
}
