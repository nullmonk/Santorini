package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"santorini/santorini"
	"strings"
)

// Visualize a game log

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: <game file> [pause]")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var b *santorini.FastBoard
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
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
			fmt.Println("\n" + line)
			fmt.Println(DumpBoardMini(b))
			continue
		}
		fmt.Println(scanner.Text())
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
