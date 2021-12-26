package santorini

import (
	"fmt"
	"strings"

	"santorini/pkg/color"
)

func (board Board) String() string {
	rows := make([]string, board.Size)
	for y := 0; y < board.Size; y++ {
		columns := make([]string, board.Size)
		for x := 0; x < board.Size; x++ {
			tile := board.GetTile(x, y)
			display := fmt.Sprintf("%s%d%s", color.GetWorkerColor(tile.team, tile.worker), tile.height, color.Reset)
			columns = append(columns, display)
		}
		rows[y] = strings.Join(columns, " ")
	}

	return strings.Join(rows, "\n")
}
