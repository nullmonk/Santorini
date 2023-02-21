package santorini

type Board interface {
	// Functions needed to control flow of the game
	PlayTurn(*Turn) (victory bool, err error)
	UndoTurn(t *Turn) error
	setTile(team, height, x, y uint8) error
	Clone() Board

	// Read only functions for bots to use

	// GetWorkers returns the Tiles that the workers currently reside on
	GetWorkers(team uint8) []Tile
	// GetSurroundingTiles returns the tiles surounding the given tiles
	GetSurroundingTiles(x, y uint8) []Tile
	// Get all the valid turns that a team can make
	ValidTurns(team uint8) []*Turn
	// Get information about a tile
	GetTile(x, y uint8) (t Tile)

	// Information about the board
	Teams() []uint8
	Hash() string
	Dimensions() (uint8, uint8)
}
