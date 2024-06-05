package santorini

import (
	"fmt"

	"go.starlark.net/starlark"
)

/* TStarlark functions for the board object*/

// Functions needed for starlark.Value
func (f FastBoard) String() string {
	return f.GameHash()
}
func (f FastBoard) Type() string {
	return "Board"
}
func (f FastBoard) Freeze() {
}
func (f FastBoard) Truth() starlark.Bool {
	return starlark.True
}
func (f FastBoard) Hash() (uint32, error) {
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

func (f FastBoard) Attr(name string) (starlark.Value, error) {
	switch name {
	case "width":
		return starlark.MakeInt(int(f.width)), nil
	case "height":
		return starlark.MakeInt(int(f.height)), nil
	case "get_tile":
		return starlark.NewBuiltin(name, f.get_tile), nil
	case "get_surrounding_tiles":
		return starlark.NewBuiltin(name, f.get_surrounding_tiles), nil
	}
	return nil, fmt.Errorf("not found")
}

func (f FastBoard) AttrNames() []string {
	return []string{"width", "height", "get_tile", "get_surrounding_tiles"}
}

func (f *FastBoard) get_tile(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var x, y starlark.Int
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &x, &y); err != nil {
		return nil, err
	}

	i, ok := x.Int64()
	if !ok {
		return nil, fmt.Errorf("int value %s is 'not exactly representable'", x.String())
	}
	j, ok := y.Int64()
	if !ok {
		return nil, fmt.Errorf("int value %s is 'not exactly representable'", y.String())
	}
	return f.GetTile(uint8(i), uint8(j)), nil
}

func (f *FastBoard) get_surrounding_tiles(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var i, j starlark.Int
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 2, &i, &j); err != nil {
		return nil, err
	}

	x64, ok := i.Int64()
	if !ok {
		return nil, fmt.Errorf("int value %s is 'not exactly representable'", i.String())
	}
	y64, ok := j.Int64()
	if !ok {
		return nil, fmt.Errorf("int value %s is 'not exactly representable'", j.String())
	}

	x := uint8(x64)
	y := uint8(y64)

	tiles := starlark.NewList(nil)
	// List all surrounding tiles
	if y > 0 {
		// North
		index := f.board[(f.width*(y-1))+x]
		tiles.Append(Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x,
			y:      y - 1,
		})
	}
	if y < f.height-1 {
		// South
		index := f.board[(f.width*(y+1))+x]
		tiles.Append(Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x,
			y:      y + 1,
		})
	}
	if x > 0 {
		// West
		index := f.board[(f.width*(y))+x-1]
		tiles.Append(Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x - 1,
			y:      y,
		})
	}
	if x < f.width-1 {
		// East
		index := f.board[(f.width*(y))+x+1]
		tiles.Append(Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x + 1,
			y:      y,
		})
	}
	if y > 0 && x < f.width-1 {
		// NorthEast
		index := f.board[(f.width*(y-1))+x+1]
		tiles.Append(Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x + 1,
			y:      y - 1,
		})
	}
	if y > 0 && x > 0 {
		// NorthWest
		index := f.board[(f.width*(y-1))+x-1]
		tiles.Append(Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x - 1,
			y:      y - 1,
		})
	}
	if y < f.height-1 && x < f.width-1 {
		// SouthEast
		index := f.board[(f.width*(y+1))+x+1]
		tiles.Append(Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x + 1,
			y:      y + 1,
		})
	}
	if y < f.height-1 && x > 0 {
		// SouthEast
		index := f.board[(f.width*(y+1))+x-1]
		tiles.Append(Tile{
			team:   index >> 3,
			height: index & 0x7,
			x:      x - 1,
			y:      y + 1,
		})
	}

	return tiles, nil
}
