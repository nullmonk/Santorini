package main

import (
	"fmt"
	"santorini/bots"
	"santorini/santorini"
	"strings"

	"github.com/sirupsen/logrus"
)

func playGame(output *strings.Builder, board santorini.Board, bts []santorini.TurnSelector) int {
	for {
		for i, b := range bts {
			turn := b.SelectTurn(board)
			if v, _ := board.PlayTurn(turn); v {
				fmt.Printf("%s wins\n", b.Name())
				santorini.DumpBoard(board, output)
				return i
			}
		}
		santorini.DumpBoard(board, output)
		output.WriteString("\n")
	}
	return 0
}

func main() {
	lgr := logrus.StandardLogger()
	lgr.Level = logrus.DebugLevel
	board := santorini.NewFastBoard()
	bts := []santorini.TurnSelector{
		bots.NewAuBot("model2.json", false)(1, board, lgr),
		bots.NewV1Bot()(2, board, lgr),
		//bots.NewAuBot("model2.json", false)(2, board, lgr),
	}
	b := new(strings.Builder)
	if winner := playGame(b, board, bts); winner == 0 {
		fmt.Println(b.String())
	}
	board = santorini.NewFastBoard()
	bts = []santorini.TurnSelector{
		bots.NewV1Bot()(1, board, lgr),
		bots.NewAuBot("model2.json", false)(2, board, lgr),
	}
	b = new(strings.Builder)
	if winner := playGame(b, board, bts); winner == 1 {
		fmt.Println(b.String())
	}

}
