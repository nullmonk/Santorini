package board

type Turn struct {
	Worker Tile
	MoveTo Tile
	Build  Tile
}

type Board interface {
	// Functions needed to play the game
	PlayTurn(Turn) (victory bool, err error)
	GetTile(x, y uint8) (t Tile)
	Clone() Board
	// Get a list of all the eligible turns that a team can take
	ValidTurns(team uint8) []Turn
}
