package board

import "fmt"

type Turn struct {
	Worker *Tile
	MoveTo *Tile
	Build  *Tile
}

func (t Turn) String() string {
	return fmt.Sprintf("%c%d%c%d%c%d",
		rune(t.Worker.x+65),
		t.Worker.y,
		rune(t.MoveTo.x+65),
		t.MoveTo.y,
		rune(t.Build.x+65),
		t.Build.y,
	)
}

type Board interface {
	// Functions needed to play the game
	PlayTurn(*Turn) (victory bool, err error)
	UndoTurn(t *Turn) error
	GetTile(x, y uint8) (t *Tile)
	Clone() Board
	Hash() string
	// Get a list of all the eligible turns that a team can take
	ValidTurns(team uint8) []*Turn
}
