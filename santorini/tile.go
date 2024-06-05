package santorini

/* implement a tile object that is accessible by the starlark script */
import (
	"encoding/json"
	"fmt"

	"go.starlark.net/starlark"
)

type Tile struct {
	team   uint8 // 0 no team, otherwise team number
	height uint8 // 0 no building, 4 capped
	x      uint8 // x position of the tile
	y      uint8 // y position of the tile
}

// Functions needed for starlark.Value
func (t Tile) String() string {
	s, _ := t.MarshalJSON()
	return string(s)
}
func (t Tile) Type() string {
	return "Tile"
}
func (t Tile) Freeze() {
}
func (t Tile) Truth() starlark.Bool {
	return starlark.True
}
func (t Tile) Hash() (uint32, error) {
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
func (t Tile) Attr(name string) (starlark.Value, error) {
	switch name {
	case "team":
		return starlark.MakeInt(int(t.team)), nil
	case "height":
		return starlark.MakeInt(int(t.height)), nil
	case "x":
		return starlark.MakeInt(int(t.x)), nil
	case "y":
		return starlark.MakeInt(int(t.y)), nil
	}
	return nil, fmt.Errorf("not found")
}

func (t Tile) AttrNames() []string {
	return []string{"team", "height", "x", "y"}
}

func (t *Tile) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Team   uint8 `json:"team,omitempty"`
		Height uint8 `json:"height,omitempty"`
		X      uint8 `json:"x"`
		Y      uint8 `json:"y"`
	}{
		t.team, t.height, t.x, t.y,
	})
}

func (t *Tile) UnmarshalJSON(b []byte) error {
	x := struct {
		Team   uint8 `json:"team,omitempty"`
		Height uint8 `json:"height,omitempty"`
		X      uint8 `json:"x"`
		Y      uint8 `json:"y"`
	}{}

	err := json.Unmarshal(b, &x)
	t.height = x.Height
	t.team = x.Team
	t.x = x.X
	t.y = x.Y
	return err
}

func (t Tile) IsCapped() bool {
	return t.height > 3
}

func (t Tile) IsOccupied() bool {
	return t.height > 3 || t.team > 0
}

func (t Tile) GetX() uint8 {
	return t.x
}

func (t Tile) GetY() uint8 {
	return t.y
}

func (t Tile) GetTeam() uint8 {
	return t.team
}

func (t Tile) GetHeight() uint8 {
	return t.height
}

func (t Tile) SameLocation(t2 Tile) bool {
	return t.x == t2.x && t.y == t2.y
}

func (t Tile) Equal(o Tile) bool {
	return t.x == o.x && t.y == o.y && t.team == o.team && t.height == o.height
}

// CanMoveTo says if t is able to move to t2
func (t Tile) CanMoveTo(t2 Tile) error {
	dist := getDistance(t, t2)
	// cant move to same spot or a far distance
	if dist == 0 || dist >= 2 {
		return fmt.Errorf("the worker cannot move to the given block")
	}
	// The worker cannot jump 2 blocks
	if t.height < t2.height && t2.height-t.height > 1 {
		return fmt.Errorf("the worker cannot jump 2 blocks")
	}
	if t2.IsOccupied() {
		return fmt.Errorf("the worker cannot move to an occupied block")
	}
	return nil
}

// HeightDiff if >0, the worker need to go up, if <0, the worker will jump down
func (t Tile) HeightDiff(t2 Tile) int {
	return int(t2.height) - int(t.height)
}
