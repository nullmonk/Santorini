package santorini

import (
	"fmt"
	"strings"

	"santorini/pkg/color"
)

func (tile Tile) String() (display string) {
	display = fmt.Sprintf("%d", tile.Height)

	if tile.Worker == nil {
		return
	}

	display = fmt.Sprintf("%s%s%s", color.GetWorkerColor(tile.Worker.Team, tile.Worker.Number), display, color.Reset)
	return
}

func (board Board) String() string {
	rows := make([]string, board.Size)
	for x := 0; x < int(board.Size); x++ {
		columns := make([]string, board.Size)
		for y := 0; y < int(board.Size); y++ {
			tile := board.GetTile(uint8(x), uint8(y))
			columns = append(columns, tile.String())
		}
		rows[x] = strings.Join(columns, " ")
	}

	return strings.Join(rows, "\n")
}
