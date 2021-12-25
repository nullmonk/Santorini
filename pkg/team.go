package santorini

type Team struct {
	WorkerOne Worker
	WorkerTwo Worker
}

type Turn struct {
	Worker
	MoveTo Tile
	Build  Tile
}

type Worker struct {
	X uint8
	Y uint8
}

func (worker Worker) GetValidTurns(board Board) (turns []Turn) {
	// Get all valid tiles to move to
	// moveTiles := board.GetMoveableTiles(worker.X, worker.Y)

	// Add a potential turn for each possible build and move
	// for _, tile := range moveTiles {

	// }

	return
}
