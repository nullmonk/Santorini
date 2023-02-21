package santorini

import (
	"fmt"
	"strconv"
)

type Turn struct {
	Worker Tile
	MoveTo Tile
	Build  Tile
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

func (t *Turn) IsWinningMove() bool {
	return t.MoveTo.GetHeight() == 3
}
