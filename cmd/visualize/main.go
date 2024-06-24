package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"santorini/santorini"
	"strings"
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

func TileIcon(t santorini.Tile) string {
	color := []string{White, Blue, Red, Purple}[t.GetTeam()]
	tileIcon := fmt.Sprintf("%d", t.GetHeight())
	if t.GetHeight() == 4 {
		tileIcon = "^"
	}
	return fmt.Sprintf("%s%v%s", color, tileIcon, Reset)
}

func printlnwithcolorreplacement(color map[string]string, line string) {
	for orig, repl := range color {
		line = strings.ReplaceAll(line, orig, repl)
	}
	fmt.Println(line)
}

// Visualize a game log

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: <game file> [pause]")
		os.Exit(1)
	}

	var team_colors map[string]string

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var b *santorini.FastBoard
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "TEAMS: ") {
			tms := strings.Split(strings.TrimPrefix(line, "TEAMS: "), " ")
			team_colors = make(map[string]string, len(tms))
			// Replace "team" with "[COLOR]team[RESET]"
			for i, tm := range tms {
				color := []string{White, Blue, Red, Purple}[i+1]
				team_colors[tm] = fmt.Sprintf("%s%v%s", color, tm, Reset)
			}
		}
		if strings.HasPrefix(line, "BOARD: ") {
			line = strings.TrimPrefix(line, "BOARD: ")
			b, err = santorini.NewBoardFromHash(line)
			if err != nil {
				log.Fatal(fmt.Errorf("error loading board from hash %s: %s", line, err))
			}
			continue
		}
		if strings.HasPrefix(line, "TURN: ") {
			turn := line[6:12]
			t, err := santorini.TurnFromString(0, turn)
			if err != nil {
				log.Fatal(fmt.Errorf("invalid turn %s: %s", turn, err))
			}
			if _, err := b.PlayTurn(t); err != nil {
				log.Fatal(fmt.Errorf("error playing turn %s: %s", turn, err))
			}
			fmt.Println()
			printlnwithcolorreplacement(team_colors, line)
			fmt.Println(DumpBoardMini(b))
			continue
		}
		printlnwithcolorreplacement(team_colors, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func DumpBoardMini(b *santorini.FastBoard) string {
	width, height := b.Dimensions()
	line := ""
	for y := uint8(0); y < height; y++ {
		for x := uint8(0); x < width; x++ {
			tile := b.GetTile(x, y)
			tileIcon := TileIcon(tile)
			line += tileIcon
		}
		line += "\n"
	}
	return strings.TrimSpace(line)
}
