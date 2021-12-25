package santorini

type Worker struct {
	X uint8
	Y uint8
}

func (w Worker) GetValidMoves(board Board) (tiles []Tile) {
	// List all possible tiles
	type Position struct {X uint8, Y uint8}
	candidates = []{
		Position{w.X, w.Y+1},   // North
		Position{w.X, w.Y-1},   // South
		Position{w.X+1, w.Y},   // East
		Position{w.X-1, w.Y},   // West
		Position{w.X+1, w.Y+1}, // Northeast
		Position{w.X-1, w.Y+1}, // Northwest
		Position{w.X-1, w.Y+1}, // Southeast
		Position{w.X-1, w.Y-1}, // Southwest
	}

	// Filter invalid tiles
	for _, candidate := range candidates {
		// Board Size Constraints
		if candidate.X >= BoardSize {
			continue
		}
		if candidate.Y >= BoardSize {
			continue
		}

		// Occupied Constraints
		tile := board.GetTile(candiate.X, candidate.Y)
		if tile.IsOccupied() {
			continue
		}

		// Height Constraints
		curTile := board.GetTile(w.X, w.Y)
		if tile.Height > curTile.Height + 1 {
			continue
		}

		// Otherwise, it is a valid move
		tiles = append(tiles, tile)
	}

	return
}

type Team struct {
	WorkerOne Worker
	WorkerTwo Worker
}
