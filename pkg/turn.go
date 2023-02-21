package santorini

import "fmt"

// Turn stores a desired state change for the team
type Turn struct {
	Team   int
	Worker int
	MoveTo Tile
	Build  Tile
}

// IsVictory returns true if the turn would result in a victory
func (t Turn) IsVictory() bool {
	return t.MoveTo.height+1 > 3
}

// GetWorkerTile locates a particular worker's tile
func (board *DefaultBoard) GetWorkerTile(team, worker int) Tile {
	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			tile := board.GetTile(x, y)
			if tile.team == team && tile.worker == worker {
				return tile
			}
		}
	}
	panic(fmt.Errorf("failed to locate team %d worker %d", team, worker))
}

// GetWorkerTile locates all tiles that the provided team has workers on
func (board *DefaultBoard) GetWorkerTiles(team int) (tiles []Tile) {
	for y := 0; y < board.Size; y++ {
		for x := 0; x < board.Size; x++ {
			tile := board.GetTile(x, y)
			if tile.team == team {
				tiles = append(tiles, tile)
			}
		}
	}
	if tiles == nil {
		panic(fmt.Errorf("failed to locate any workers for team %d", team))
	}
	return
}

func (board *DefaultBoard) GetValidTurns(team int) (turns []Turn) {
	if team == 0 {
		return
	}
	// Get worker tiles
	workerTiles := board.GetWorkerTiles(team)

	// Check moves for each worker
	for _, workerTile := range workerTiles {
		// Get all valid tiles to move to
		moves := board.GetMoveableTiles(workerTile)

		// Add a potential turn for each possible build and move
		for _, move := range moves {
			builds := board.GetBuildableTiles(workerTile.team, workerTile.worker, move)
			for _, build := range builds {
				turns = append(turns, Turn{
					Team:   team,
					Worker: workerTile.worker,
					MoveTo: board.GetTile(move.x, move.y),
					Build:  board.GetTile(build.x, build.y),
				})
			}
		}
	}

	board.Teams[team] = len(turns) != 0
	return
}
