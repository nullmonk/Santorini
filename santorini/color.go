package santorini

import (
	"fmt"
	"runtime"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

func TileIcon(t Tile) string {
	color := []string{White, Blue, Red, Purple}[t.team]
	tileIcon := fmt.Sprintf("%d", t.height)
	if t.height == 4 {
		tileIcon = "^"
	}
	return fmt.Sprintf("%s%v%s", color, tileIcon, Reset)
}
