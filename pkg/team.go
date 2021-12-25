package santorini

type Worker struct {
	X uint8
	Y uint8
}

func (w Worker) GetValidMoves(board Board) (tiles []Tile) {
	// List all possible tiles
	type Position struct {
		X uint8
		Y uint8
	}
	candidates := []Position{
		{w.X, w.Y + 1},     // North
		{w.X, w.Y - 1},     // South
		{w.X + 1, w.Y},     // East
		{w.X - 1, w.Y},     // West
		{w.X + 1, w.Y + 1}, // Northeast
		{w.X - 1, w.Y + 1}, // Northwest
		{w.X - 1, w.Y + 1}, // Southeast
		{w.X - 1, w.Y - 1}, // Southwest
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
		tile := board.GetTile(candidate.X, candidate.Y)
		if tile.IsOccupied() {
			continue
		}

		// Height Constraints
		curTile := board.GetTile(w.X, w.Y)
		if tile.Height > curTile.Height+1 {
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
