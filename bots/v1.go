package bots

import (
	"santorini/santorini"
)

/*
func v1RankAll() *santorini.Turn {
	// Auto win/defend moves are handled by Bot already, we can focus on the good good
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
	return &bb.turns[len(bb.turns)-1] // use the last move (Highest ranked)
}
*/

func NewV1Bot() santorini.BotInitializer {
	return NewBasicBot("V1Bot", rankMove, false)
}

func rankMove(bot *Bot, b santorini.Board, turn *santorini.Turn) int {
	rank := 0
	width, height := b.Dimensions()

	if turn.Worker.GetHeight() < turn.MoveTo.GetHeight() {
		// Moving up is generally good
		rank += 100
	} else {
		height := turn.Worker.GetHeight() - turn.MoveTo.GetHeight()
		// if we have to jump down, atleast dont do it near an edge
		if turn.MoveTo.GetX() == 0 || turn.MoveTo.GetX() == width || turn.MoveTo.GetY() == 0 || turn.MoveTo.GetY() == height {
			if turn.MoveTo.GetHeight() == 0 {
				rank -= 15
			}
		}
		if height == 2 {
			rank -= 100
		} else if height == 1 {
			rank -= 50
		}
	}

	if (turn.Build.GetX() == width || turn.Build.GetX() == 0) && (turn.Build.GetY() == height || turn.Build.GetY() == 0) {
		// Corner strat og
		rank += 2
	}

	// Dont build 2 up (unless capping, which is already handled)
	if turn.MoveTo.HeightDiff(turn.Build) > 0 {
		rank -= 50
	}

	if turn.Build.GetHeight() == 3 {
		rank -= 200 // never cap
	}

	surroundingBuild := b.GetSurroundingTiles(turn.Build.GetX(), turn.Build.GetY())
	for _, tile := range surroundingBuild {
		// Building next to the enemy is bad
		tt := tile.GetTeam()
		if tt > 0 && tt != bot.Team() {
			if turn.Build.GetHeight() == 2 {
				if tile.GetHeight() > 0 {
					// Building where the enemy can win is super bad
					rank -= 111111111
				}
				// Dont build things the enemy can just cap
				rank -= 50
			}

			if tile.HeightDiff(turn.Build) >= 1 {
				rank += 30
			}
			rank -= 30
		}
		// + rank for each building next to us (prevent holes from being built)

		if tile.GetHeight() > 0 {
			rank += 5 * int(tile.GetHeight())
		}
	}

	for _, tile := range b.GetSurroundingTiles(turn.MoveTo.GetX(), turn.MoveTo.GetY()) {
		// Try not to move directly next to my buddy (its good to be spread out)
		if tile.GetTeam() > 0 && tile.GetTeam() == bot.Team() {
			rank -= 30
		}
		// If there is a level 3 tile next to our moveto, and we can build a level 3 tile, we win
		if tile.GetHeight() == 3 {
			if turn.MoveTo.GetHeight() == 2 && turn.Build.GetHeight() == 2 {
				rank += 1000
			}
			rank += 75
		}
	}

	// Keep building up next to ourselves
	if turn.Build.GetHeight() > 0 && turn.MoveTo.GetHeight() > 0 {
		rank += 10 * int(turn.Build.GetHeight())
	}
	return rank
}

// If a worker is close to being trapped, have it escape
/*
func escapeTraps() *santorini.Turn {
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
*/
