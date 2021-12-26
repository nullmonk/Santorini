package santorini

type Turn struct {
	Worker *Worker
	MoveTo Tile
	Build  Tile
}

type Worker struct {
	Team   int
	Number int

	X uint8
	Y uint8
}

func (worker *Worker) GetValidTurns(board Board) (turns []Turn) {
	// Get all valid tiles to move to
	moves := board.GetMoveableTiles(worker.X, worker.Y)

	// Add a potential turn for each possible build and move
	for _, move := range moves {
		builds := board.GetBuildableTiles(move.x, move.y, *worker)
		for _, build := range builds {
			turns = append(turns, Turn{
				Worker: worker,
				MoveTo: move,
				Build:  build,
			})
		}
	}

	return
}
