package santorini

import (
	"fmt"
	"strconv"

	"go.starlark.net/starlark"
)

type Turn struct {
	Worker Tile
	MoveTo Tile
	Build  Tile
	Rank   int // Used by bots to rank how important this move is
}

func TurnFromString(team uint8, s string) (*Turn, error) {
	if len(s) != 6 {
		return nil, fmt.Errorf("not a valid turn string")
	}
	turn := &Turn{
		Worker: Tile{
			team: team,
		},
		MoveTo: Tile{},
		Build:  Tile{},
	}
	turn.Worker.x = s[0] - 65
	i, err := strconv.Atoi(string(s[1]))
	if err != nil {
		return nil, err
	}
	turn.Worker.y = uint8(i)

	turn.MoveTo.x = s[2] - 65
	i, err = strconv.Atoi(string(s[3]))
	if err != nil {
		return nil, err
	}
	turn.MoveTo.y = uint8(i)

	turn.Build.x = s[4] - 65
	i, err = strconv.Atoi(string(s[5]))
	turn.Build.y = uint8(i)

	return turn, err
}

func (t *Turn) IsWinningMove() bool {
	return t.MoveTo.GetHeight() == 3
}

// Functions needed for starlark.Value
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

func (t Turn) Type() string {
	return "Turn"
}
func (t Turn) Freeze() {
}
func (t Turn) Truth() starlark.Bool {
	return starlark.True
}
func (t Turn) Hash() (uint32, error) {
	return 0, fmt.Errorf("cannot hash")
}

// Functions needed for starlark.HasAttr
/*
type HasAttrs interface {
	Value
	Attr(name string) (Value, error) // returns (nil, nil) if attribute not present
	AttrNames() []string             // callers must not modify the result.
}
*/

// All the things accessible from this object
func (t Turn) Attr(name string) (starlark.Value, error) {
	switch name {
	case "worker":
		return &t.Worker, nil
	case "move":
		return &t.MoveTo, nil
	case "build":
		return &t.Build, nil
	case "is_winning_move":
		if t.MoveTo.height == 3 {
			return starlark.True, nil
		} else {
			return starlark.False, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (t Turn) AttrNames() []string {
	return []string{"worker", "move", "build", "is_winning_move"}
}
